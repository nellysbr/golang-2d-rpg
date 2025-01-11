package entities

type PlayerClass interface {
	ClassName() string
	AttackRange() float64
}

// Classe Melee
type MeleeClass struct{}

// Classe Ranged
type RangedClass struct {
	RangeBonus float64
}

type Player struct {
	*Sprite
	Health      uint
	Experience  uint
	Speed       float64
	PlayerClass PlayerClass
}

func (r *RangedClass) ClassName() string {
	return "Ranged"
}

func (r *RangedClass) AttackRange() float64 {
	return 10.0 + r.RangeBonus
}

func (m *MeleeClass) ClassName() string {
	return "Melee"
}

func (m *MeleeClass) AttackRange() float64 {
	return 1.5
}
