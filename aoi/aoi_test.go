package aoi

import (
	"fmt"
	"testing"
)

func TestAoi(t *testing.T) {
	scene := &Scene{
		towers:  make(map[int64]*Tower),
		entitys: make(map[int64]*Entity),
	}

	james := scene.CreateEntity("james", &Position{x: 99, y: 99})
	jack := scene.CreateEntity("jack", &Position{x: 101, y: 99})

	fmt.Println("james:", james.tower.id, "jack:", jack.tower.id)

	_ = jack
	james.MoveTo(&Position{x: 9, y: 9})
	// jack.MoveTo(&Position{x: 12, y: 5})
}
