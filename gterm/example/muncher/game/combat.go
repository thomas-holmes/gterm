package game

type CombatSystem struct {
	Messaging
}

func (combat CombatSystem) fight(player *Player, monster *Monster) {
	switch {
	case player.Level == monster.Level:
		monster.Kill()
		player.GainExp((monster.Level + 1) / 2)
	case player.Level > monster.Level:
		monster.Kill()
		player.GainExp((monster.Level + 1) / 4)
	case player.Level < monster.Level:
		player.Damage(monster.Level)
	}
}

func (combat CombatSystem) Notify(message Message, data interface{}) {
	switch message {
	case PlayerAttack:
		if d, ok := data.(PlayerAttackMessage); ok {
			combat.fight(d.Player, d.Monster)
		}
	}
}
