package models

import (
	"fmt"
)

type Guild struct {
	ID    int
	Emoji string
	Name  string

	BattleID  int
	PlayerIDs []int

	Prestige int
	Wins     int
	Coins    int

	ChatID int64
}

func (g *Guild) Title() string {
	return fmt.Sprintf("[%s]%s", g.Emoji, g.Name)
}

var Guilds = []Guild{
	Guild{ID: 1, Emoji: RhinoEmoji, Name: "Rhino"},
	Guild{ID: 2, Emoji: ScorpionEmoji, Name: "Scorpion"},
	Guild{ID: 3, Emoji: LizardEmoji, Name: "Lizard"},
}

func GuildByTitle(title string) (Guild, error) {
	for _, guild := range Guilds {
		if guild.Title() == title {
			return guild, nil
		}
	}
	return Guild{}, fmt.Errorf("Guild %s not found", title)
}

func GuildByEmoji(emoji string) (Guild, error) {
	for _, guild := range Guilds {
		if guild.Emoji == emoji {
			return guild, nil
		}
	}
	return Guild{}, fmt.Errorf("Guild %s not found", emoji)
}

func (g *Guild) InBattle() bool {
	return g.BattleID != 0
}
