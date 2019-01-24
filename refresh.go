package main

import (
	"github.com/yanzay/log"
	"github.com/yanzay/picasso/templates"
	"github.com/yanzay/tbot"
)

func (mid *middlewares) refresh(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		if m.Data == "/refresh" {
			mid.store.Reset(int64(m.From.ID))
			player, err := mid.store.GetPlayer(m.From.ID)
			if err != nil {
				log.Errorf("unable to get player: %v", err)
				return
			}
			templates.RenderPlayer(m, player)
			return
		}
		f(m)
	}
}
