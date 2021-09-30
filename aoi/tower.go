package aoi

import "fmt"

type Tower struct {
	scene   *Scene // 所属场景
	id      int64  // 场景内的灯塔唯一id
	rect    Rect
	entitys map[int64]*Entity
}

type Rect struct {
	pos    Position // 矩形中心坐标
	width  float64  // 矩形长度
	height float64  // 矩形宽度
}

// type Circle struct {
// 	point  Position // 圆心坐标
// 	radius float64  // 圆的半径
// }

// // 判断坐标系中,矩形与圆是否相交
// func isRectTouchCircle(rect *Rect, circle *Circle) bool {
// 	// 1. 偏移矩形变换坐标到(0,0), 并将圆翻转到第一象限
// 	Cx := math.Abs(circle.point.x - rect.point.x)
// 	Cy := math.Abs(circle.point.y - rect.point.y)
// 	Rw := rect.width / 2
// 	Rh := rect.height / 2
// 	// 2. 判断
// 	if Cy <= Rh {
// 		return Cx < Rw+circle.radius
// 	} else if Cx <= Rw {
// 		return Cy < Rh+circle.radius
// 	} else {
// 		return math.Sqrt(math.Pow(Rw, 2)+math.Pow(Rh, 2))+circle.radius < math.Sqrt(math.Pow(Cx, 2)+math.Pow(Cy, 2))
// 	}
// }

func (t *Tower) BroadCast(action string) {
	fmt.Println("Tower", t.id, "BroadCast:", action)
}
