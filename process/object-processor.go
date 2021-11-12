package process

import (
	"bytes"
	"image/color"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// ObjectProcessInfo
// -----------
// Object process params
type ObjectProcessInfo struct {
	Actions []ObjectProcess

	// Use on image circle or rounded-corner
	LastImageFormatType *ImageFormatType
}

type ObjectProcess struct {
	Action ObjectProcessAction

	ImageQualityMin *int64
	ImageQualityMax *int64

	ImageHeight     *int64
	ImageWidth      *int64
	ImageColor      *color.RGBA
	ImageResizeMode *ImageResizeMode

	ImagePositionX *int64
	ImagePositionY *int64

	ImageFormatType *ImageFormatType

	ImageRadius *int64

	ImageValue *int64
}

func (info *ObjectProcessInfo) IsProcessImage() bool {
	actions := info.Actions
	length := len(actions)
	for i := range actions {
		actionL := actions[i].Action
		actionR := actions[length-i-1].Action
		isImageProcess := isProcessImage(actionL) || isProcessImage(actionR)
		if isImageProcess {
			return isImageProcess
		}
	}
	return false
}

func isProcessImage(action ObjectProcessAction) bool {
	return action == ImageCropAction ||
		action == ImageResizeAction ||
		action == ImageCompressAction ||
		action == ImageFormatAction ||
		action == ImageCircleCropAction ||
		action == ImageRoundedCornersCropAction ||
		action == ImageBrightAction ||
		action == ImageContrastAction ||
		action == ImageRotateAction ||
		action == ImageSharpenAction
}

// ObjectProcessAction
// -----------
// Object process action
type ObjectProcessAction int

const (
	ImageCropAction ObjectProcessAction = iota
	ImageResizeAction
	ImageCompressAction
	ImageFormatAction
	ImageCircleCropAction
	ImageRoundedCornersCropAction
	ImageBrightAction
	ImageContrastAction
	ImageRotateAction
	ImageSharpenAction
)

func parseObjectProcessInfo(processQuery string) ObjectProcessInfo {
	fileProcess := strings.TrimSpace(processQuery)
	objectProcessInfo := ObjectProcessInfo{}
	if fileProcess != "" {
		processParams := strings.Split(fileProcess, "/")
		switch processParams[0] {
		case "image":
			parseImageProcessInfo(processParams, &objectProcessInfo)
			break
		}
	}
	return objectProcessInfo
}

func parseImageProcessInfo(processParams []string, objectProcessInfo *ObjectProcessInfo) {
	for i := 1; i < len(processParams); i++ {
		param := strings.TrimSpace(processParams[i])
		if param != "" {
			params := strings.Split(param, ",")
			switch strings.TrimSpace(params[0]) {
			case "resize":
				info := &ObjectProcess{
					Action: ImageResizeAction,
				}
				parseResizeImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "quality":
				info := &ObjectProcess{
					Action: ImageCompressAction,
				}
				parseCompressImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "crop":
				info := &ObjectProcess{
					Action: ImageCropAction,
				}
				parseCropImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "format":
				info := &ObjectProcess{
					Action: ImageFormatAction,
				}
				parseFormatImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				(*objectProcessInfo).LastImageFormatType = info.ImageFormatType
				break
			case "circle":
				info := &ObjectProcess{
					Action: ImageCircleCropAction,
				}
				parseCircleCropImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "rounded-corners":
				info := &ObjectProcess{
					Action: ImageRoundedCornersCropAction,
				}
				parseCircleCropImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "bright":
				info := &ObjectProcess{
					Action: ImageBrightAction,
				}
				parseBrightImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "contrast":
				info := &ObjectProcess{
					Action: ImageContrastAction,
				}
				parseBrightImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "rotate":
				info := &ObjectProcess{
					Action: ImageRotateAction,
				}
				parseRotateImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			case "sharpen":
				info := &ObjectProcess{
					Action: ImageSharpenAction,
				}
				parseSharpenImageInfo(params, info)
				(*objectProcessInfo).Actions = append((*objectProcessInfo).Actions, *info)
				break
			}
		}
	}
}

func parseResizeImageInfo(params []string, info *ObjectProcess) {
	parseProcessParams(params, func(name string, value *string) {
		switch name {
		case "w":
			i := convImageProcessParamToInt64(value, 4096, 1, 4096)
			(*info).ImageWidth = &i
			break
		case "h":
			i := convImageProcessParamToInt64(value, 4096, 1, 4096)
			(*info).ImageHeight = &i
			break
		case "m":
			var ImageResizeMode ImageResizeMode
			switch *value {
			case "lfit":
				ImageResizeMode = lfit
				break
			case "mfit":
				ImageResizeMode = mfit
				break
			case "fill":
				ImageResizeMode = fill
				break
			case "pad":
				ImageResizeMode = pad
				break
			case "fixed":
				ImageResizeMode = fixed
				break
			}
			(*info).ImageResizeMode = &ImageResizeMode
			break
		case "color":
			c := hexToRGBA(*value)
			(*info).ImageColor = &c
			break
		}
	})
}

func parseCompressImageInfo(params []string, info *ObjectProcess) {
	parseProcessParams(params, func(name string, value *string) {
		switch name {
		case "q":
			v := *value
			vs := strings.Split(v, "-")
			min := convImageProcessParamToInt64(&(vs[0]), 100, 1, 100)
			max := min
			if len(vs) > 1 {
				max = convImageProcessParamToInt64(&(vs[1]), 100, 1, 100)
				if max < min {
					max, min = min, max
				}
			} else {
				if min > 40 && min <= 100 {
					min = min - 4
				} else if min <= 40 {
					max = min + 4
				}
			}
			(*info).ImageQualityMin = &min
			(*info).ImageQualityMax = &max
			break
		}
	})
}

func parseCropImageInfo(params []string, info *ObjectProcess) {
	parseProcessParams(params, func(name string, value *string) {
		switch name {
		case "w":
			i := convImageProcessParamToInt64(value, 4096, 1, 4096)
			(*info).ImageWidth = &i
			break
		case "h":
			i := convImageProcessParamToInt64(value, 4096, 1, 4096)
			(*info).ImageHeight = &i
			break
		case "x":
			i := convImageProcessParamToInt64(value, 4096, 0, 4096)
			(*info).ImagePositionX = &i
			break
		case "y":
			i := convImageProcessParamToInt64(value, 4096, 0, 4096)
			(*info).ImagePositionY = &i
			break
		}
	})
}

func parseFormatImageInfo(params []string, info *ObjectProcess) {
	if len(params) > 1 {
		var iType ImageFormatType
		switch params[1] {
		case "jpg":
		case "jpeg":
			iType = jpegType
			break
		case "png":
			iType = pngType
			break
		case "gif":
			iType = gifType
			break
		case "bmp":
			iType = bmpType
			break
		case "webp":
			iType = webpType
			break
		}
		(*info).ImageFormatType = &iType
	}
}

func parseCircleCropImageInfo(params []string, info *ObjectProcess) {
	parseProcessParams(params, func(name string, value *string) {
		switch name {
		case "r":
			i := convImageProcessParamToInt64(value, 0, 1, 4096)
			(*info).ImageRadius = &i
			break
		}
	})
}

func parseBrightImageInfo(params []string, info *ObjectProcess) {
	if len(params) > 1 {
		n := params[1]
		v := convImageProcessParamToInt64(&n, 0, -100, 100)
		if v != 0 {
			info.ImageValue = &v
		}
	}
}

func parseRotateImageInfo(params []string, info *ObjectProcess) {
	if len(params) > 1 {
		n := params[1]
		v := convImageProcessParamToInt64(&n, 0, 0, 360)
		if v != 0 {
			info.ImageValue = &v
		}
	}
}

func parseSharpenImageInfo(params []string, info *ObjectProcess) {
	if len(params) > 1 {
		n := params[1]
		v := convImageProcessParamToInt64(&n, 50, 50, 399)
		if v != 0 {
			info.ImageValue = &v
		}
	}
}

func convImageProcessParamToInt64(value *string, defaultValue, min, max int64) int64 {
	v := *value
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		i = defaultValue
	}
	i = int64(math.Min(math.Max(float64(i), float64(min)), float64(max)))
	return i
}

