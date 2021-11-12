package process

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

func BrightImage(buffer []byte, brightValue *int64, simpleType string) []byte {
	if brightValue == nil {
		return buffer
	}

	v := int32(*brightValue)
	if v == 0 {
		return buffer
	}

	isPNG, isJPEG, isBMP, isGIF, isWebp := checkImageType(simpleType)

	if !(isPNG || isJPEG || isGIF || isBMP || isWebp) { // not support type
		return buffer
	}

	imgSrc, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return buffer
	}

	bounds := imgSrc.Bounds()
	sW := bounds.Dx()
	sH := bounds.Dy()

	rgbImg := image.NewRGBA(bounds)

	for y := 0; y < sH; y++ {
		for x := 0; x < sW; x++ {
			r, g, b, a := imgSrc.At(x, y).RGBA()
			rgbImg.SetRGBA(x, y, color.RGBA{
				A: _fixColor(int32(a)),
				R: _fixColor(int32(r>>8) + v),
				G: _fixColor(int32(g>>8) + v),
				B: _fixColor(int32(b>>8) + v),
			})
		}
	}

	return _saveImage(buffer, rgbImg, isPNG, isJPEG, isGIF, isBMP, isWebp)
}

func ContrastImage(buffer []byte, contrastValue *int64, simpleType string) []byte {
	if contrastValue == nil {
		return buffer
	}

	v := int32(*contrastValue)
	if v == 0 {
		return buffer
	}

	isPNG, isJPEG, isBMP, isGIF, isWebp := checkImageType(simpleType)

	if !(isPNG || isJPEG || isGIF || isBMP || isWebp) { // not support type
		return buffer
	}

	imgSrc, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return buffer
	}

	t := defContrastThreshold

	bounds := imgSrc.Bounds()
	sW := bounds.Dx()
	sH := bounds.Dy()

	rgbImg := image.NewRGBA(bounds)

	for y := 0; y < sH; y++ {
		for x := 0; x < sW; x++ {
			r, g, b, a := imgSrc.At(x, y).RGBA()
			rgbImg.SetRGBA(x, y, color.RGBA{
				A: _fixColor(int32(a)),
				R: _fixColor(_computeContrast(int32(r>>8), t, v)),
				G: _fixColor(_computeContrast(int32(g>>8), t, v)),
				B: _fixColor(_computeContrast(int32(b>>8), t, v)),
			})
		}
	}

	return _saveImage(buffer, rgbImg, isPNG, isJPEG, isGIF, isBMP, isWebp)
}

func _computeContrast(col int32, th int32, con int32) int32 {
	return col + (col-th)*con/255
}
