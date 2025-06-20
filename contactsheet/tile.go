package contactsheet

import (
	"fmt"
	"image"
	"strings"

	"golang.org/x/image/draw"
)

type TileGenerator interface {
	GenerateTile(img image.Image, config FixedTileConfig) (image.Image, error)
}

type Fit struct{}
type Crop struct{}

func NewTileGenerator(tileMode string) (TileGenerator, error) {
	// fmt.Println("tilemode:", tileMode)
	if tileMode == "fit" {
		return &Fit{}, nil
	} else if tileMode == "crop" {
		return &Crop{}, nil
	} else {
		return nil, fmt.Errorf("invalid tileMode")
	}
}

func getInterpolator(config FixedTileConfig) (draw.Interpolator, error) {
	switch strings.ToLower(config.Interpolator) {
	case "n", "nearestneighbor":
		return draw.NearestNeighbor, nil
	case "a", "approxbilinear":
		return draw.ApproxBiLinear, nil
	case "b", "bilinear":
		return draw.BiLinear, nil
	case "c", "catmullrom":
		return draw.CatmullRom, nil
	default:
		return nil, fmt.Errorf("invalid interpolator")
	}
}

func (strategy *Fit) GenerateTile(img image.Image, config FixedTileConfig) (image.Image, error) {
	rect := image.Rect(0, 0, config.TileWidth+config.Padding*2, config.TileHeight+config.Padding*2)
	ret := image.NewRGBA(rect)
	bg := image.NewUniform(config.TileBackgroundColor)
	ip, err := getInterpolator(config)
	if err != nil {
		return nil, err
	}
	draw.Copy(ret, image.Point{}, bg, ret.Bounds(), draw.Over, nil)

	thumbnail, err := fitImageKeepAspectRatio(img, config.TileWidth, config.TileHeight, ip)
	if err != nil {
		return nil, err
	}

	centeringOffset := image.Pt(0, 0)
	if thumbnail.Bounds().Dx() < config.TileWidth {
		centeringOffset.X = (config.TileWidth - thumbnail.Bounds().Dx()) / 2
	}
	if thumbnail.Bounds().Dy() < config.TileHeight {
		centeringOffset.Y = (config.TileHeight - thumbnail.Bounds().Dy()) / 2
	}

	draw.Copy(ret, image.Pt(config.Padding, config.Padding).Add(centeringOffset), thumbnail, thumbnail.Bounds(), draw.Over, nil)

	return ret, nil
}

func (strategy *Crop) GenerateTile(img image.Image, config FixedTileConfig) (image.Image, error) {
	rect := image.Rect(0, 0, config.TileWidth+config.Padding*2, config.TileHeight+config.Padding*2)
	ret := image.NewRGBA(rect)
	bg := image.NewUniform(config.TileBackgroundColor)
	ip, err := getInterpolator(config)
	if err != nil {
		return nil, err
	}
	draw.Copy(ret, image.Point{}, bg, ret.Bounds(), draw.Src, nil)

	thumbnail, err := cropImageKeepAspectRatio(img, config.TileWidth, config.TileHeight, ip)
	if err != nil {
		return nil, err
	}

	draw.Copy(ret, image.Pt(config.Padding, config.Padding), thumbnail, thumbnail.Bounds(), draw.Src, nil)

	return ret, nil
}

func fitImageKeepAspectRatio(img image.Image, width int, height int, interpolator draw.Interpolator) (image.Image, error) {
	ratio := float64(width) / float64(img.Bounds().Dx())
	if float64(img.Bounds().Dy())*ratio > float64(height) {
		ratio = float64(height) / float64(img.Bounds().Dy())
	}

	dst := image.NewRGBA(image.Rect(0, 0, int(float64(img.Bounds().Dx())*ratio), int(float64(img.Bounds().Dy())*ratio)))
	interpolator.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst, nil
}

func cropRect(srcWidth, srcHeight, dstWidth, dstHeight int) image.Rectangle {
	srcRatio := float64(srcWidth) / float64(srcHeight)
	dstRatio := float64(dstWidth) / float64(dstHeight)
	if srcRatio > dstRatio {
		// srcが横長
		cropWidth := int(float64(srcHeight) * dstRatio)
		cropLeft := (srcWidth - cropWidth) / 2
		return image.Rect(cropLeft, 0, cropLeft+cropWidth, srcHeight)
	} else {
		// srcが縦長
		cropHeight := int(float64(srcWidth) / dstRatio)
		cropTop := (srcHeight - cropHeight) / 2
		return image.Rect(0, cropTop, srcWidth, cropTop+cropHeight)
	}
}

func cropImageKeepAspectRatio(img image.Image, width int, height int, interpolator draw.Interpolator) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	interpolator.Scale(dst, dst.Bounds(), img, cropRect(img.Bounds().Dx(), img.Bounds().Dy(), width, height), draw.Over, nil)
	return dst, nil
}
