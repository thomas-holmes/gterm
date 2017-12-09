package main

import "log"

type CombatSystem struct {
	Messaging
}

type Combatant interface {
	Combatant() *Creature
}

func (combat CombatSystem) fight(a Entity, d Entity) {
	aCombatant, ok := a.(Combatant)
	if !ok {
		return
	}
	dCombatant, ok := d.(Combatant)
	if !ok {
		return
	}
	attacker, defender := aCombatant.Combatant(), dCombatant.Combatant()

	defender.Damage(attacker.Level)

	if defender.HP.Current == 0 {
		combat.Broadcast(KillEntity, KillEntityMessage{Attacker: a, Defender: d})
	}
}

func (combat CombatSystem) Notify(message Message, data interface{}) {
	switch message {
	case AttackEntity:
		if d, ok := data.(AttackEntityMesasge); ok {
			log.Printf("Got a fight message, %+v, %+v, %+v", d, d.Attacker, d.Defender)
			combat.fight(d.Attacker, d.Defender)
		}
	}
}
