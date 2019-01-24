package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/yanzay/log"
	"github.com/yanzay/picasso/models"
	"github.com/yanzay/picasso/storage"
)

var (
	dbFile = flag.String("db", "data.db", "Database file")
	number = flag.Int("n", 1000, "Number of players")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	store := storage.New(*dbFile)
	for i := 0; i < *number; i++ {
		player := genPlayer()
		log.Infof("Saving player: %#v", player)
		store.SetPlayer(player)
	}
}

func genPlayer() *models.Player {
	p := &models.Player{
		ID:        rand.Int(),
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Blocked:   false,
		Coins:     rand.Intn(1000000),
		Implants: models.Implants{
			Battery:  rand.Intn(100) + 10,
			Exoframe: rand.Intn(100) + 10,
			Shield:   rand.Intn(100) + 10,
			Weapon:   rand.Intn(100) + 10,
		},
		Equipment: models.Equipment{
			CoinMiner:      rand.Intn(50) + 20,
			Backpack:       rand.Intn(50) + 20,
			ResourceMiner1: rand.Intn(50) + 20,
			ResourceMiner2: rand.Intn(50) + 20,
			ResourceMiner3: rand.Intn(50) + 20,
		},
		Bot: true,
	}
	return p
}
