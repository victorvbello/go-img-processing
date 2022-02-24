package image

import "image"

func GetCenterPoint(rect image.Rectangle) image.Point {
	return image.Point{(rect.Min.X + rect.Max.X) / 2, (rect.Min.Y + rect.Max.Y) / 2}
}
