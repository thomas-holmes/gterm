package game

type CombatSystem struct {
	Messaging
}

func (combat CombatSystem) fight(attacker *Creature, defender *Creature) {
	switch {
	case attacker.Level == defender.Level:
		combat.Broadcast(KillEntity, KillEntityMessage{Attacker: attacker, Defender: defender})
	case attacker.Level > defender.Level:
		combat.Broadcast(KillEntity, KillEntityMessage{Attacker: attacker, Defender: defender})
	case attacker.Level < defender.Level:
		attacker.Damage(defender.Level)
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
