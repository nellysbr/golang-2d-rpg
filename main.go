package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Image *ebiten.Image
	X, Y  float64
}

// Interface para classes de jogadores
type PlayerClass interface {
	ClassName() string
	AttackRange() float64
}

// Classe Melee
type MeleeClass struct{}

func (m *MeleeClass) ClassName() string {
	return "Melee"
}

func (m *MeleeClass) AttackRange() float64 {
	return 1.5
}

// Classe Ranged
type RangedClass struct {
	RangeBonus float64
}

func (r *RangedClass) ClassName() string {
	return "Ranged"
}

func (r *RangedClass) AttackRange() float64 {
	return 10.0 + r.RangeBonus
}

// Struct base do jogador
type Player struct {
	*Sprite
	Health      uint
	Experience  uint
	Speed       float64
	PlayerClass PlayerClass
}

type Enemy struct {
	*Sprite
	FollowsPlayer   bool
	CanAttackPlayer bool
	CanAttackEnemy  bool
}

type Coin struct {
	*Sprite
	AmtXp uint
}

type Game struct {
	player     *Player
	enemies    []*Enemy
	experience []*Coin
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.X, g.player.Y)

	screen.DrawImage(g.player.Image.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image), opts)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(sprite.Image.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), opts)

		opts.GeoM.Reset()
	}

	for _, sprite := range g.experience {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(sprite.Image.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image), opts)

		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/samurai.png")
	if err != nil {
		log.Fatal(err)
	}

	lionImage, _, err := ebitenutil.NewImageFromFile("assets/lion.png")
	if err != nil {
		log.Fatal(err)
	}

	experienceImage, _, err := ebitenutil.NewImageFromFile("assets/goldcoin.png")
	if err != nil {
		log.Fatal(err)
	}

	// Inicialização do jogador com uma classe Melee
	rangedClass := &RangedClass{}
	player := &Player{
		Sprite: &Sprite{
			Image: playerImage,
			X:     100,
			Y:     100,
		},
		Health:      100,
		Experience:  0,
		Speed:       2.5,
		PlayerClass: rangedClass,
	}

	// Criação do jogo
	game := Game{
		player: player,
		enemies: []*Enemy{
			{
				Sprite: &Sprite{
					Image: lionImage,
					X:     150,
					Y:     100,
				},
				FollowsPlayer:   true,
				CanAttackPlayer: false,
				CanAttackEnemy:  false,
			},
		},
		experience: []*Coin{
			{
				Sprite: &Sprite{
					Image: experienceImage,
					X:     150,
					Y:     100,
				},
				AmtXp: 100,
			},
		},
	}

	// Log para testar o alcance do ataque
	log.Printf("Player Class: %s, Attack Range: %.2f\n", player.PlayerClass.ClassName(), player.PlayerClass.AttackRange())

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
