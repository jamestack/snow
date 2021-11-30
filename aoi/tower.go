package aoi

import "fmt"

type Tower struct {
	scene   *Scene // 所属场景
	id      int64  // 场景内的灯塔唯一id
	entitys map[int64]*Entity
}

func (t *Tower) BroadCast(action string) {
	fmt.Println("Tower", t.id, "BroadCast:", action)
}
