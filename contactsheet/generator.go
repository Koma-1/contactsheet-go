package contactsheet

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "golang.org/x/image/bmp"

	"golang.org/x/image/draw"
)

type ContactsheetGenerator struct {
	dst     draw.Image
	index   int
	counter int
	config  FixedTileConfig
	layout  Layout
	tilegen TileGenerator
}

func (gen *ContactsheetGenerator) push(img image.Image) error {
	src, err := gen.tilegen.GenerateTile(img, gen.config)
	if err != nil {
		return err
	}
	if gen.dst == nil {
		gen.newDst()
	}
	x, y, newpage := gen.layout.NextPosition(src.Bounds().Dx(), src.Bounds().Dy())

	if newpage {
		err := gen.save()
		if err != nil {
			return err
		}
		gen.newDst()
	}

	draw.Copy(gen.dst, image.Pt(x, y), src, src.Bounds(), draw.Over, nil)
	gen.counter++
	return nil
}

func (gen *ContactsheetGenerator) save() error {
	gen.index++
	filename := gen.config.OutputPrefix + strconv.FormatInt(int64(gen.index), 10) + ".png"
	path := filepath.Join(gen.config.OutputDirectory, filename)

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := png.Encode(f, gen.dst); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	log.Printf("save: %s\n", filename)
	return nil
}

func (gen *ContactsheetGenerator) newDst() {
	w, h := gen.layout.TotalSize()
	rect := image.Rect(0, 0, w, h)
	gen.dst = image.NewRGBA(rect)
	bg := image.NewUniform(gen.config.BackgroundColor)
	draw.Copy(gen.dst, image.Point{}, bg, gen.dst.Bounds(), draw.Over, nil)
}

func (gen *ContactsheetGenerator) CopyDst() (draw.Image, error) {
	dst := image.NewRGBA(gen.dst.Bounds())
	draw.Copy(dst, image.Point{}, gen.dst, gen.dst.Bounds(), draw.Over, nil)

	return dst, nil
}

func (gen *ContactsheetGenerator) ResetDst() {
	gen.dst = nil
}

func (gen *ContactsheetGenerator) GenerateFromDir() error {
	files, err := os.ReadDir(gen.config.InputDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !checkExt(file.Name()) {
			continue
		}
		path := filepath.Join(gen.config.InputDirectory, file.Name())
		img, err := decodeImage(path)
		if err != nil {
			return err
		}
		if err := gen.push(img); err != nil {
			return err
		}
	}
	if gen.counter == 0 {
		log.Println("No valid images found")
		return nil
	}
	if err := gen.save(); err != nil {
		return err
	}
	return nil
}

func NewGenerator(config FixedTileConfig) (*ContactsheetGenerator, error) {
	tilegen, err := NewTileGenerator(config.TileMode)
	if err != nil {
		return nil, err
	}
	layout, err := NewFixedGridLayout(config)
	if err != nil {
		return nil, err
	}
	return &ContactsheetGenerator{
		dst:     nil,
		index:   0,
		counter: 0,
		config:  config,
		layout:  &layout,
		tilegen: tilegen,
	}, nil
}

func decodeImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, err
}

func checkExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpeg" || ext == ".jpg" || ext == ".bmp"
}
