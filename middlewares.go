package main

import (
	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
	"github.com/yanzay/tbot/model"

	"github.com/yanzay/picasso/storage"
)

type middlewares struct{ store storage.Storage }

func setMiddlewares(bot *tbot.Server, store storage.Storage) *middlewares {
	m := &middlewares{store}
	bot.AddMiddleware(m.login)
	bot.AddMiddleware(m.refresh)
	bot.AddMiddleware(func(f tbot.HandlerFunction) tbot.HandlerFunction {
		return func(m *tbot.Message) {
			log.Infof("[%d] New message: %v", m.ChatID, m.Data)
			f(m)
		}
	})
	bot.AddMiddleware(func(f tbot.HandlerFunction) tbot.HandlerFunction {
		return func(m *tbot.Message) {
			if m.ChatType == model.ChatTypeSuperGroup {
				log.Infof("Guild chat message %d: %s", m.ChatID, m.Data)
				guildChatHandler(m)
				return
			}
			f(m)
		}
	})
	return m
}
