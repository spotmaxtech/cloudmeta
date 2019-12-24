package main

import "testing"

func Test_validInstance(t *testing.T) {
	type args struct {
		inst InstanceProduct
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validInstance(tt.args.inst); got != tt.want {
				t.Errorf("validInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}