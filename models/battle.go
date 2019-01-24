package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/yanzay/log"
)

type Battle struct {
	ID           int
	Start        time.Time
	Attack       []*Player
	Defence      []*Player
	AttackGuild  string
	DefenceGuild string
}

type BattleResult struct {
	PlayerID int
	Win      bool
	Attack   bool
	Coins    int
	Prestige int
	Winners  []*Player
	Losers   []*Player
}

type BattlePower struct {
	Main    int
	Special int
}

type BattleReward struct {
	Coins    int
	Prestige int
}

func NewBattle(attack *Player, defence *Player) *Battle {
	b := &Battle{
		ID:      rand.Int()/2 + 1,
		Start:   time.Now(),
		Attack:  []*Player{attack},
		Defence: []*Player{defence},
	}
	if !defence.NoGuild() {
		b.AttackGuild = attack.Guild
		b.DefenceGuild = defence.Guild
	}
	return b
}

func (b *Battle) PlayerIDs() []int {
	ids := []int{}
	for _, attacker := range b.Attack {
		ids = append(ids, attacker.ID)
	}
	for _, defender := range b.Defence {
		ids = append(ids, defender.ID)
	}
	return ids
}

func (b *Battle) Guild() bool {
	return b.DefenceGuild != ""
}

func (b *Battle) Results() []*BattleResult {
	results := []*BattleResult{}
	attack := attackPower(b.Attack)
	defence := defencePower(b.Defence)
	log.Infof("Battle results:\n%s", b.String())
	log.Infof("Attack: %d/%d", attack.Main, attack.Special)
	log.Infof("Defence: %d/%d", defence.Main, defence.Special)
	var defenceCoef = 1.0
	var attackCoef = 1.0
	if attack.Special > defence.Special {
		defenceCoef = 0.9
	}
	if attack.Special < defence.Special {
		attackCoef = 0.9
	}
	var winners, losers []*Player
	var winAttack bool
	attackRand := randomness()
	defenceRand := randomness()
	log.Infof("Attack rand: %f, defence rand: %f", attackRand, defenceRand)
	if float64(attack.Main)*attackCoef*attackRand > float64(defence.Main)*defenceCoef*defenceRand {
		winners = b.Attack
		losers = b.Defence
		winAttack = true
		log.Infof("Attack win!")
	} else {
		winners = b.Defence
		losers = b.Attack
		log.Infof("Attack lose!")
	}
	reward := calculateReward(losers)
	for _, winner := range winners {
		result := &BattleResult{
			Win:      true,
			Attack:   winAttack,
			PlayerID: winner.ID,
			Coins:    reward.Coins / len(winners),
			Prestige: reward.Prestige / len(winners),
			Winners:  winners,
			Losers:   losers,
		}
		results = append(results, result)
	}
	for _, loser := range losers {
		result := &BattleResult{
			Win:      false,
			Attack:   !winAttack,
			PlayerID: loser.ID,
			Coins:    loser.Coins / 3,
			Prestige: loser.Prestige / 100,
			Winners:  winners,
			Losers:   losers,
		}
		results = append(results, result)
	}
	return results
}

func attackPower(players []*Player) *BattlePower {
	total := &BattlePower{}
	for _, player := range players {
		total.Main += player.BattleEnergy
		total.Special += player.AttackEnergy
	}
	return total
}

func defencePower(players []*Player) *BattlePower {
	total := &BattlePower{}
	for _, player := range players {
		total.Main += player.BattleEnergy
		total.Special += player.DefenceEnergy
	}
	return total
}

func calculateReward(players []*Player) *BattleReward {
	reward := &BattleReward{}
	for _, player := range players {
		reward.Coins += player.Coins / 3
		reward.Prestige += player.Prestige / 100
	}
	return reward
}

func (b Battle) String() string {
	attackers := make([]string, 0, len(b.Attack))
	defenders := make([]string, 0, len(b.Defence))
	for _, player := range b.Attack {
		attackers = append(attackers, player.FullName())
	}
	for _, player := range b.Defence {
		defenders = append(defenders, player.FullName())
	}
	return fmt.Sprintf("Attack: %s\nDefence: %s", strings.Join(attackers, ", "), strings.Join(defenders, ", "))
}

func randomness() float64 {
	return 1 + (rand.Float64()-0.5)/10
}
