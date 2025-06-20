package contactsheet

type Layout interface {
	NextPosition(w, h int) (int, int, bool)
	TotalSize() (int, int)
	Reset()
}

type Margin interface {
	margin() (int, int, int, int, int, int)
}

type FixedGridLayout struct {
	tileWidth   int
	tileHeight  int
	cols        int
	rows        int
	innerMargin int
	outerMargin int
	index       int
}

func (layout *FixedGridLayout) NextPosition(w, h int) (int, int, bool) {
	isReset := layout.index == layout.cols*layout.rows
	if isReset {
		layout.Reset()
	}
	row := layout.index / layout.cols
	col := layout.index % layout.cols
	x := layout.outerMargin + layout.tileWidth*col + layout.innerMargin*col
	y := layout.outerMargin + layout.tileHeight*row + layout.innerMargin*row

	layout.index++
	return x, y, isReset
}

func (layout *FixedGridLayout) TotalSize() (int, int) {
	w := layout.tileWidth*layout.cols + layout.innerMargin*(layout.cols-1) + layout.outerMargin*2
	h := layout.tileHeight*layout.rows + layout.innerMargin*(layout.rows-1) + layout.outerMargin*2
	return w, h
}

func (layout *FixedGridLayout) Reset() {
	layout.index = 0
}

func NewFixedGridLayout(config FixedTileConfig) (FixedGridLayout, error) {
	return FixedGridLayout{
		tileWidth:   config.TileWidth + config.Padding*2,
		tileHeight:  config.TileHeight + config.Padding*2,
		cols:        config.TileCols,
		rows:        config.TileRows,
		innerMargin: config.InnerMargin,
		outerMargin: config.OuterMargin,
		index:       0,
	}, nil
}
