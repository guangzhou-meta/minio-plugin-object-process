package process

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

func SharpenImage(buffer []byte, sharpenValue *int64, simpleType string) []byte {
	if sharpenValue == nil {
		return buffer
	}

	v := *sharpenValue
	if v < 1 {
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
	width := bounds.Dx()
	height := bounds.Dy()

	imgTmp := image.NewRGBA(bounds)

	_laplacianSharpen(imgSrc, imgTmp, width, height, int32(v))

	return _saveImage(buffer, imgTmp, isPNG, isJPEG, isGIF, isBMP, isWebp)
}

func _parseColor(src image.Image, x int, y int) (int, int, int, int) {
	r0, g0, b0, a0 := src.At(x, y).RGBA()
	r := int(r0 >> 8)
	g := int(g0 >> 8)
	b := int(b0 >> 8)
	a := int(a0 >> 8)
	return r, g, b, a
}

func _laplacianSharpen(src image.Image, tmp *image.RGBA, width, height int, sharpen int32) {
	maxHeight := height - 1
	maxWidth := width - 1
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := _parseColor(src, x, y)
			if y > 0 && y < maxHeight && x > 0 && x < maxWidth {
				var rChannel []int
				var gChannel []int
				var bChannel []int
				for x0 := -1; x0 < 2; x0++ {
					for y0 := -1; y0 < 2; y0++ {
						r0, g0, b0, _ := _parseColor(src, x0+x, y0+y)
						rChannel = append(rChannel, r0)
						gChannel = append(gChannel, g0)
						bChannel = append(bChannel, b0)
					}
				}
				r = __laplacianSharpen(rChannel, sharpen, 15)
				g = __laplacianSharpen(gChannel, sharpen, 15)
				b = __laplacianSharpen(bChannel, sharpen, 15)
			}
			col := color.RGBA{
				A: uint8(a),
				R: _fixColor(int32(r)),
				G: _fixColor(int32(g)),
				B: _fixColor(int32(b)),
			}
			tmp.SetRGBA(x, y, col)
		}
	}
}

func __laplacianSharpen(channel []int, sharpen, threshold int32) int {
	gaussianMat := []int{
		1, 2, 1, 2, 4, 2, 1, 2, 1,
	}

	sum := 0
	for i := range channel {
		c := channel[i]
		sum += c * gaussianMat[i]
	}

	src := channel[4]
	blur := uint8(sum >> 4)
	texture := int32(_fixColor(int32(src - int(blur))))
	if texture > threshold {
		detail := (texture * sharpen) >> 5
		return int(_fixColor(int32(src) + detail))
	}
	return src
}
