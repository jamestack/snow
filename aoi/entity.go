package aoi

import (
	"fmt"
)

type Entity struct {
	scene *Scene // 所属场景
	tower *Tower // 所属灯塔
	id    int64  // 实体唯一id
	name  string
	pos   Position               // 实体的位置
	Attrs map[string]interface{} // 属性列表
}

type Position struct {
	x, y float64
}

// 设置位置
func (e *Entity) SetPosition(x, y float64) {
	e.pos.x = x
	e.pos.y = y
}

// 设置属性
func (e *Entity) SetAttr(key string, val interface{}) {
	e.Attrs[key] = val
}

// 获取属性
func (e *Entity) GetAttr(key string) (val interface{}) {
	return e.Attrs[key]
}

// 广播行为
func (e *Entity) DoBroadAction(event string) {
	fmt.Println("DoBroadAction =>", event)
}

// 移动
func (e *Entity) MoveTo(targetPosition *Position) {
	e.BroadCast(fmt.Sprintln(e.name, "Move To:", targetPosition.x, targetPosition.y))
}

func (e *Entity) SubBroadcast() {

}

// 广播给周围实体
func (e *Entity) BroadCast(event string) {

}
