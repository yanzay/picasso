package templates

import (
	"fmt"
	"strings"

	"github.com/yanzay/picasso/models"
)

type BattleTemplates struct{}

func (BattleTemplates) RenderWinResult(p *models.Player, res *models.BattleResult) string {
	s := []string{}
	s = append(s, fmt.Sprintf("Congratulations, %s, you win!", p.FullName()))
	if res.Coins > 0 {
		s = append(s, fmt.Sprintf("Your reward is %d%s.", res.Coins, CoinEmoji))
	}
	if res.Prestige > 0 {
		s = append(s, fmt.Sprintf("Your prestige increased by %d%s.", res.Prestige, PrestigeEmoji))
	}
	winners := "Winners: " + playerList(res.Winners)
	losers := "Losers: " + playerList(res.Losers)
	return strings.Join(s, " ") + fmt.Sprintf("\n%s\n%s", winners, losers)
}

func (BattleTemplates) RenderLoseResult(p *models.Player, res *models.BattleResult) string {
	s := []string{}
	s = append(s, fmt.Sprintf("Unfortunately, %s, you lose.", p.FullName()))
	if res.Coins > 0 {
		s = append(s, fmt.Sprintf("In battle you have lost %d%s.", res.Coins, CoinEmoji))
	}
	if res.Prestige > 0 {
		s = append(s, fmt.Sprintf("Your prestige decreased by %d%s.", res.Prestige, PrestigeEmoji))
	}
	winners := "Winners: " + playerList(res.Winners)
	losers := "Losers: " + playerList(res.Losers)
	return strings.Join(s, " ") + fmt.Sprintf("\n%s\n%s", winners, losers)
}

func (BattleTemplates) RenderAttackedBy(p *models.Player) string {
	return fmt.Sprintf("You are under attack! %s attacked you.", p.FullName())
}

func (BattleTemplates) RenderJoinedGuildBattle(p *models.Player, attack bool) string {
	var action string
	if attack {
		action = "attack"
	} else {
		action = "defence"
	}
	return fmt.Sprintf("%s joined the %s", p.FullName(), action)
}

func (BattleTemplates) RenderGuildAttack(attacker, defender *models.Guild) string {
	return fmt.Sprintf("For the %s! Your guild attacked %s, join and fight!", attacker.Title(), defender.Title())
}

func (BattleTemplates) RenderGuildDefence(attacker *models.Guild) string {
	if attacker == nil {
		return fmt.Sprintf("Your guild is under attack! Aggressor is crazy loner! Join and defend!")
	}
	return fmt.Sprintf("Your guild is under attack! Aggressor: %s. Join and defend!", attacker.Title())
}
