package main

import (
	"testing"

	"github.com/yanzay/tbot"
	"github.com/yanzay/tbot/model"

	"github.com/yanzay/picasso/storage"
)

func TestLogin(t *testing.T) {
	ch := make(chan *model.Message)
	go func() {
		for {
			<-ch
		}
	}()

	newMessage := func(text string) *tbot.Message {
		m := &tbot.Message{Message: &model.Message{Data: text, From: model.User{ID: 1}}}
		m.SetReplyChannel(ch)
		return m
	}

	store := storage.New("_test.db")
	m := &middlewares{store}
	triggered := false
	f := func(*tbot.Message) { triggered = true }

	m.login(f)(newMessage("/start"))
	m.login(f)(newMessage("Jon"))
	m.login(f)(newMessage("Snow"))
	m.login(f)(newMessage("ℹ️ Info"))

	if triggered != true {
		t.Error("Should trigger handler")
	}
}
