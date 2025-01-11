package entities

type Enemy struct {
	*Sprite
	FollowsPlayer   bool
	CanAttackPlayer bool
	CanAttackEnemy  bool
}
