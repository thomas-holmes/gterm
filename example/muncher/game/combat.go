package game

type CombatSystem struct {
	Messaging
}

func (combat CombatSystem) fight(attacker *Creature, defender *Creature) {
	defender.Damage(attacker.Level)

	if defender.HP.Current == 0 {
		combat.Broadcast(KillEntity, KillEntityMessage{Attacker: attacker, Defender: defender})
	}
}

func (combat CombatSystem) Notify(message Message, data interface{}) {
	switch message {
	case CreatureAttack:
		if d, ok := data.(CreatureAttackMessage); ok {
			combat.fight(d.Attacker, d.Defender)
		}
	}
}
