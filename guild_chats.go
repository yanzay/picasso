package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/picasso/models"
	"github.com/yanzay/tbot"
	"github.com/yanzay/tbot/model"
)

func guildChatHandler(m *tbot.Message) {
	handleAdminCommand(m)
	if m.Type != model.MessageNewMembers {
		return
	}
	for _, user := range m.Users {
		if !userInGuild(user.ID, m.ChatID) {
			err := app.bot.SendRaw("kickChatMember", map[string]string{
				"chat_id":    fmt.Sprint(m.ChatID),
				"user_id":    fmt.Sprint(m.From.ID),
				"until_date": fmt.Sprint(time.Now().Add(40 * time.Second).Unix()),
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func handleAdminCommand(m *tbot.Message) {
	if m.Type != model.MessageText {
		return
	}
	if m.From.UserName != *admin {
		return
	}
	if !strings.HasPrefix(m.Data, "/guildchat") {
		return
	}
	idstr := strings.TrimPrefix(m.Data, "/guildchat ")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		m.Reply("invalid id")
		return
	}
	guild, err := app.store.GetGuild(id)
	if err != nil {
		log.Errorf("unable to get guild %d: %v", id, err)
		m.Reply("Something went wrong")
		return
	}
	guild.ChatID = m.ChatID
	err = app.store.SetGuild(guild)
	if err != nil {
		log.Errorf("unable to set guild %d: %v", id, err)
		m.Reply("Something went wrong")
		return
	}
	m.Reply("OK")
}

func userInGuild(userID int, chatID int64) bool {
	guilds, err := app.store.GetAllGuilds()
	if err != nil {
		log.Errorf("unable to get guilds: %v", err)
		return false
	}
	var guild *models.Guild
	for _, g := range guilds {
		if g.ChatID == chatID {
			guild = g
			break
		}
	}
	if guild == nil {
		log.Warningf("not a guild supergroup! %d", chatID)
		return false
	}
	for _, playerID := range guild.PlayerIDs {
		if playerID == userID {
			return true
		}
	}
	return false
}
