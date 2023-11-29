package utility

import "testing"

func TestInArray(t *testing.T) {
	type args struct {
		search interface{}
		array  interface{}
		deep   bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test inarray", args: args{
				search: 1,
				array:  []string{"1", "2", "3"},
				deep:   true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InArray(tt.args.search, tt.args.array, tt.args.deep); got != tt.want {
				t.Errorf("InArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
