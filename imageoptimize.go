package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"bytes"

	"io/ioutil"

	"math"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	originalImage, _ := OpenFile("sample.png")
	//resizedImage, _ := ResizeAndCompress(originalImage, 500, 500, AspectFit, VerticalAlignmentCenter, HorizontalAlignmentCenter)
	resizedImage, _ := ThumbnailAndCompress(originalImage, 100, 100)
	file, _ := os.Create(`sample-resized.png`)
	defer file.Close()
	file.Write(resizedImage)
}

type ContentMode int

const (
	ScaleToFill ContentMode = iota
	AspectFit
	AspectFill
)

type VerticalAlignment int

const (
	VerticalAlignmentTop VerticalAlignment = iota
	VerticalAlignmentBottom
	VerticalAlignmentCenter
)

type HorizontalAlignment int

const (
	HorizontalAlignmentLeft HorizontalAlignment = iota
	HorizontalAlignmentRight
	HorizontalAlignmentCenter
)

type OriginalImage struct {
	filePath    string
	mineType    string
	decImage    image.Image // for jpg,png
	gifImage    *gif.GIF    // for animation gif
	imageConfig image.Config
}

func GenerateOriginalImage(file []byte) (OriginalImage, error) {
	var err error
	mineType := http.DetectContentType(file)
	imageFile := bytes.NewReader(file)

	var decImage image.Image
	var gifImage *gif.GIF
	var imageConfig image.Config
	switch mineType {
	case "image/jpeg":
		decImage, err = jpeg.Decode(imageFile)
		if err != nil {
			return OriginalImage{}, err
		}
		_, err = imageFile.Seek(io.SeekStart, 0)
		if err != nil {
			return OriginalImage{}, err
		}
		imageConfig, err = jpeg.DecodeConfig(imageFile)
		if err != nil {
			return OriginalImage{}, err
		}
	case "image/png":
		decImage, err = png.Decode(imageFile)
		if err != nil {
			return OriginalImage{}, err
		}
		_, err = imageFile.Seek(io.SeekStart, 0)
		if err != nil {
			return OriginalImage{}, err
		}
		imageConfig, err = png.DecodeConfig(imageFile)
		if err != nil {
			return OriginalImage{}, err
		}
	case "image/gif":
		gifImage, err = gif.DecodeAll(imageFile)
		if err != nil {
			return OriginalImage{}, err
		}
		imageConfig = gifImage.Config
		if err != nil {
			return OriginalImage{}, err
		}
	default:
		return OriginalImage{}, errors.New("Unsupported file type")
	}

	originalImage := OriginalImage{
		mineType:    mineType,
		decImage:    decImage,
		gifImage:    gifImage,
		imageConfig: imageConfig,
	}
	return originalImage, nil
}

func OpenFile(filePath string) (OriginalImage, error) {
	imageFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return OriginalImage{}, err
	}
	return GenerateOriginalImage(imageFile)
}

//
func scaleToFillSize(imageConfig image.Config, resizedWidth uint, resizedHeight uint) (uint, uint) {
	return resizedWidth, resizedHeight
}
func aspectFitSize(imageConfig image.Config, resizedWidth uint, resizedHeight uint) (uint, uint) {
	scaleWidth := float64(resizedWidth) / float64(imageConfig.Width)
	scaleHeight := float64(resizedHeight) / float64(imageConfig.Height)
	if scaleWidth > scaleHeight {
		return uint(float64(imageConfig.Width) * scaleHeight), resizedHeight
	}
	return resizedWidth, uint(float64(imageConfig.Height) * scaleWidth)
}
func aspectFillSize(imageConfig image.Config, resizedWidth uint, resizedHeight uint) (uint, uint) {
	scaleWidth := float64(resizedWidth) / float64(imageConfig.Width)
	scaleHeight := float64(resizedHeight) / float64(imageConfig.Height)
	if scaleWidth > scaleHeight {
		return resizedWidth, uint(float64(imageConfig.Height) * scaleWidth)
	}
	return uint(float64(imageConfig.Width) * scaleHeight), resizedHeight
}

func calcAlignment(widthDiff int, heightDiff int, verticalAlignment VerticalAlignment, horizontalAlignment HorizontalAlignment) (int, int) {
	switch horizontalAlignment {
	case HorizontalAlignmentLeft:
		widthDiff = 0
	case HorizontalAlignmentCenter:
		widthDiff = widthDiff / 2
	case HorizontalAlignmentRight:
		widthDiff = widthDiff
	}
	switch verticalAlignment {
	case VerticalAlignmentTop:
		heightDiff = 0
	case VerticalAlignmentCenter:
		heightDiff = heightDiff / 2
	case VerticalAlignmentBottom:
		heightDiff = heightDiff
	}
	return widthDiff, heightDiff
}

