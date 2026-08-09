package pathing

import "testing"

func TestUint128ShiftRight(t *testing.T) {
	tests := []struct {
		name string
		u    uint128
		n    uint
		want uint128
	}{
		{
			name: "zero",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    0,
			want: uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
		},
		{
			name: "one bit",
			u:    uint128{lo: 0, hi: 1},
			n:    1,
			want: uint128{lo: 0x80_00_00_00_00_00_00_00, hi: 0},
		},
		{
			name: "one byte",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    8,
			want: uint128{lo: 0x21_12_34_56_78_9a_bc_de, hi: 0x00_0f_ed_cb_a9_87_65_43},
		},
		{
			name: "32 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    32,
			want: uint128{lo: 0x87_65_43_21_12_34_56_78, hi: 0x00_00_00_00_0f_ed_cb_a9},
		},
		{
			name: "63 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    63,
			want: uint128{lo: 0x1f_db_97_53_0e_ca_86_42, hi: 0},
		},
		{
			name: "64 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    64,
			want: uint128{lo: 0x0f_ed_cb_a9_87_65_43_21, hi: 0},
		},
		{
			name: "65 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    65,
			want: uint128{lo: 0x07_f6_e5_d4_c3_b2_a1_90, hi: 0},
		},
		{
			name: "96 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    96,
			want: uint128{lo: 0x00_00_00_00_0f_ed_cb_a9, hi: 0},
		},
		{
			name: "127 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    127,
			want: uint128{lo: 0, hi: 0},
		},
		{
			name: "128 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    128,
			want: uint128{lo: 0, hi: 0},
		},
		{
			name: "one bit from lo",
			u:    uint128{lo: 1, hi: 0},
			n:    1,
			want: uint128{lo: 0, hi: 0},
		},
		{
			name: "one bit from hi",
			u:    uint128{lo: 0, hi: 0x80_00_00_00_00_00_00_00},
			n:    1,
			want: uint128{lo: 0, hi: 0x40_00_00_00_00_00_00_00},
		},
		{
			name: "carry from hi",
			u:    uint128{lo: 0, hi: 1},
			n:    4,
			want: uint128{lo: 0x10_00_00_00_00_00_00_00, hi: 0},
		},
		{
			name: "carry across word boundary",
			u:    uint128{lo: 0, hi: 0x80_00_00_00_00_00_00_00},
			n:    8,
			want: uint128{lo: 0, hi: 0x00_80_00_00_00_00_00_00},
		},
		{
			name: "63 bits / b",
			u:    uint128{lo: 0, hi: 1},
			n:    63,
			want: uint128{lo: 2, hi: 0},
		},
		{
			name: "64 bits with zero hi",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0},
			n:    64,
			want: uint128{lo: 0, hi: 0},
		},
		{
			name: "127 bits with top bit set",
			u:    uint128{lo: 0, hi: 0x80_00_00_00_00_00_00_00},
			n:    127,
			want: uint128{lo: 1, hi: 0},
		},
		{
			name: "128 bits / b",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    128,
			want: uint128{lo: 0, hi: 0},
		},
		{
			name: "129 bits",
			u:    uint128{lo: 0x12_34_56_78_9a_bc_de_f0, hi: 0x0f_ed_cb_a9_87_65_43_21},
			n:    129,
			want: uint128{lo: 0, hi: 0},
		},
		{
			name: "cross word boundary",
			u:    uint128{lo: 0, hi: 0x80_00_00_00_00_00_00_00},
			n:    65,
			want: uint128{lo: 0x40_00_00_00_00_00_00_00, hi: 0},
		},
	}

	for _, test := range tests {
		have := test.u.ShiftRight(test.n)
		want := test.want
		if have != want {
			t.Fatalf("%s: ShiftRight(%d)\nhave: %v\nwant: %v",
				test.name, test.n, have, want)
		}
	}
}
