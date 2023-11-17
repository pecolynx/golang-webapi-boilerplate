package main

import "testing"

func Test_a(t *testing.T) {
	type input struct {
		x int
	}
	type output struct {
		value int
	}
	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "positive value",
			input: input{
				x: 3,
			},
			output: output{
				value: 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := a(tt.input.x); got != tt.output.value {
				t.Errorf("a() = %v, want %v", got, tt.output.value)
			}
		})
	}
}
