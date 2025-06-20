package cmd

import (
	"image/color"
	"testing"
)

func TestColor1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  color.Color
	}{
		{name: "1", input: "black", want: color.Black},
		{name: "2", input: "white", want: color.White},
		{name: "3", input: "#123456", want: color.RGBA{0x12, 0x34, 0x56, 0xFF}},
		{name: "4", input: "#12345678", want: color.NRGBA{0x12, 0x34, 0x56, 0x78}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor(tt.input)
			if err != nil {
				t.Error(err)
			}
			if got != tt.want {
				t.Errorf("%+v, want %+v", got, tt.want)
			}
		})
	}
}
func TestColor2(t *testing.T) {
	got, err := parseColor("#123456789")
	if err == nil {
		t.Errorf("%+v, want nil", got)
	}
}
