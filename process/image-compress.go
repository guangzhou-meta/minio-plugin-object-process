package process

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ultimate-guitar/go-imagequant"
	"golang.org/x/image/bmp"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
)

func CompressImage(buffer []byte, qualityMin, qualityMax *int64, simpleType string) []byte {
	isPNG, isJPEG, _, isGIF, _ := checkImageType(simpleType)

	if !(isPNG || isJPEG || isGIF) { // not support type
		return buffer
	}

	min := defCompressMin
	max := defCompressMax
	if qualityMin != nil {
		min = int(*qualityMin)
	}
	if qualityMax != nil {
		max = int(*qualityMax)
	}

	if isPNG {
		buffer = compressPng(buffer, min, max)
	} else if isJPEG {
		buffer = compressJpeg(buffer, min, max)
	} else if isGIF {
		buffer = compressGif(buffer, min, max)
	}

	return buffer
}

func compressJpeg(buffer []byte, min, max int) []byte {
	refImg, err := jpeg.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return buffer
	}

	quality := min + int(float64(max-min)*0.5)

	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	err = jpeg.Encode(writer, refImg, &jpeg.Options{Quality: quality})
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	_ = writer.Flush()
	return buf.Bytes()
}

func compressPng(buffer []byte, min, max int) []byte {
	refImg, err := png.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return buffer
	}

	bounds := refImg.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	cmpImg := _compressPng(refImg, width, height, min, max)
	if cmpImg == nil {
		return buffer
	}
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	err = png.Encode(writer, cmpImg)
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	_ = writer.Flush()
	return buf.Bytes()
}

func _compressPng(refImg image.Image, width, height, min, max int) image.Image {
	attr, err := imagequant.NewAttributes()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer attr.Release()
	_ = attr.SetSpeed(10)
	_ = attr.SetQuality(min, max)
	rgba32data := string(imagequant.ImageToRgba32(refImg))
	iqm, err := imagequant.NewImage(attr, rgba32data, width, height, 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer iqm.Release()
	quantizeRes, err := iqm.Quantize(attr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer quantizeRes.Release()
	rgb8data, err := quantizeRes.WriteRemappedImage()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	cmpImg := imagequant.Rgb8PaletteToGoImage(quantizeRes.GetImageWidth(), quantizeRes.GetImageHeight(), rgb8data, quantizeRes.GetPalette())
	if cmpImg == nil {
		fmt.Println(err)
		return nil
	}
	return cmpImg
}

func compressGif(buffer []byte, min, max int) []byte {
	refImg, err := gif.DecodeAll(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	imageArr := refImg.Image
	ch := make(chan goCompressGif)
	for i := range imageArr {
		item := imageArr[i]
		go _compressGif(ch, item, i, min, max)
	}
	resultCount := 0
	for {
		item := <-ch
		imageArr[item.Index] = &item.Result
		resultCount += 1
		if resultCount == len(imageArr) {
			break
		}
	}
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	err = gif.EncodeAll(writer, refImg)
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	_ = writer.Flush()
	return buf.Bytes()
}

func _compressGif(ch chan goCompressGif, item *image.Paletted, index, min, max int) {
	bounds := item.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	cmpImg := _compressPng(item, width, height, min, max)
	if cmpImg != nil {
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, cmpImg, bounds.Min, draw.Over)
		item = palettedImage
	}
	ch <- goCompressGif{
		Index:  index,
		Result: *item,
	}
}

func compressBmp(buffer []byte, min, max int) []byte {
	refImg, err := bmp.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	quality := min + int(float64(max-min)*0.5)

	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	err = jpeg.Encode(writer, refImg, &jpeg.Options{Quality: quality})
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	_ = writer.Flush()

	refImg, err = jpeg.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	buf = bytes.NewBuffer(nil)
	writer = bufio.NewWriter(buf)
	err = bmp.Encode(writer, refImg)
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	_ = writer.Flush()

	return buf.Bytes()
}
