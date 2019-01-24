package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/picasso/models"
	"github.com/yanzay/picasso/storage"
)

var (
	dbFile = flag.String("db", "data.db", "Database file")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	store := storage.New(*dbFile)
	for _, guild := range models.Guilds {
		existing, err := store.GetGuild(guild.ID)
		if err != nil {
			log.Infof("Error getting guild: %v", err)
			log.Infof("New guild: %s", guild.Name)
			guild.PlayerIDs = []int{}
			err := store.SetGuild(&guild)
			if err != nil {
				log.Fatalf("unable to save guild: %v", err)
			}
			continue
		}
		existing.PlayerIDs = []int{}
		store.SetGuild(existing)
	}
	players, err := store.GetAllPlayers()
	if err != nil {
		log.Fatalf("unable to get all players: %v", err)
	}
	for _, player := range players {
		if !player.NoGuild() {
			guild, err := models.GuildByEmoji(player.Guild)
			if err != nil {
				log.Fatalf("unable to get guild by emoji %s: %v", player.Guild, err)
			}
			player.GuildID = guild.ID
			log.Infof("Saving player: %#v", *player)
			err = store.SetPlayer(player)
			if err != nil {
				log.Fatalf("unable to save player: %v", err)
			}
			g, err := store.GetGuild(guild.ID)
			if err != nil {
				log.Fatal(err)
			}
			g.PlayerIDs = append(g.PlayerIDs, player.ID)
			err = store.SetGuild(g)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
