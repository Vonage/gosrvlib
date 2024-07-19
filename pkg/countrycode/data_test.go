package countrycode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()

	data := New()

	require.NotNil(t, data)

	require.Len(t, data.dStatusByID, lenEnumStatus)
	require.Len(t, data.dStatusIDByName, lenEnumStatus)
	require.Len(t, data.dRegionByID, lenEnumRegion)
	require.Len(t, data.dRegionIDByCode, lenEnumRegion)
	require.Len(t, data.dRegionIDByName, lenEnumRegion)
	require.Len(t, data.dSubRegionByID, lenEnumSubRegion)
	require.Len(t, data.dSubRegionIDByCode, lenEnumSubRegion)
	require.Len(t, data.dSubRegionIDByName, lenEnumSubRegion)
	require.Len(t, data.dIntermediateRegionByID, lenEnumIntRegion)
	require.Len(t, data.dIntermediateRegionIDByCode, lenEnumIntRegion)
	require.Len(t, data.dIntermediateRegionIDByName, lenEnumIntRegion)
	require.Len(t, data.dCountryNamesByAlpha2ID, 249)
	require.Len(t, data.dCountryKeyByAlpha2ID, 676)
	require.Len(t, data.dAlpha2IDByAlpha3ID, 249)
	require.Len(t, data.dAlpha2IDByNumericID, 249)
	require.Len(t, data.dAlpha2IDsByRegionID, lenEnumRegion)
	require.Len(t, data.dAlpha2IDsBySubRegionID, lenEnumSubRegion)
	require.Len(t, data.dAlpha2IDsByIntermediateRegionID, lenEnumIntRegion)
	require.Len(t, data.dAlpha2IDsByStatusID, lenEnumStatus)
	require.Len(t, data.dAlpha2IDsByTLD, 250)
}

