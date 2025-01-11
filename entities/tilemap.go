package entities

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

type TilemapJSON struct {
	Layers   []TilemapLayerJSON `json:"layers"`
	Tilesets []map[string]any   `json:"tilesets"`
}

func (t *TilemapJSON) GenTilesets(mapPath string) ([]Tileset, error) {
	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		tilesetPath := resolveTilesetPath(mapPath, tilesetData["source"].(string))
		tileset, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}

		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}

func resolveTilesetPath(basePath, relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(filepath.Dir(basePath), relativePath)
}
