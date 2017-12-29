package main

import (
	"fmt"
	"log"
)

type CombatSystem struct {
	Messaging
}

func (combat CombatSystem) fight(a Entity, d Entity) {
	attacker, ok := a.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", a)
		return
	}
	defender, ok := d.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", d)
		return
	}

	log.Printf("Fighting with %+v", attacker.Equipment)
	defender.Damage(attacker.Equipment.Weapon.Power)
	logString := fmt.Sprintf("%v hits %v for %v damage!", attacker.Name, defender.Name, attacker.Equipment.Weapon.Power)

	combat.Broadcast(GameLogAppend, GameLogAppendMessage{[]string{logString}})

	// This should be done by the entity instead of here, I think?
	// I think this used to attribute the experience gain on death. Maybe
	// need a more sophisticated combat/exp tracking system instead based
	// on damage dealt & proximity?
	if defender.HP.Current == 0 {
		combat.Broadcast(KillEntity, KillEntityMessage{Attacker: a, Defender: d})
	}
}

func (combat CombatSystem) zap(a Entity, d Entity, s Spell) {
	attacker, ok := a.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", a)
		return
	}
	defender, ok := d.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", d)
		return
	}

	log.Printf("Spell attacking with %+v", s)
	defender.Damage(s.Power)
	logString := fmt.Sprintf("%v hits %v with %v for %v damage!", attacker.Name, defender.Name, s.Name, s.Power)

	combat.Broadcast(GameLogAppend, GameLogAppendMessage{[]string{logString}})

	// This should be done by the entity instead of here, I think?
	// I think this used to attribute the experience gain on death. Maybe
	// need a more sophisticated combat/exp tracking system instead based
	// on damage dealt & proximity?
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
	case SpellAttackEntity:
		if d, ok := data.(SpellAttackEntityMessage); ok {
			log.Printf("Got a spell attack, %v, %+v, %+v, %+v", d.Spell.Name, d, d.Attacker, d.Defender)
			combat.zap(d.Attacker, d.Defender, d.Spell)
		}
	}
}
