package main

import (
	"flag"
	"os"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"

	"github.com/yanzay/picasso/game"
	"github.com/yanzay/picasso/storage"
	"github.com/yanzay/picasso/templates"
)

var (
	dbFile     = flag.String("data", "data.db", "Database file")
	minute     = flag.Int("m", 60, "Game minute duration in seconds")
	local      = flag.Bool("local", false, "Start without webhook")
	webhookURL = flag.String("url", "https://picasso.yanzay.com/", "Webhook base URL")
	admin      = flag.String("admin", "yanzay", "Admin username")
)

const (
	supportChatURL = "https://t.me/joinchat/BYwd6URAR2sUmGEr1TuNKg"
	wikiURL        = "http://telegra.ph/Picasso-Game---Quick-Start-10-14"
)

type application struct {
	bot   *tbot.Server
	store storage.Storage
}

var app *application

func init() {
	flag.Parse()
}

func main() {
	app = &application{}

	store := storage.New(*dbFile)
	app.store = store

	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := createBot(token, *webhookURL, store, *local)
	if err != nil {
		log.Fatalf("unable to create bot: %v", err)
	}
	app.bot = bot

	log.Infof("Starting game...")
	nots := make(chan *notification, 1000)
	go notifier(bot.Send, nots)
	sendFunc := func(chatID int, text string) {
		nots <- &notification{id: chatID, message: text}
	}
	game := game.New(store, sendFunc, templates.BattleTemplates{})
	game.Start()
	mids := setMiddlewares(bot, store)
	setHandlers(bot, mids, store, game)

	log.Infof("Staring server...")
	log.Fatal(bot.ListenAndServe())
}

func createBot(token string, webhookURL string, store storage.Storage, local bool) (*tbot.Server, error) {
	mux := tbot.NewRouterMux(store)
	if local {
		return tbot.NewServer(token, tbot.WithMux(mux))
	}
	return tbot.NewServer(token,
		tbot.WithMux(mux),
		tbot.WithWebhook(webhookURL+token, ":8014"))
}

type notification struct {
	id      int
	message string
}

func notifier(sendFunc func(int64, string) error, ch chan *notification) {
	for {
		n := <-ch
		start := time.Now()
		err := sendFunc(int64(n.id), n.message)
		log.Infof("Message latency: %s", time.Since(start))
		if err != nil {
			log.Errorf("unable to send: %v", err)
		}
	}
}
