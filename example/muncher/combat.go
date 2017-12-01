package main

type CombatSystem struct {
	Messaging
}

func (combat CombatSystem) fight(a Entity, d Entity) {
	attacker, defender := a.Fighter(), d.Fighter()

	defender.Damage(attacker.Level)

	if defender.HP.Current == 0 {
		combat.Broadcast(KillEntity, KillEntityMessage{Attacker: a, Defender: d})
	}
}

func (combat CombatSystem) Notify(message Message, data interface{}) {
	switch message {
	case AttackEntity:
		if d, ok := data.(AttackEntityMesasge); ok {
			combat.fight(d.Attacker, d.Defender)
		}
	}
}