// resize and compress
func ResizeAndCompress(originalImage OriginalImage, width uint, height uint, contentMode ContentMode, verticalAlignment VerticalAlignment, horizontalAlignment HorizontalAlignment) ([]byte, error) {
	var resizedImageWidth, resizedImageHeight uint
	switch contentMode {
	case ScaleToFill:
		resizedImageWidth, resizedImageHeight = scaleToFillSize(originalImage.imageConfig, width, height)
	case AspectFit:
		resizedImageWidth, resizedImageHeight = aspectFitSize(originalImage.imageConfig, width, height)
	case AspectFill:
		resizedImageWidth, resizedImageHeight = aspectFillSize(originalImage.imageConfig, width, height)
	}

	resultImage := new(bytes.Buffer)

	switch originalImage.mineType {
	case "image/jpeg", "image/png":
		resizedImage := resize.Resize(resizedImageWidth, resizedImageHeight, originalImage.decImage, resize.Lanczos3)

		widthDiff, heightDiff := calcAlignment(int(resizedImageWidth-width), int(resizedImageHeight-height), verticalAlignment, horizontalAlignment)
		img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		// set white background color
		bgColor := color.RGBA{255, 255, 255, 0}
		rect := img.Rect
		for h := rect.Min.Y; h < rect.Max.Y; h++ {
			for v := rect.Min.X; v < rect.Max.X; v++ {
				img.Set(v, h, bgColor)
			}
		}

		resizedImageRect := image.Rectangle{image.Pt(0, 0), image.Pt(int(width), int(height))}
		draw.Draw(img, resizedImageRect, resizedImage, image.Pt(widthDiff, heightDiff), draw.Over)

		var err error
		switch originalImage.mineType {
		case "image/jpeg":
			err = jpeg.Encode(resultImage, img, nil)
		case "image/png":
			err = png.Encode(resultImage, img)
		}
		if err != nil {
			return nil, err
		}

	case "image/gif":
		var originalWidth, originalHeight int
		// resize
		for index, frame := range originalImage.gifImage.Image {
			rect := frame.Bounds()

			// Add colors from original gif image
			var tmpPalette color.Palette
			for x := 1; x <= rect.Dx(); x++ {
				for y := 1; y <= rect.Dy(); y++ {
					if !contains(tmpPalette, originalImage.gifImage.Image[index].At(x, y)) {
						tmpPalette = append(tmpPalette, originalImage.gifImage.Image[index].At(x, y))
					}
				}
			}

			if index == 0 {
				originalWidth = rect.Dx()
				originalHeight = rect.Dy()
			}

			// remove margins
			formatedPalette := image.NewPaletted(image.Rectangle{image.Pt(0, 0), image.Pt(int(originalWidth), int(originalHeight))}, tmpPalette)
			formatedImageRect := image.Rect(0, 0, originalWidth, originalHeight)
			draw.Draw(formatedPalette, formatedImageRect, frame.SubImage(rect), image.Pt(0, 0), draw.Over)

			// resize
			widthDiff, heightDiff := calcAlignment(int(resizedImageWidth-width), int(resizedImageHeight-height), verticalAlignment, horizontalAlignment)
			tmpImage := formatedPalette.SubImage(formatedPalette.Bounds())
			resizedPalette := image.NewPaletted(image.Rectangle{image.Pt(0, 0), image.Pt(int(width), int(height))}, tmpPalette)
			resizedImageRect := image.Rect(0, 0, int(width), int(height))
			resizedImage := resize.Resize(resizedImageWidth, resizedImageHeight, tmpImage, resize.NearestNeighbor)
			draw.Draw(resizedPalette, resizedImageRect, resizedImage, image.Pt(widthDiff, heightDiff), draw.Over)

			originalImage.gifImage.Image[index] = resizedPalette
		}

		// Set size to resized size
		originalImage.gifImage.Config.Width = int(width)
		originalImage.gifImage.Config.Height = int(height)
		err := gif.EncodeAll(resultImage, originalImage.gifImage)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unsupported file type")
	}
	return resultImage.Bytes(), nil
}

func ThumbnailAndCompress(originalImage OriginalImage, maxWidth uint, maxHeight uint) ([]byte, error) {
	widthScale := float64(maxWidth) / float64(originalImage.imageConfig.Width)
	heightScale := float64(maxHeight) / float64(originalImage.imageConfig.Height)
	if maxWidth == 0 {
		widthScale = 1
	}
	if maxHeight == 0 {
		heightScale = 1
	}
	minScale := math.Min(widthScale, heightScale)
	minScale = math.Min(1, minScale)
	width := uint(float64(originalImage.imageConfig.Width) * float64(minScale))
	height := uint(float64(originalImage.imageConfig.Height) * float64(minScale))
	return ResizeAndCompress(originalImage, width, height, ScaleToFill, VerticalAlignmentCenter, HorizontalAlignmentCenter)
}

func Compress(originalImage OriginalImage) ([]byte, error) {
	return ResizeAndCompress(originalImage, uint(originalImage.imageConfig.Width), uint(originalImage.imageConfig.Height), ScaleToFill, VerticalAlignmentCenter, HorizontalAlignmentCenter)
}

// Check if color is already in the Palette
func contains(colorPalette color.Palette, c color.Color) bool {
	for _, tmpColor := range colorPalette {
		if tmpColor == c {
			return true
		}
	}
	return false
}
