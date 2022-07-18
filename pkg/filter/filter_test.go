package filter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func strPtr(v string) *string {
	return &v
}

func TestParseRules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    [][]Rule
		wantErr bool
	}{
		{
			name: "success",
			json: `[
			  [
				{ "field": "name", "type": "exact", "value": "doe" },
				{ "field": "age", "type": "exact", "value": 42 }
			  ],
			  [
				{ "field": "address.country", "type": "regexp", "value": "EN|FR" }
			  ]
			]`,
			want: [][]Rule{
				{
					{Field: "name", Type: "exact", Value: "doe"},
					{Field: "age", Type: "exact", Value: 42.0},
				},
				{
					{Field: "address.country", Type: "regexp", Value: "EN|FR"},
				},
			},
			wantErr: false,
		},
		{
			name:    "error - invalid json",
			json:    `[`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := ParseRules(tt.json)

			if tt.wantErr {
				require.Error(t, err, "ParseRules() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, r, "Filtered = %v, want %v", r, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rules   [][]Rule
		opts    []Option
		wantErr bool
	}{
		{
			name: "success",
			rules: [][]Rule{{{
				Field: "Somefield",
				Type:  "exact",
				Value: "some value",
			}}},
			opts: []Option{
				func(v *processor) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "error - option error",
			rules: [][]Rule{{{
				Field: "Somefield",
				Type:  "exact",
				Value: "some value",
			}}},
			opts: []Option{
				func(v *processor) error {
					return errors.New("test error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p, err := New(tt.rules, tt.opts...)

			if tt.wantErr {
				require.Error(t, err, "New() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, p)
			}
		})
	}
}

func TestFilter_Apply(t *testing.T) {
	t.Parallel()

	type simpleStruct struct {
		StringField    string `json:"string_field"`
		IntField       int
		Float64Field   float64
		StringPtrField *string
		unexported     string
	}

	type complexStruct struct {
		Internal simpleStruct `json:"internal"`
	}

	type complexStructWithPtr struct {
		Internal *simpleStruct
	}

	type embeddingStruct struct {
		simpleStruct
	}

	trueRegex := Rule{
		Field: "",
		Type:  "regexp",
		Value: ".*",
	}
	falseRegex := Rule{
		Field: "",
		Type:  "regexp",
		Value: "$a",
	}

	tests := []struct {
		name     string
		filter   [][]Rule
		opts     []Option
		elements interface{}
		want     interface{}
		wantErr  bool
	}{
		{
			name: "success - nested string equal",
			elements: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
				{
					Internal: simpleStruct{
						StringField: "value 2",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "Internal.StringField",
				Type:  "exact",
				Value: "value 1",
			}}},
			want: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
			},
		},
		{
			name: "success - nested string different",
			elements: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
				{
					Internal: simpleStruct{
						StringField: "value 2",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "Internal.StringField",
				Type:  "different",
				Value: "value 1",
			}}},
			want: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 2",
					},
				},
			},
		},
		{
			name: "success - nested regex",
			elements: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
				{
					Internal: simpleStruct{
						StringField: "value 2",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "Internal.StringField",
				Type:  "regexp",
				Value: ".* 1",
			}}},
			want: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
			},
		},
		{
			name: "success - int equal",
			elements: &[]simpleStruct{
				{
					IntField: 42,
				},
				{
					IntField: 43,
				},
			},
			filter: [][]Rule{{{
				Field: "IntField",
				Type:  "exact",
				Value: 42,
			}}},
			want: &[]simpleStruct{
				{
					IntField: 42,
				},
			},
		},
		{
			name: "success - float64 equal",
			elements: &[]simpleStruct{
				{
					Float64Field: 42,
				},
				{
					Float64Field: 43,
				},
			},
			filter: [][]Rule{{{
				Field: "Float64Field",
				Type:  "exact",
				Value: 42,
			}}},
			want: &[]simpleStruct{
				{
					Float64Field: 42,
				},
			},
		},
		{
			name: "success - nil equal",
			elements: &[]simpleStruct{
				{
					StringPtrField: strPtr("value 1"),
				},
				{
					StringPtrField: nil,
				},
			},
			filter: [][]Rule{{{
				Field: "StringPtrField",
				Type:  "exact",
				Value: nil,
			}}},
			want: &[]simpleStruct{
				{
					StringPtrField: nil,
				},
			},
		},
		{
			name: "success - invalid filter value type", // TODO report error or filter?
			elements: &[]simpleStruct{
				{
					StringField: "value 1",
				},
			},
			filter: [][]Rule{{{
				Field: "StringField",
				Type:  "exact",
				Value: 42,
			}}},
			want: &[]simpleStruct{},
		},
		{
			name: "success - regexp with an int", // TODO report error or filter?
			elements: &[]simpleStruct{
				{
					IntField: 42,
				},
			},
			filter: [][]Rule{{{
				Field: "IntField",
				Type:  "regexp",
				Value: "42",
			}}},
			want: &[]simpleStruct{},
		},
		{
			name: "success - mismatched array",
			elements: &[]interface{}{
				complexStructWithPtr{
					Internal: &simpleStruct{
						StringField: "value 1",
					},
				},
				complexStruct{
					Internal: simpleStruct{
						StringField: "value 2",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "Internal.StringField",
				Type:  "exact",
				Value: "value 1",
			}}},
			want: &[]interface{}{
				complexStructWithPtr{
					Internal: &simpleStruct{
						StringField: "value 1",
					},
				},
			},
		},
		{
			name: "success - with field tags",
			elements: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
				{
					Internal: simpleStruct{
						StringField: "value 2",
					},
				},
			},
			opts: []Option{WithFieldNameTag("json")},
			filter: [][]Rule{{{
				Field: "internal.string_field",
				Type:  "exact",
				Value: "value 1",
			}}},
			want: &[]complexStruct{
				{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
			},
		},
		{
			name: "success - with embedding struct",
			elements: &[]embeddingStruct{
				{
					simpleStruct: simpleStruct{
						StringField: "value 1",
					},
				},
				{
					simpleStruct: simpleStruct{
						StringField: "value 2",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "StringField",
				Type:  "exact",
				Value: "value 1",
			}}},
			want: &[]embeddingStruct{
				{
					simpleStruct: simpleStruct{
						StringField: "value 1",
					},
				},
			},
		},
		{
			name:     "success - with root field selector",
			elements: &[]int{41, 42, 43},
			filter: [][]Rule{{{
				Field: "",
				Type:  "exact",
				Value: 42,
			}}},
			want: &[]int{42},
		},
		{
			name:     "success - with empty AND filter",
			elements: &[]int{41, 42, 43},
			filter:   [][]Rule{},
			want:     &[]int{41, 42, 43},
		},
		{
			name:     "success - with empty OR filter",
			elements: &[]int{41, 42, 43},
			filter:   [][]Rule{{}},
			want:     &[]int{},
		},
		{
			name:     "combination - true AND true",
			elements: &[]string{"a"},
			filter:   [][]Rule{{trueRegex}, {trueRegex}},
			want:     &[]string{"a"},
		},
		{
			name:     "combination - true AND false",
			elements: &[]string{"a"},
			filter:   [][]Rule{{trueRegex}, {falseRegex}},
			want:     &[]string{},
		},
		{
			name:     "combination - false AND true",
			elements: &[]string{"a"},
			filter:   [][]Rule{{falseRegex}, {trueRegex}},
			want:     &[]string{},
		},
		{
			name:     "combination - false AND false",
			elements: &[]string{"a"},
			filter:   [][]Rule{{falseRegex}, {falseRegex}},
			want:     &[]string{},
		},
		{
			name:     "combination - true OR false",
			elements: &[]string{"a"},
			filter:   [][]Rule{{trueRegex, falseRegex}},
			want:     &[]string{"a"},
		},
		{
			name:     "combination - true OR true",
			elements: &[]string{"a"},
			filter:   [][]Rule{{trueRegex, trueRegex}},
			want:     &[]string{"a"},
		},
		{
			name:     "combination - false OR true",
			elements: &[]string{"a"},
			filter:   [][]Rule{{falseRegex, trueRegex}},
			want:     &[]string{"a"},
		},
		{
			name:     "combination - false OR false",
			elements: &[]string{"a"},
			filter:   [][]Rule{{falseRegex, falseRegex}},
			want:     &[]string{},
		},
		{
			name:     "combination - (false OR true) AND (true OR false)",
			elements: &[]string{"a"},
			filter:   [][]Rule{{falseRegex, trueRegex}, {trueRegex, falseRegex}},
			want:     &[]string{"a"},
		},
		{
			name:     "combination - (false OR true) AND (false OR false)",
			elements: &[]string{"a"},
			filter:   [][]Rule{{falseRegex, trueRegex}, {falseRegex, falseRegex}},
			want:     &[]string{},
		},
		{
			name: "error - with field tag not found",
			elements: &[]interface{}{
				complexStruct{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
			},
			opts: []Option{WithFieldNameTag("json")},
			filter: [][]Rule{{{
				Field: "internal.invalid_field",
				Type:  "exact",
				Value: "value 1",
			}}},
			wantErr: true,
		},
		{
			name:     "error - not a pointer",
			elements: 42,
			filter: [][]Rule{{{
				Type: "exact",
			}}},
			wantErr: true,
		},
		{
			name:     "error - not a slice",
			elements: &simpleStruct{},
			filter: [][]Rule{{{
				Type: "exact",
			}}},
			wantErr: true,
		},
		{
			name: "error - unexported field",
			elements: &[]interface{}{
				complexStruct{
					Internal: simpleStruct{
						unexported: "value 1",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "Internal.unexported",
				Type:  "exact",
				Value: "value 1",
			}}},
			wantErr: true,
		},
		{
			name: "error - with nil item",
			elements: &[]interface{}{
				nil,
			},
			filter: [][]Rule{{{
				Field: "Somefield",
				Type:  "exact",
				Value: "value 1",
			}}},
			wantErr: true,
		},
		{
			name: "error - nested path not found",
			elements: &[]interface{}{
				complexStruct{
					Internal: simpleStruct{
						StringField: "value 1",
					},
				},
			},
			filter: [][]Rule{{{
				Field: "Internal.InvalidField",
				Type:  "exact",
				Value: "value 1",
			}}},
			wantErr: true,
		},
		{
			name: "error - nested path inside a basic type",
			elements: &[]interface{}{
				simpleStruct{
					StringField: "value 1",
				},
			},
			filter: [][]Rule{{{
				Field: "StringField.InvalidField",
				Type:  "exact",
				Value: "value 1",
			}}},
			wantErr: true,
		},
		{
			name: "error - invalid regex",
			elements: &[]simpleStruct{
				{
					StringField: "value 1",
				},
			},
			filter: [][]Rule{{{
				Field: "StringField",
				Type:  "regexp",
				Value: "(",
			}}},
			wantErr: true,
		},
		{
			name: "error - not a string",
			elements: &[]simpleStruct{
				{
					StringField: "value 1",
				},
			},
			filter: [][]Rule{{{
				Field: "StringField",
				Type:  "regexp",
				Value: 1,
			}}},
			wantErr: true,
		},
		{
			name: "error - invalid filter type",
			elements: &[]simpleStruct{
				{
					StringField: "value 1",
				},
			},
			filter: [][]Rule{{{
				Field: "StringField",
				Type:  "invalid filter type",
				Value: "value 1",
			}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f, err := New(tt.filter, tt.opts...)
			require.NoError(t, err)

			err = f.Apply(tt.elements)

			if tt.wantErr {
				require.Error(t, err, "Apply() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, tt.elements, "Filtered = %v, want %v", tt.elements, tt.want)
			}
		})
	}
}

func benchmarkFilterApply(b *testing.B, n int, json string, opts ...Option) {
	b.Helper()

	type simpleStruct struct {
		IntField     int
		Float64Field float64
		SomeField1   interface{}
		SomeField2   interface{}
		SomeField3   interface{}
		StringField  string `json:"string_field"`
	}

	rules, err := ParseRules(json)
	require.NoError(b, err)

	filter, err := New(rules, opts...)
	require.NoError(b, err)

	data := make([]simpleStruct, n)
	for i := 0; i < n; i++ {
		data[i] = simpleStruct{
			StringField: "hello world", // TODO use faker ?
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		dataCopy := make([]simpleStruct, len(data))
		copy(dataCopy, data)

		b.StartTimer()

		err := filter.Apply(&dataCopy)
		require.NoError(b, err)
	}
}

func BenchmarkFilter_Apply_Exact_100(b *testing.B) {
	benchmarkFilterApply(
		b,
		100,
		`[[{"field": "StringField", "type": "exact", "value": "hello world"}]]`,
	)
}

func BenchmarkFilter_Apply_Exact_1000(b *testing.B) {
	benchmarkFilterApply(
		b,
		1000,
		`[[{"field": "StringField", "type": "exact", "value": "hello world"}]]`,
	)
}

func BenchmarkFilter_Apply_Exact_10000(b *testing.B) {
	benchmarkFilterApply(
		b,
		10000,
		`[[{"field": "StringField", "type": "exact", "value": "hello world"}]]`,
	)
}

func BenchmarkFilter_Apply_Regexp_1000(b *testing.B) {
	benchmarkFilterApply(
		b,
		1000,
		`[[{"field": "StringField", "type": "regexp", "value": "hello.*"}]]`,
	)
}

func BenchmarkFilter_Apply_WithTagField_1000(b *testing.B) {
	benchmarkFilterApply(
		b,
		1000,
		`[[{"field": "string_field", "type": "exact", "value": "hello world"}]]`,
		WithFieldNameTag("json"),
	)
}
