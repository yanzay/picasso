package main

import "github.com/yanzay/tbot"

func (mid *middlewares) onlyAdmin(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		if m.From.UserName == *admin {
			f(m)
		} else {
			m.Reply("Access Denied")
		}
	}
}