type objectParamsProcessHandler func(name string, value *string)

func parseProcessParams(params []string, processHandler objectParamsProcessHandler) {
	for i := 1; i < len(params); i++ {
		p := strings.Split(strings.TrimSpace(params[i]), "_")
		if len(p) > 1 {
			processHandler(p[0], &(p[1]))
		} else {
			processHandler(p[0], nil)
		}
	}
}

func hexToRGBA(col string) color.RGBA {
	col = strings.ReplaceAll(col, "#", "")
	c := color.RGBA{
		A: 255,
		R: 0,
		G: 0,
		B: 0,
	}
	l := len(col)
	if l == 6 || l == 8 {
		if len(col) == 8 { // 带alpha值
			alpha := col[0:2] // 截取alpha值
			a, _ := strconv.ParseUint(alpha, 16, 8)
			c.A = uint8(a)
			col = col[2:]
		} else {
			c.A = 255
		}
		red := col[0:2] // 红色
		r, _ := strconv.ParseUint(red, 16, 8)
		c.R = uint8(r)
		green := col[2:4] // 绿色
		g, _ := strconv.ParseUint(green, 16, 8)
		c.G = uint8(g)
		blue := col[4:6] // 蓝色
		b, _ := strconv.ParseUint(blue, 16, 8)
		c.B = uint8(b)
	}
	return c
}

func ProcessObject(objectReader io.Reader, processQuery string) (resultReader io.Reader, contentLength string, contentType string) {
	contentLength = ""
	contentType = ""
	resultReader = objectReader

	if processQuery == "" || objectReader == nil {
		return
	}

	processInfo := parseObjectProcessInfo(processQuery)

	if !processInfo.IsProcessImage() {
		return
	}

	buffer, err := io.ReadAll(objectReader)
	if err == nil { // Process object
		objectType := checkObjectType(buffer)
		switch {
		case objectType.IsImage: // Process image
			ProcessImage(processInfo, objectType, &buffer, &contentType)
			break
		}
	}
	contentLength = strconv.FormatInt(int64(len(buffer)), 10)
	resultReader = bytes.NewReader(buffer)
	return
}

func checkObjectType(buf []byte) ObjectTypeInfo {
	l := int(math.Min(512, float64(len(buf))))
	headerBuf := buf[0:l]
	headerStr := http.DetectContentType(headerBuf)

	arr := strings.Split(headerStr, "/")

	t := arr[0]
	result := ObjectTypeInfo{
		ContentType: headerStr,
		TypeGroup:   t,
		SimpleType:  arr[1],
		IsImage:     t == "image",
	}
	return result
}