func Test_statusByID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      int
		exp     *enumData
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     &enumData{code: "0", name: "Unassigned"},
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumStatus - 1,
			exp:     &enumData{code: "6", name: "Formerly assigned"},
			wantErr: false,
		},
		{
			name:    "out of range too small",
			in:      -1,
			exp:     nil,
			wantErr: true,
		},
		{
			name:    "out of range too big",
			in:      lenEnumStatus,
			exp:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.statusByID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_statusIDByName(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "Unassigned",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "Formerly assigned",
			exp:     lenEnumStatus - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.statusIDByName(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_regionByID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      int
		exp     *enumData
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     &enumData{code: "", name: ""},
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumRegion - 1,
			exp:     &enumData{code: "150", name: "Europe"},
			wantErr: false,
		},
		{
			name:    "out of range too small",
			in:      -1,
			exp:     nil,
			wantErr: true,
		},
		{
			name:    "out of range too big",
			in:      lenEnumRegion,
			exp:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.regionByID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_regionIDByCode(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "150",
			exp:     lenEnumRegion - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "invalid",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.regionIDByCode(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_regionIDByName(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "Europe",
			exp:     lenEnumRegion - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "invalid",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.regionIDByName(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_subRegionByID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      int
		exp     *enumData
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     &enumData{code: "", name: ""},
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumSubRegion - 1,
			exp:     &enumData{code: "419", name: "Latin America and the Caribbean"},
			wantErr: false,
		},
		{
			name:    "out of range too small",
			in:      -1,
			exp:     nil,
			wantErr: true,
		},
		{
			name:    "out of range too big",
			in:      lenEnumSubRegion,
			exp:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.subRegionByID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_subRegionIDByCode(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "419",
			exp:     lenEnumSubRegion - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "invalid",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.subRegionIDByCode(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_subRegionIDByName(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "Latin America and the Caribbean",
			exp:     lenEnumSubRegion - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "invalid",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.subRegionIDByName(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_intermediateRegionByID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      int
		exp     *enumData
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     &enumData{code: "", name: ""},
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumIntRegion - 1,
			exp:     &enumData{code: "029", name: "Caribbean"},
			wantErr: false,
		},
		{
			name:    "out of range too small",
			in:      -1,
			exp:     nil,
			wantErr: true,
		},
		{
			name:    "out of range too big",
			in:      lenEnumIntRegion,
			exp:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.intermediateRegionByID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_intermediateRegionIDByCode(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "029",
			exp:     lenEnumIntRegion - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "invalid",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.intermediateRegionIDByCode(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_intermediateRegionIDByName(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      string
		exp     uint8
		wantErr bool
	}{
		{
			name:    "first",
			in:      "",
			exp:     0,
			wantErr: false,
		},
		{
			name:    "last",
			in:      "Caribbean",
			exp:     lenEnumIntRegion - 1,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      "invalid",
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.intermediateRegionIDByName(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_countryNamesByAlpha2ID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint16
		exp     *names
		wantErr bool
	}{
		{
			name: "first",
			in:   0x0024,
			exp: &names{
				en: "Andorra",
				fr: "Andorre (l')",
			},
			wantErr: false,
		},
		{
			name: "last",
			in:   0x0357,
			exp: &names{
				en: "Zimbabwe",
				fr: "Zimbabwe (le)",
			},
			wantErr: false,
		},
		{
			name:    "invalid zero",
			in:      0,
			exp:     nil,
			wantErr: true,
		},
		{
			name:    "invalid max",
			in:      0xFFFF,
			exp:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.countryNamesByAlpha2ID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_countryKeyByAlpha2ID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint16
		exp     uint64
		wantErr bool
	}{
		{
			name:    "first",
			in:      0x0021,
			exp:     0x2084000000000000,
			wantErr: false,
		},
		{
			name:    "last",
			in:      0x035A,
			exp:     0x2D68000000000000,
			wantErr: false,
		},
		{
			name:    "first official",
			in:      0x0024,
			exp:     0x10902E20294C0240, // 0 | 001 | 00001 00100 | 00001 01110 00100 | 0000010100 | 101 | 00110 | 000 | 00001 00100 | 0000
			wantErr: false,
		},
		{
			name:    "last official",
			in:      0x0357,
			exp:     0x1D5F572D98613570, // 0 | 001 | 11010 10111 | 11010 10111 00101 | 1011001100 | 001 | 10000 | 100 | 11010 10111 | 0000
			wantErr: false,
		},
		{
			name:    "invalid zero",
			in:      0,
			exp:     0,
			wantErr: true,
		},
		{
			name:    "invalid max",
			in:      0xFFFF,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.countryKeyByAlpha2ID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDByAlpha3ID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint16
		exp     uint16
		wantErr bool
	}{
		{
			name:    "first",
			in:      0x05C4,
			exp:     0x0024,
			wantErr: false,
		},
		{
			name:    "last",
			in:      0x6AE5,
			exp:     0x0357,
			wantErr: false,
		},
		{
			name:    "invalid zero",
			in:      0,
			exp:     0,
			wantErr: true,
		},
		{
			name:    "invalid max",
			in:      0xFFFF,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDByAlpha3ID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDByNumericID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint16
		exp     uint16
		wantErr bool
	}{
		{
			name:    "first",
			in:      0x0014,
			exp:     0x0024,
			wantErr: false,
		},
		{
			name:    "last",
			in:      0x02CC,
			exp:     0x0357,
			wantErr: false,
		},
		{
			name:    "invalid zero",
			in:      0,
			exp:     0,
			wantErr: true,
		},
		{
			name:    "invalid max",
			in:      0xFFFF,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDByNumericID(tt.in)

			require.Equal(t, tt.exp, got)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDsByRegionID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint8
		exp     int
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     429,
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumRegion - 1,
			exp:     51,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      lenEnumRegion,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDsByRegionID(tt.in)

			require.Len(t, got, tt.exp)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDsBySubRegionID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint8
		exp     int
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     429,
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumSubRegion - 1,
			exp:     52,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      lenEnumSubRegion,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDsBySubRegionID(tt.in)

			require.Len(t, got, tt.exp)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDsByIntermediateRegionID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint8
		exp     int
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     571,
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumIntRegion - 1,
			exp:     28,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      lenEnumIntRegion,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDsByIntermediateRegionID(tt.in)

			require.Len(t, got, tt.exp)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDsByStatusID(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint8
		exp     int
		wantErr bool
	}{
		{
			name:    "first",
			in:      0,
			exp:     321,
			wantErr: false,
		},
		{
			name:    "last",
			in:      lenEnumStatus - 1,
			exp:     14,
			wantErr: false,
		},
		{
			name:    "official",
			in:      1,
			exp:     249,
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      lenEnumStatus,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDsByStatusID(tt.in)

			require.Len(t, got, tt.exp)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_alpha2IDsByTLD(t *testing.T) {
	t.Parallel()

	data := New()

	tests := []struct {
		name    string
		in      uint16
		exp     int
		wantErr bool
	}{
		{
			name:    "first",
			in:      0x0024,
			exp:     1,
			wantErr: false,
		},
		{
			name:    "last",
			in:      0x0357,
			exp:     1,
			wantErr: false,
		},
		{
			name:    "usa",
			in:      0x02B3,
			exp:     2,
			wantErr: false,
		},
		{
			name:    "zero",
			in:      0,
			exp:     426,
			wantErr: false,
		},
		{
			name:    "invalid max",
			in:      0xFFFF,
			exp:     0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := data.alpha2IDsByTLD(tt.in)

			require.Len(t, got, tt.exp)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
