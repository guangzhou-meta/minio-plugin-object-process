package process

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
)

func RotateImage(buffer []byte, rotateValue *int64, simpleType string) []byte {
	if rotateValue == nil {
		return buffer
	}

	v := float64(*rotateValue)
	if int(v)%360 == 0 {
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

	hW := float64(sW) / 2
	hH := float64(sH) / 2

	sin := math.Sin(v * math.Pi / 180)
	cos := math.Cos(v * math.Pi / 180)

	ltX, ltY := _computeRotatePosition(-hW, hH, sin, cos)
	rtX, rtY := _computeRotatePosition(hW, hH, sin, cos)
	lbX, lbY := _computeRotatePosition(-hW, -hH, sin, cos)
	rbX, rbY := _computeRotatePosition(hW, -hH, sin, cos)

	maxWidth := int(math.Max(math.Abs(rbX-ltX), math.Abs(rtX-lbX)))
	maxHeight := int(math.Max(math.Abs(rbY-ltY), math.Abs(rtY-lbY)))

	imgTmp := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))

	halfWidth := float64(maxWidth) / 2
	halfHeight := float64(maxHeight) / 2
	sinR := math.Sin((360 - v) * math.Pi / 180)
	cosR := math.Cos((360 - v) * math.Pi / 180)
	for y := 0; y < maxHeight; y++ {
		for x := 0; x < maxWidth; x++ {
			tX := int((float64(x)-halfWidth)*cosR + (-float64(y)+halfHeight)*sinR)
			tY := int(-(float64(x)-halfWidth)*sinR + (-float64(y)+halfHeight)*cosR)
			tXF := float64(tX)
			tYF := float64(tY)
			colour := color.RGBA{
				A: 255,
				R: 255,
				G: 255,
				B: 255,
			}
			if !(tXF > hW || tXF < -hW || tYF > hH || tYF < -hH) {
				tXN := int(tXF + hW)
				tYN := int(math.Abs(tYF - hH))
				r, g, b, a := imgSrc.At(tXN, tYN).RGBA()
				colour.A = uint8(a)
				colour.R = uint8(r >> 8)
				colour.G = uint8(g >> 8)
				colour.B = uint8(b >> 8)
			}
			imgTmp.SetRGBA(x, y, colour)
		}
	}

	return _saveImage(buffer, imgTmp, isPNG, isJPEG, isGIF, isBMP, isWebp)
}

func _computeRotatePosition(x, y, sin, cos float64) (float64, float64) {
	x0 := x*cos + y*sin
	y0 := -x*sin + y*cos
	return x0, y0
}
