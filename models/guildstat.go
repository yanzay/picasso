package models

type GuildStat struct {
	Title         string
	Leader        *Player
	PlayersNumber int
	Prestige      int
	Wins          int
	Coins         int
}
