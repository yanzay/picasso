package game

import (
	"fmt"
	"os"
	"testing"

	"github.com/yanzay/log"
	"github.com/yanzay/picasso/models"
	"github.com/yanzay/picasso/storage"
	"github.com/yanzay/picasso/templates"
)

var (
	sendMessage = func(id int, text string) { fmt.Println(id, text) }
	store       = storage.New("_test.db")
	game        = New(store, sendMessage, templates.BattleTemplates{})
)

func TestMain(m *testing.M) {
	log.Writer = EmptyWriter{}
	createPlayers(store)
	status := m.Run()
	deletePlayers(store)
	os.Exit(status)
}

func createPlayers(store storage.Storage) {
	for i := 0; i < 10; i++ {
		player := models.NewPlayer(i)
		store.SetPlayer(player)
	}
}

func deletePlayers(store storage.Storage) {
	players, _ := store.GetAllPlayers()
	for _, player := range players {
		store.DeletePlayer(player.ID)
	}
}
