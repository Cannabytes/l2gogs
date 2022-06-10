package models

import (
	"l2gogameserver/gameserver/interfaces"
	"math"
	"sync"
)

type WorldRegion struct {
	TileX         int32
	TileY         int32
	TileZ         int32
	CharsInRegion sync.Map //TODO переделать на мапу с RW мьютексом ци шо ци каво
	NpcInRegion   sync.Map //TODO переделать на мапу с RW мьютексом ци шо ци каво
}

func NewWorldRegion(x, y, z int32) *WorldRegion {
	var newRegion WorldRegion
	newRegion.TileX = x
	newRegion.TileY = y
	newRegion.TileZ = z
	return &newRegion
}

func (w *WorldRegion) AddVisibleChar(character interfaces.CharacterI) {
	w.CharsInRegion.Store(character.ObjectID(), character)
}
func (w *WorldRegion) DeleteVisibleChar(character interfaces.CharacterI) {
	w.CharsInRegion.Delete(character.ObjectID())
}

func (w *WorldRegion) AddVisibleNpc(npc Npc) {
	w.NpcInRegion.Store(npc.ObjId, npc)
}

func (w *WorldRegion) GetNeighbors() []interfaces.WorldRegioner {
	return GetNeighbors(int(w.TileX), int(w.TileY), int(w.TileZ), 1, 1)
}

func (w *WorldRegion) GetCharsInRegion() []interfaces.CharacterI {
	result := make([]interfaces.CharacterI, 0, 64)
	w.CharsInRegion.Range(func(key, value interface{}) bool {
		result = append(result, value.(*Character))
		return true
	})

	return result
}

func (w *WorldRegion) GetNpcInRegion() []interfaces.Npcer {
	result := make([]interfaces.Npcer, 0, 64)
	w.NpcInRegion.Range(func(key, value interface{}) bool {
		result = append(result, value.(Npc))
		return true
	})

	return result
}
func Contains(regions []interfaces.WorldRegioner, region interfaces.WorldRegioner) bool {
	for _, v := range regions {
		if v == region {
			return true
		}
	}
	return false
}
func GetAroundPlayer(c interfaces.Positionable) []interfaces.CharacterI {
	currentRegion := c.GetCurrentRegion()
	if nil == currentRegion {
		return nil
	}
	result := make([]interfaces.CharacterI, 0, 64)

	for _, v := range currentRegion.GetNeighbors() {
		result = append(result, v.GetCharsInRegion()...)
	}
	return result
}
func GetAroundPlayerObjId(c *Character) []int32 {
	currentRegion := c.GetCurrentRegion()
	if nil == currentRegion {
		return nil
	}
	result := make([]int32, 0, 64)

	for _, v := range currentRegion.GetNeighbors() {
		for _, vv := range v.GetCharsInRegion() {
			result = append(result, vv.ObjectID())
		}
	}
	return result
}

func GetAroundPlayersInRadius(c interfaces.CharacterI, radius int32) []*Character {
	currentRegion := c.GetCurrentRegion()
	if nil == currentRegion {
		return nil
	}
	result := make([]*Character, 0, 64)

	sqrad := radius * radius

	for _, v := range currentRegion.GetNeighbors() {
		charInRegion := v.GetCharsInRegion()
		for _, vv := range charInRegion {
			char, ok := vv.(*Character)
			if !ok {
				continue
			}
			if char.ObjectID() == c.ObjectID() {
				continue
			}
			dx := math.Abs(float64(char.Coordinates.X - c.GetX()))
			if dx > float64(radius) {
				continue
			}

			//toDO тут Y должен быть ? проверить на других сборках
			dy := math.Abs(float64(char.Coordinates.X - c.GetX()))
			if dy > float64(radius) {
				continue
			}

			if dx*dx+dy*dy > float64(sqrad) {
				continue
			}

			result = append(result, char)

		}
	}
	return result
}

func GetAroundPlayersObjIdInRadius(c *Character, radius int32) []int32 {
	currentRegion := c.CurrentRegion
	if nil == currentRegion {
		return nil
	}
	result := make([]int32, 0, 64)

	sqrad := radius * radius

	for _, v := range currentRegion.GetNeighbors() {
		charInRegion := v.GetCharsInRegion()
		for _, vv := range charInRegion {
			char, ok := vv.(*Character)
			if !ok {
				continue
			}
			if char.ObjectID() == c.ObjectID() {
				continue
			}
			dx := math.Abs(float64(char.Coordinates.X - c.Coordinates.X))
			if dx > float64(radius) {
				continue
			}

			dy := math.Abs(float64(char.Coordinates.X - c.Coordinates.X))
			if dy > float64(radius) {
				continue
			}

			if dx*dx+dy*dy > float64(sqrad) {
				continue
			}

			result = append(result, char.ObjectID())

		}
	}
	return result
}
