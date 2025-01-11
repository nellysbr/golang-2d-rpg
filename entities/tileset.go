package entities

import (
	"encoding/json"
	"golang-2d-rpg/utils"
	"image"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid

	srcX := id % 22
	srcY := id / 22

	srcX *= 16
	srcY *= 16

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+16, srcY+16,
		),
	).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	if id < 0 || id >= len(d.imgs) {
		log.Printf("Warning: Tile ID %d (original: %d) is out of range for DynTileset", id, id+d.gid)
		return nil // Ou retorne uma imagem padrão/vazia
	}
	return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {
	log.Printf("Creating new tileset from path: %s with gid: %d", path, gid)
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynTilesetJSON.Tiles {
			tileJSONPath := utils.ResolveTilesetPath(path, tileJSON.Path)
			log.Printf("Loading tile image from: %s", tileJSONPath)
			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				log.Printf("Error loading tile image: %v", err)
				return nil, err
			}
			dynTileset.imgs = append(dynTileset.imgs, img)
		}
		log.Printf("Created DynTileset with %d images", len(dynTileset.imgs))
		return &dynTileset, nil
	}

	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	tileJSONPath := utils.ResolveTilesetPath(path, uniformTilesetJSON.Path)
	log.Printf("Loading tileset image from: %s", tileJSONPath)
	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		log.Printf("Error loading tileset image: %v", err)
		return nil, err
	}
	uniformTileset.img = img
	uniformTileset.gid = gid

	log.Printf("Created UniformTileset with image: %s", tileJSONPath)
	return &uniformTileset, nil
}
