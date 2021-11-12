package process

import (
	"bufio"
	"bytes"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
)

func CropImage(buffer []byte, cropWidth, cropHeight, cropX, cropY *int64, simpleType string) []byte {
	if cropWidth == nil && cropHeight == nil {
		return buffer
	}
	var w = 0
	var h = 0
	if cropWidth != nil {
		w = int(*cropWidth)
		if w < 1 {
			w = 0
		}
	}
	if cropHeight != nil {
		h = int(*cropHeight)
		if h < 1 {
			h = 0
		}
	}

	if w == 0 && h == 0 {
		return buffer
	}

	var x = 0
	var y = 0
	if cropX != nil {
		x = int(*cropX)
	}
	if cropY != nil {
		y = int(*cropY)
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

	cropImg := _cropImage(imgSrc, w, h, x, y)

	if cropImg == nil {
		return buffer
	}

	return _saveImage(buffer, cropImg, isPNG, isJPEG, isGIF, isBMP, isWebp)
}

func _cropImage(imgSrc image.Image, w int, h int, x int, y int) image.Image {
	var cropImg image.Image
	if rgbImg, ok := imgSrc.(*image.YCbCr); ok {
		cropImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := imgSrc.(*image.RGBA); ok {
		cropImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := imgSrc.(*image.NRGBA); ok {
		cropImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	}
	return cropImg
}

func CircleCropImage(buffer []byte, cropRadius *int64, formatType *ImageFormatType, simpleType string) (bufferT []byte, resultType string) {
	resultType = simpleType
	bufferT = buffer
	if cropRadius == nil {
		return
	}
	radius := int(math.Max(0, float64(*cropRadius)))
	if radius == 0 {
		return
	}

	isPNG, isJPEG, isBMP, isGIF, isWebp := checkImageType(simpleType)

	if !(isPNG || isJPEG || isGIF || isBMP || isWebp) { // not support type
		return
	}

	imgSrc, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return
	}
	bounds := imgSrc.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// fixed radius
	radius = int(math.Min(float64(radius), math.Min(float64(width), float64(height))*0.5))

	w := radius * 2
	cropImg := _cropImage(imgSrc, w, w, int(float64(width)*0.5)-radius, int(float64(height)*0.5)-radius)
	if cropImg == nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	err = png.Encode(writer, cropImg)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = writer.Flush()
	cropImg, err = png.Decode(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		fmt.Println(err)
		return
	}

	isFormatJpeg := formatType == nil || *formatType != pngType

	render := image.NewRGBA(image.Rect(0, 0, w, w))
	if isFormatJpeg {
		draw.Draw(render, render.Bounds(), &image.Uniform{
			C: color.RGBA{
				A: 255,
				R: 255,
				G: 255,
				B: 255,
			},
		}, image.Point{}, draw.Src)
	}
	draw.DrawMask(render, cropImg.Bounds(), cropImg, image.Point{}, &roundedCorner{
		Width: w, Height: w, Radius: radius,
	}, image.Point{}, draw.Over)

	buf = bytes.NewBuffer(nil)
	writer = bufio.NewWriter(buf)
	if isFormatJpeg {
		err = jpeg.Encode(writer, render, nil)
		resultType = "jpeg"
	} else {
		err = png.Encode(writer, render)
		resultType = "png"
	}
	if err != nil {
		fmt.Println(err)
		return buffer, simpleType
	}
	_ = writer.Flush()

	return buf.Bytes(), resultType
}

func RoundedCornerCropImage(buffer []byte, cropRadius *int64, formatType *ImageFormatType, simpleType string) (bufferT []byte, resultType string) {
	resultType = simpleType
	bufferT = buffer
	if cropRadius == nil {
		return
	}
	radius := int(math.Max(0, float64(*cropRadius)))
	if radius == 0 {
		return
	}

	isPNG, isJPEG, isBMP, isGIF, isWebp := checkImageType(simpleType)

	if !(isPNG || isJPEG || isGIF || isBMP || isWebp) { // not support type
		return
	}

	imgSrc, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return
	}
	bounds := imgSrc.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// fixed radius
	radius = int(math.Min(float64(radius), math.Min(float64(width), float64(height))*0.5))

	isFormatJpeg := formatType == nil || *formatType != pngType

	render := image.NewRGBA(bounds)
	if isFormatJpeg {
		draw.Draw(render, render.Bounds(), &image.Uniform{
			C: color.RGBA{
				A: 255,
				R: 255,
				G: 255,
				B: 255,
			},
		}, image.Point{}, draw.Src)
	}
	draw.DrawMask(render, imgSrc.Bounds(), imgSrc, image.Point{}, &roundedCorner{
		Width: width, Height: height, Radius: radius,
	}, image.Point{}, draw.Over)

	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	if isFormatJpeg {
		err = jpeg.Encode(writer, render, nil)
		resultType = "jpeg"
	} else {
		err = png.Encode(writer, render)
		resultType = "png"
	}
	if err != nil {
		fmt.Println(err)
		return buffer, simpleType
	}
	_ = writer.Flush()

	return buf.Bytes(), resultType
}
