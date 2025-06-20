package cmd

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

func parseColor(str string) (color.Color, error) {
	if str = strings.ToLower(str); str == "black" {
		return color.Black, nil
	} else if str == "white" {
		return color.White, nil
	} else if matched, err := regexp.MatchString("^#[0-9a-f]{6}$", str); err == nil && matched {
		num, err := strconv.ParseUint(strings.Trim(str, "#"), 16, 24)
		if err != nil {
			return nil, err
		}
		return color.RGBA{uint8(num >> 16 & 0xFF), uint8(num >> 8 & 0xFF), uint8(num & 0xFF), 0xFF}, nil
	} else if matched, err := regexp.MatchString("^#[0-9a-f]{8}$", str); err == nil && matched {
		num, err := strconv.ParseUint(strings.Trim(str, "#"), 16, 32)
		if err != nil {
			return nil, err
		}
		return color.NRGBA{uint8(num >> 24 & 0xFF), uint8(num >> 16 & 0xFF), uint8(num >> 8 & 0xFF), uint8(num & 0xFF)}, nil
	}
	return nil, fmt.Errorf("parseColor: invalid color format: %s", str)
}
