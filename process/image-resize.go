package process

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
)

func ResizeImage(buffer []byte, resizeWidth, resizeHeight *int64, resizeMode *ImageResizeMode, padColor *color.RGBA, simpleType string) []byte {
	if resizeWidth == nil && resizeHeight == nil {
		return buffer
	}
	var w = 0
	var h = 0
	if resizeWidth != nil {
		w = int(*resizeWidth)
		if w < 1 {
			w = 0
		}
	}
	if resizeHeight != nil {
		h = int(*resizeHeight)
		if h < 1 {
			h = 0
		}
	}

	if w == 0 && h == 0 {
		return buffer
	}

	wF := float64(w)
	hF := float64(h)

	var m = lfit
	if resizeMode != nil {
		m = *resizeMode
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
	sWF := float64(sW)
	sHF := float64(sH)
	wR := wF / sWF
	hR := hF / sHF

	// Resize
	switch m {
	case fixed:
		imgSrc = resize.Resize(uint(w), uint(h), imgSrc, resize.NearestNeighbor)
		break
	case lfit:
		ratio := math.Min(wR, hR)
		imgSrc = resize.Resize(uint(sWF*ratio), uint(sHF*ratio), imgSrc, resize.NearestNeighbor)
		break
	case pad:
		ratio := math.Min(wR, hR)
		rW := uint(sWF * ratio)
		rH := uint(sHF * ratio)
		imgSrc = resize.Resize(rW, rH, imgSrc, resize.NearestNeighbor)
		imgSrc = drawWrap(imgSrc, w, h, int(math.Abs((wF-float64(rW))*0.5)), int(math.Abs((hF-float64(rH))*0.5)), padColor)
		break
	case mfit:
		ratio := math.Max(wR, hR)
		rH := uint(sHF * ratio)
		imgSrc = resize.Resize(uint(sWF*ratio), rH, imgSrc, resize.NearestNeighbor)
		break
	case fill:
		ratio := math.Max(wR, hR)
		rW := uint(sWF * ratio)
		rH := uint(sHF * ratio)
		imgSrc = resize.Resize(rW, rH, imgSrc, resize.NearestNeighbor)
		imgSrc = _cropImage(imgSrc, w, h, int(math.Abs((wF-float64(rW))*0.5)), int(math.Abs((hF-float64(rH))*0.5)))
		break
	}

	//buf := bytes.NewBuffer(nil)
	//writer := bufio.NewWriter(buf)
	//switch {
	//case isPNG:
	//	err = png.Encode(writer, imgSrc)
	//	break
	//case isJPEG:
	//	err = jpeg.Encode(writer, imgSrc, nil)
	//	break
	//case isGIF:
	//	err = gif.Encode(writer, imgSrc, nil)
	//	break
	//case isBMP:
	//	err = bmp.Encode(writer, imgSrc)
	//	break
	//case isWebp:
	//	err = webp.Encode(writer, imgSrc, &webp.Options{Lossless: true, Quality: 100})
	//	break
	//}
	//if err != nil {
	//	fmt.Println(err)
	//	return buffer
	//}
	//_ = writer.Flush()
	//
	//return buf.Bytes()
	return _saveImage(buffer, imgSrc, isPNG, isJPEG, isGIF, isBMP, isWebp)
}

func drawWrap(src image.Image, w int, h int, x int, y int, padColor *color.RGBA) image.Image {
	wrapImg := image.NewRGBA(image.Rect(0, 0, w, h))
	var col color.RGBA
	if padColor == nil {
		col = color.RGBA{
			A: 255,
			R: 0,
			G: 0,
			B: 0,
		}
	} else {
		col = *padColor
	}
	draw.Draw(wrapImg, wrapImg.Bounds(), &image.Uniform{C: col}, image.Point{}, draw.Src)
	bounds := src.Bounds()
	p := image.Pt(x, y)
	draw.Draw(wrapImg, bounds.Add(p), src, bounds.Min, draw.Src)
	return wrapImg
}
