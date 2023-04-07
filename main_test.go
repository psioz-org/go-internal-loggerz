package main

import "testing"

func Test_sum(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case 1",
			args: args{
				a: 1,
				b: 1,
			},
			want: 2,
		},
		{
			name: "case 1",
			args: args{
				a: 99,
				b: 10000,
			},
			want: 10099,
		},
		{
			name: "case 1",
			args: args{
				a: 0,
				b: 1,
			},
			want: 1,
		},
		{
			name: "case 1",
			args: args{
				a: -9,
				b: 1,
			},
			want: -8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
