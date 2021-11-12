package process

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"regexp"
)

// ImageResizeMode
// -----------
// Image resize mode
type ImageResizeMode int

const (
	lfit ImageResizeMode = iota
	mfit
	fill
	pad
	fixed
)

// ImageFormatType
// -----------
// Image format type
type ImageFormatType int

const (
	jpegType ImageFormatType = iota
	pngType
	gifType
	bmpType
	webpType
)

const (
	defCompressMin             = 40
	defCompressMax             = 90
	defContrastThreshold int32 = 128
)

// goCompressGif
type goCompressGif struct {
	Index  int
	Result image.Paletted
}

func checkImageType(simpleType string) (bool, bool, bool, bool, bool) {
	v := []byte(simpleType)
	isPNG := regexp.MustCompile(`(?i)png`).Match(v)
	isJPEG := regexp.MustCompile(`(?i)jpeg`).Match(v)
	isBMP := regexp.MustCompile(`(?i)bmp`).Match(v)
	isGIF := regexp.MustCompile(`(?i)gif`).Match(v)
	isWebp := regexp.MustCompile(`(?i)webp`).Match(v)
	return isPNG, isJPEG, isBMP, isGIF, isWebp
}

// roundedCorner
type roundedCorner struct {
	Width  int
	Height int
	Radius int
}

func _saveImage(buffer []byte, rgbImg image.Image, isPNG bool, isJPEG bool, isGIF bool, isBMP bool, isWebp bool) []byte {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	var err error
	switch {
	case isPNG:
		err = png.Encode(writer, rgbImg)
		break
	case isJPEG:
		err = jpeg.Encode(writer, rgbImg, nil)
		break
	case isGIF:
		err = gif.Encode(writer, rgbImg, nil)
		break
	case isBMP:
		err = bmp.Encode(writer, rgbImg)
		break
	case isWebp:
		err = webp.Encode(writer, rgbImg, &webp.Options{Lossless: true, Quality: 100})
		break
	}
	if err != nil {
		fmt.Println(err)
		return buffer
	}
	_ = writer.Flush()

	return buf.Bytes()
}

func ProcessImage(processInfo ObjectProcessInfo, objectType ObjectTypeInfo, buffer *[]byte, contentType *string) {
	bf := *buffer
	ct := *contentType
	actions := processInfo.Actions
	for i := range actions {
		simpleType := objectType.SimpleType
		action := actions[i]
		switch action.Action {
		case ImageCropAction: // crop
			bf = CropImage(bf, action.ImageWidth, action.ImageHeight, action.ImagePositionX, action.ImagePositionY, simpleType)
			break
		case ImageResizeAction: // resize
			bf = ResizeImage(bf, action.ImageWidth, action.ImageHeight, action.ImageResizeMode, action.ImageColor, simpleType)
			break
		case ImageCompressAction: // compress
			bf = CompressImage(bf, action.ImageQualityMin, action.ImageQualityMax, simpleType)
			break
		case ImageFormatAction: // format
			buf, resultType := FormatImage(bf, action.ImageFormatType, simpleType)
			bf = buf
			objectType.SimpleType = resultType
			ct = fmt.Sprintf("image/%s", resultType)
			break
		case ImageCircleCropAction: // circle crop
			buf, resultType := CircleCropImage(bf, action.ImageRadius, processInfo.LastImageFormatType, simpleType)
			bf = buf
			objectType.SimpleType = resultType
			ct = fmt.Sprintf("image/%s", resultType)
			break
		case ImageRoundedCornersCropAction: // rounded-corner crop
			buf, resultType := RoundedCornerCropImage(bf, action.ImageRadius, processInfo.LastImageFormatType, simpleType)
			bf = buf
			objectType.SimpleType = resultType
			ct = fmt.Sprintf("image/%s", resultType)
			break
		case ImageBrightAction: // bright
			bf = BrightImage(bf, action.ImageValue, simpleType)
			break
		case ImageContrastAction: // contrast
			bf = ContrastImage(bf, action.ImageValue, simpleType)
			break
		case ImageRotateAction: // rotate
			bf = RotateImage(bf, action.ImageValue, simpleType)
			break
		case ImageSharpenAction: // sharpen
			bf = SharpenImage(bf, action.ImageValue, simpleType)
			break
		}
	}
	*buffer = bf
	*contentType = ct
}
