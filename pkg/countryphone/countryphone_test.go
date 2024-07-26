package countryphone

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// default data
	data := New(nil)

	require.NotNil(t, data)

	// custom data

	pdata := PrefixData{
		"1279": &NumInfo{
			Type: 1,
			Geo: []*GeoInfo{
				{
					Alpha2: "US",
					Area:   "California",
					Type:   1,
				},
			},
		},
	}

	data = New(pdata)

	require.NotNil(t, data)
}

func TestData_NumberInfo(t *testing.T) {
	t.Parallel()

	numInfo1 := &NumInfo{
		Type: 1,
		Geo: []*GeoInfo{
			{
				Alpha2: "US",
				Area:   "California",
				Type:   1,
			},
		},
	}

	numInfo2 := &NumInfo{
		Type: 1,
		Geo: []*GeoInfo{
			{
				Alpha2: "CA",
				Area:   "Quebec",
				Type:   2,
			},
		},
	}

	pdata := PrefixData{
		"1279": numInfo1,
		"1367": numInfo2,
	}

	data := New(pdata)

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
			name:    "first",
			prefix:  "1279000",
			want:    numInfo1,
			wantErr: false,
		},
		{
			name:    "last",
			prefix:  "1367000",
			want:    numInfo2,
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
			require.Equal(t, tt.want, got)
		})
	}
}

func TestData_NumberType(t *testing.T) {
	t.Parallel()

	data := New(PrefixData{})

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

	data := New(PrefixData{})

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
