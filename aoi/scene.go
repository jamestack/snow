package aoi

type Scene struct {
	towers  map[int64]*Tower
	entitys map[int64]*Entity
	_id     int64
}

// 每个格子的宽度
var TowerSize float64 = 100

// 视野半径
var ViewSize float64 = 10

func getTowerNxNy(posX, posY float64) (nx, ny int32) {
	if posX >= 0 {
		nx = int32(posX / TowerSize)
	} else {
		nx = int32(posX/TowerSize) - 1
	}
	if posY >= 0 {
		ny = int32(posY / TowerSize)
	} else {
		ny = int32(posY/TowerSize) - 1
	}

	return
}

func getTowerIdByNx(nx, ny int32) int64 {
	return int64(nx)<<32 + int64(ny)
}

func getTowerIdByPos(posX, posY float64) int64 {
	nx, ny := getTowerNxNy(posX, posY)
	return getTowerIdByNx(nx, ny)
}

func getTowerPos(nx, ny int32) (posX, posY float64) {
	return float64(nx)*TowerSize + TowerSize/2, float64(ny)*TowerSize + TowerSize/2
}

// 获取感兴趣的灯塔id
func getInterestTowerIds(posX, posY float64) (ids []int64) {
	nx, ny := getTowerNxNy(posX, posY)
	tx, ty := getTowerPos(nx, ny)

	// 计算实体塔内相对位置
	posX -= tx
	posY -= ty

	xR := TowerSize/2 - ViewSize
	xL := ViewSize - TowerSize/2
	yT := xR
	yL := xL

	// 右
	if posX > xR {
		ids = append(ids, getTowerIdByNx(nx+1, ny))
	}

	// 下
	if posY < yL {
		ids = append(ids, getTowerIdByNx(nx, ny-1))
	}

	// 左
	if posX < xL {
		ids = append(ids, getTowerIdByNx(nx-1, ny))
	}

	// 上
	if posY > yT {
		ids = append(ids, getTowerIdByNx(nx, ny+1))
	}

	// 右上
	if posX > xR && posY > yT {
		ids = append(ids, getTowerIdByNx(nx+1, ny+1))
	}

	// 右下
	if posX > xR && posY < yL {
		ids = append(ids, getTowerIdByNx(nx+1, ny-1))
	}

	// 左下
	if posX < xL && posY < yL {
		ids = append(ids, getTowerIdByNx(nx-1, ny-1))
	}

	// 左上
	if posX < xL && posY > yT {
		ids = append(ids, getTowerIdByNx(nx-1, ny-1))
	}

	return ids
}

func (s *Scene) GetTowerById(id int64) *Tower {
	return s.towers[id]
}

// 通过坐标获取灯塔(是否初始化)
func (s *Scene) GetTowerByPos(posX, posY float64) *Tower {
	return s.towers[getTowerIdByPos(posX, posY)]
}

func (s *Scene) GetInterestTowers(posX, posY float64) (list []*Tower) {
	ids := getInterestTowerIds(posX, posY)
	for _, id := range ids {
		list = append(list, s.GetTowerById(id))
	}
	return
}

func (s *Scene) newId() int64 {
	s._id += 1
	return s._id
}

// func (s *Scene) CreateEntity(name string, pos Position) *Entity {
// 	tower := s.GetTowerById()

// 	entity := &Entity{
// 		scene: s,
// 		tower: tower,
// 		id:    s.newId(),
// 		name:  name,
// 		pos:   pos,
// 		Attrs: map[string]interface{}{},
// 	}

// 	tower.entitys[entity.id] = entity

// 	return entity
// }

//
func (s *Scene) GetEntityById(id int64) *Entity {
	return s.entitys[id]
}

// 按照名称查询实体
func (s *Scene) FindEntityByName(name string) (res []*Entity) {
	res = make([]*Entity, 0)
	for _, e := range s.entitys {
		if e.name == name {
			res = append(res, e)
		}
	}
	return
}
