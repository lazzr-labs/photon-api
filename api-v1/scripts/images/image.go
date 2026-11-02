package images

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"

	"github.com/kovidgoyal/imaging"
	"github.com/kovidgoyal/imaging/prism/meta/autometa"
)

const (
	maxInputBytes  = 20 * 1024 * 1024
	maxImagePixels = 50_000_000
	jpegQuality    = 82
)

func OptimizeJPEG(reader io.Reader, maxDimension int) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(reader, maxInputBytes+1))
	if err != nil {
		return nil, errors.New("unable to read image")
	}
	if len(data) > maxInputBytes {
		return nil, errors.New("image exceeds 20 MiB limit")
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("unsupported or corrupt image")
	}
	if config.Width <= 0 || config.Height <= 0 {
		return nil, errors.New("image dimensions must be positive")
	}
	if config.Width > maxImagePixels/config.Height {
		return nil, errors.New("image exceeds 50,000,000 pixel limit")
	}

	metadata, _, err := autometa.Load(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("unsupported or corrupt image")
	}
	if metadata != nil && metadata.HasFrames {
		return nil, errors.New("GIFs and animated images are not supported; please upload a still image")
	}

	decodedImage, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true), imaging.Backends(imaging.GO_IMAGE))
	if err != nil {
		return nil, errors.New("unsupported or corrupt image")
	}

	bounds := decodedImage.Bounds()
	if bounds.Dx() > maxDimension || bounds.Dy() > maxDimension {
		decodedImage = imaging.Fit(decodedImage, maxDimension, maxDimension, imaging.Lanczos)
		bounds = decodedImage.Bounds()
	}

	jpegImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(jpegImage, jpegImage.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(jpegImage, jpegImage.Bounds(), decodedImage, bounds.Min, draw.Over)

	var output bytes.Buffer
	if err := jpeg.Encode(&output, jpegImage, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, errors.New("unable to encode image")
	}

	return output.Bytes(), nil
}
