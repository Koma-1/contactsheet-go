package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Koma-1/contactsheet-go/contactsheet"
)

type userConfig struct {
	width               int
	height              int
	rows                int
	cols                int
	innerMargin         int
	outerMargin         int
	padding             int
	backgroundColor     string
	tileBackgroundColor string
	indir               string
	outdir              string
	prefix              string
	tileMode            string
	interpolator        string
}

var defaultConfig = userConfig{
	width:               320,
	height:              320,
	rows:                5,
	cols:                5,
	innerMargin:         10,
	outerMargin:         20,
	padding:             10,
	backgroundColor:     "white",
	tileBackgroundColor: "#d3d3d3",
	prefix:              "out_",
	tileMode:            "fit",
	interpolator:        "CatmullRom",
}

func (c *userConfig) validate() error {
	if c.width <= 0 {
		return fmt.Errorf("width must be positive, got %d", c.width)
	}
	if c.height <= 0 {
		return fmt.Errorf("height must be positive, got %d", c.height)
	}
	if c.rows <= 0 {
		return fmt.Errorf("rows must be positive, got %d", c.rows)
	}
	if c.cols <= 0 {
		return fmt.Errorf("cols must be positive, got %d", c.cols)
	}
	if c.innerMargin < 0 {
		return fmt.Errorf("innerMargin must be nonnegative, got %d", c.innerMargin)
	}
	if c.outerMargin < 0 {
		return fmt.Errorf("outerMargin must be nonnegative, got %d", c.outerMargin)
	}
	if c.padding < 0 {
		return fmt.Errorf("padding must be nonnegative, got %d", c.padding)
	}
	if _, err := os.Stat(c.indir); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("input directory does not exist: %s", c.indir)
	}
	if _, err := os.Stat(c.outdir); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("output directory does not exist: %s", c.outdir)
	}
	return nil
}

func convertToFixedTileConfig(c userConfig) (contactsheet.FixedTileConfig, error) {
	// fmt.Println(fc)
	bg, err := parseColor(c.backgroundColor)
	if err != nil {
		return contactsheet.FixedTileConfig{}, fmt.Errorf("invalid background color: %w", err)
	}

	tbg, err := parseColor(c.tileBackgroundColor)
	if err != nil {
		return contactsheet.FixedTileConfig{}, fmt.Errorf("invalid tile background color: %w", err)
	}

	return contactsheet.FixedTileConfig{
		TileWidth:           c.width,
		TileHeight:          c.height,
		TileRows:            c.rows,
		TileCols:            c.cols,
		InnerMargin:         c.innerMargin,
		OuterMargin:         c.outerMargin,
		Padding:             c.padding,
		BackgroundColor:     bg,
		TileBackgroundColor: tbg,
		InputDirectory:      c.indir,
		OutputDirectory:     c.outdir,
		OutputPrefix:        c.prefix,
		Interpolator:        c.interpolator,
		TileMode:            c.tileMode,
	}, nil
}
