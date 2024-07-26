package countryphone

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// custom data
	indata := InData{
		"US": &InCountryData{
			CC: "1",
			Groups: []InPrefixGroup{
				{
					Name:       "Alaska",
					Type:       1,
					PrefixType: 1,
					Prefixes:   []string{"907"},
				},
			},
		},
	}

	data := New(indata)

	require.NotNil(t, data)
}

func TestNew_default(t *testing.T) {
	t.Parallel()

	data := New(nil)

	require.NotNil(t, data)
}

func TestData_NumberInfo(t *testing.T) {
	t.Parallel()

	indata := InData{
		"US": &InCountryData{
			CC: "1",
			Groups: []InPrefixGroup{
				{
					Name:       "Alaska",
					Type:       1,
					PrefixType: 1,
					Prefixes:   []string{"907"},
				},
				{
					Name:       "Arizona",
					Type:       1,
					PrefixType: 1,
					Prefixes:   []string{"480", "520", "602", "623", "928"},
				},
			},
		},
		"CA": &InCountryData{
			CC: "1",
			Groups: []InPrefixGroup{
				{
					Name:       "Manitoba",
					Type:       2,
					PrefixType: 1,
					Prefixes:   []string{"204", "431", "584"},
				},
				{
					Name:       "Nunavut",
					Type:       2,
					PrefixType: 1,
					Prefixes:   []string{"867"},
				},
			},
		},
		"JP": &InCountryData{
			CC: "81",
		},
		"__": &InCountryData{
			CC: "7",
			Groups: []InPrefixGroup{
				{
					Name:       "TEST",
					Type:       5,
					PrefixType: 7,
					Prefixes:   []string{},
				},
			},
		},
	}

	data := New(indata)

	require.NotNil(t, data)

	tests := []struct {
		name    string
		prefix  string
		want    *NumInfo
		wantErr bool
	}{
		{
			name:    "empty",
			prefix:  "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no match",
			prefix:  "999999",
			want:    nil,
			wantErr: true,
		},
		{
			name:   "US & CA",
			prefix: "100000",
			want: &NumInfo{
				Type: 0,
				Geo: []*GeoInfo{
					{
						Alpha2: "US",
						Area:   "",
						Type:   0,
					},
					{
						Alpha2: "CA",
						Area:   "",
						Type:   0,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "US - Alaska",
			prefix: "1907000",
			want: &NumInfo{
				Type: 1,
				Geo: []*GeoInfo{
					{
						Alpha2: "US",
						Area:   "Alaska",
						Type:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "US - Arizona",
			prefix: "1623000",
			want: &NumInfo{
				Type: 1,
				Geo: []*GeoInfo{
					{
						Alpha2: "US",
						Area:   "Arizona",
						Type:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "CA - Manitoba",
			prefix: "1431000",
			want: &NumInfo{
				Type: 1,
				Geo: []*GeoInfo{
					{
						Alpha2: "CA",
						Area:   "Manitoba",
						Type:   2,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "JP",
			prefix: "81234567890",
			want: &NumInfo{
				Type: 0,
				Geo: []*GeoInfo{
					{
						Alpha2: "JP",
						Area:   "",
						Type:   0,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "Artificial without prefix",
			prefix: "7123",
			want: &NumInfo{
				Type: 7,
				Geo: []*GeoInfo{
					{
						Alpha2: "__",
						Area:   "",
						Type:   0,
					},
					{
						Alpha2: "__",
						Area:   "TEST",
						Type:   5,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.NumberInfo(tt.prefix)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, got)

				return
			}

			require.NoError(t, err)
			require.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestData_NumberType(t *testing.T) {
	t.Parallel()

	data := New(InData{})

	require.NotNil(t, data)

	tests := []struct {
		name    string
		num     int
		want    string
		wantErr bool
	}{
		{
			name:    "out of bounds < 0",
			num:     -1,
			want:    "",
			wantErr: true,
		},
		{
			name:    "out of bounds > max",
			num:     8,
			want:    "",
			wantErr: true,
		},
		{
			name:    "first",
			num:     0,
			want:    "",
			wantErr: false,
		},
		{
			name:    "last",
			num:     7,
			want:    "other",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.NumberType(tt.num)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, got)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestData_AreaType(t *testing.T) {
	t.Parallel()

	data := New(InData{})

	require.NotNil(t, data)

	tests := []struct {
		name    string
		num     int
		want    string
		wantErr bool
	}{
		{
			name:    "out of bounds < 0",
			num:     -1,
			want:    "",
			wantErr: true,
		},
		{
			name:    "out of bounds > max",
			num:     6,
			want:    "",
			wantErr: true,
		},
		{
			name:    "first",
			num:     0,
			want:    "",
			wantErr: false,
		},
		{
			name:    "last",
			num:     5,
			want:    "other",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.AreaType(tt.num)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, got)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
