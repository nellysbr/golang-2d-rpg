package main

import (
	"golang-2d-rpg/entities"
	"golang-2d-rpg/utils"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	experience  []*entities.Coin
	tilemapJson *entities.TilemapJSON
	tilesets    []entities.Tileset

	tilemapImage *ebiten.Image
	cam          *utils.Camera
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y += 2
	}

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if g.player.X < sprite.X {
				sprite.X -= 0.5
			} else {
				sprite.X += 0.5
			}

			if g.player.Y < sprite.Y {
				sprite.Y -= 0.5
			} else {
				sprite.Y += 0.5
			}
		}
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 640, 480)
	g.cam.Constrain(float64(g.tilemapJson.Layers[0].Width)*16.0, float64(g.tilemapJson.Layers[0].Height)*16, 320, 640)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	opts := ebiten.DrawImageOptions{}

	for layerIndex, layer := range g.tilemapJson.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			if layerIndex >= len(g.tilesets) {
				log.Printf("Warning: Layer index %d is out of range for tilesets", layerIndex)
				continue
			}

			img := g.tilesets[layerIndex].Img(id)
			if img == nil {
				continue // Skip this tile if the image is nil
			}

			opts.GeoM.Reset()
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(img, &opts)
		}
	}

	opts.GeoM.Reset()
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	screen.DrawImage(g.player.Image.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image), &opts)

	for _, sprite := range g.enemies {
		opts.GeoM.Reset()
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(sprite.Image.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), &opts)
	}

	for _, sprite := range g.experience {
		opts.GeoM.Reset()
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(sprite.Image.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), &opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/players/samurai.png")
	if err != nil {
		log.Fatal(err)
	}

	lionImage, _, err := ebitenutil.NewImageFromFile("assets/enemys/lion.png")
	if err != nil {
		log.Fatal(err)
	}

	experienceImage, _, err := ebitenutil.NewImageFromFile("assets/misc/goldcoin.png")
	if err != nil {
		log.Fatal(err)
	}

	mapPath := "assets/maps/spawnMap.json"
	tilemapJson, err := entities.NewTilemapJSON(mapPath)
	if err != nil {
		log.Fatal(err)
	}

	rangedClass := &entities.RangedClass{}
	player := &entities.Player{
		Sprite: &entities.Sprite{
			Image: playerImage,
			X:     100,
			Y:     100,
		},
		Health:      100,
		Experience:  0,
		Speed:       2.5,
		PlayerClass: rangedClass,
	}

	tilesets, err := tilemapJson.GenTilesets(mapPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d tilesets", len(tilesets))

	game := Game{
		player: player,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Image: lionImage,
					X:     150,
					Y:     100,
				},
				FollowsPlayer:   true,
				CanAttackPlayer: false,
				CanAttackEnemy:  false,
			},
		},
		experience: []*entities.Coin{
			{
				Sprite: &entities.Sprite{
					Image: experienceImage,
					X:     150,
					Y:     100,
				},
				AmtXp: 100,
			},
		},
		tilemapJson: tilemapJson,
		tilesets:    tilesets,
		cam:         utils.NewCamera(0, 0),
	}

	log.Printf("Player Class: %s, Attack Range: %.2f\n", player.PlayerClass.ClassName(), player.PlayerClass.AttackRange())

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
