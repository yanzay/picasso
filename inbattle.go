package main

import (
	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

func (mid *middlewares) inBattle(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		player, err := mid.store.GetPlayer(m.From.ID)
		if err != nil {
			log.Errorf("can't get player: %q", err)
			return
		}
		if player.InBattle() {
			m.Reply("You can't do that while in battle!")
			return
		}
		f(m)
	}
}
