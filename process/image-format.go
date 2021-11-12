package process

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

func FormatImage(buffer []byte, formatType *ImageFormatType, simpleType string) (bufferT []byte, resultType string) {
	resultType = simpleType
	bufferT = buffer
	if formatType == nil {
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

	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)

	switch *formatType {
	case pngType:
		err = png.Encode(writer, imgSrc)
		resultType = "png"
		break
	case jpegType:
		err = jpeg.Encode(writer, imgSrc, nil)
		resultType = "jpeg"
		break
	case gifType:
		err = gif.Encode(writer, imgSrc, nil)
		resultType = "gif"
		break
	case bmpType:
		err = bmp.Encode(writer, imgSrc)
		resultType = "bmp"
		break
	case webpType:
		err = webp.Encode(writer, imgSrc, &webp.Options{Lossless: true, Quality: 100})
		resultType = "webp"
		break
	}
	if err != nil {
		fmt.Println(err)
		return buffer, simpleType
	}
	_ = writer.Flush()

	return buf.Bytes(), resultType
}
