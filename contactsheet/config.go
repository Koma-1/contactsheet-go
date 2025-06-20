package contactsheet

import (
	"image/color"
)

type FixedTileConfig struct {
	TileWidth           int
	TileHeight          int
	TileRows            int
	TileCols            int
	InnerMargin         int
	OuterMargin         int
	Padding             int
	BackgroundColor     color.Color
	TileBackgroundColor color.Color
	InputDirectory      string
	OutputDirectory     string
	OutputPrefix        string
	Interpolator        string
	TileMode            string
}
