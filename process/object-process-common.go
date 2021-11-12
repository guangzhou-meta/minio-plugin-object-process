package process

import (
	"image"
	"image/color"
	"math"
)

// ObjectTypeInfo
// -----------
// Object content-type information
type ObjectTypeInfo struct {
	ContentType string
	TypeGroup   string
	SimpleType  string
	IsImage     bool
}

func (rc *roundedCorner) ColorModel() color.Model {
	return color.AlphaModel
}

func (rc *roundedCorner) Bounds() image.Rectangle {
	return image.Rect(0, 0, rc.Width, rc.Height)
}

func (rc *roundedCorner) At(x, y int) color.Color {
	r := rc.Radius
	if r == 0 {
		return color.Alpha{A: 255}
	}
	rF := float64(r * r)
	w := rc.Width - 1
	h := rc.Height - 1
	if x < r && y < r && (rF < math.Pow(float64(r-x), 2)+math.Pow(float64(r-y), 2)) { // top left
		return color.Alpha{A: 0}
	} else if x > (w-r) && y < r && (rF < math.Pow(float64(w-r-x), 2)+math.Pow(float64(r-y), 2)) { // top right
		return color.Alpha{A: 0}
	} else if x < r && y > (h-r) && (rF < math.Pow(float64(r-x), 2)+math.Pow(float64(h-r-y), 2)) { // bottom left
		return color.Alpha{A: 0}
	} else if x > (w-r) && y > (h-r) && (rF < math.Pow(float64(w-r-x), 2)+math.Pow(float64(h-r-y), 2)) { // bottom right
		return color.Alpha{A: 0}
	}
	return color.Alpha{A: 255}
}

func _fixColor(c int32) uint8 {
	if c > 255 {
		return 255
	}
	if c < 0 {
		return 0
	}
	return uint8(c)
}
