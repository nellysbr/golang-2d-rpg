package entities

import (
	"encoding/json"
	"os"
)

type TimemapLayerJson struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJson struct {
	Layers []TimemapLayerJson `json:"layers"`
}

func NewTilemapJson(filepath string) (*TilemapJson, error) {

	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJson TilemapJson
	err = json.Unmarshal(contents, &tilemapJson)

	if err != nil {
		return nil, err
	}

	return &tilemapJson, nil
}
