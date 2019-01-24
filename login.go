package main

import (
	"regexp"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"

	"github.com/yanzay/picasso/models"
	"github.com/yanzay/picasso/templates"
)

// Question is a generic question structure for surveys
type Question struct {
	Key               string
	Prompt            string
	Options           []string
	ValidationRule    string
	ValidationComment string
}

func (q *Question) isValidAnswer(answer string) (string, bool) {
	if q.ValidationRule != "" {
		match, err := regexp.MatchString(q.ValidationRule, answer)
		if err != nil {
			log.Errorf("error matching validation rule: %q", err)
		}
		if !match {
			return q.ValidationComment, false
		}
	}
	return "", true
}

var questions = map[string]*Question{
	"first_name": {
		Key:               "first_name",
		Prompt:            "Choose first name for your character:",
		ValidationRule:    "^[A-Z][A-Za-z]*$",
		ValidationComment: "First name should start with capital letter and contain only letters",
	},
	"last_name": {
		Key:               "last_name",
		Prompt:            "Choose last name for your character:",
		ValidationRule:    "^[A-Z][A-Za-z]*$",
		ValidationComment: "Last name should start with capital letter and contain only letters",
	},
}

func (mid *middlewares) login(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		player, err := mid.store.GetPlayer(m.From.ID)
		if err != nil {
			log.Errorf("can't get player: %q", err)
			return
		}
		if player.Blocked {
			m.Reply("You're blocked. Have a nice life!")
			return
		}
		if player.IsFull() {
			f(m)
			return
		}
		mid.survey(m, player)
	}
}

func (mid *middlewares) survey(m *tbot.Message, player *models.Player) {
	survey, err := mid.store.GetSurvey("login", m.From.ID)
	if err != nil {
		log.Errorf("can't get survey: %q", err)
		return
	}

	// already asked a question
	if survey.Asking != "" {
		comment := setAnswer(player, questions[survey.Asking], m.Text())
		if comment != "" {
			m.Reply(comment)
		} else {
			survey.Asking = ""
			mid.store.SetSurvey("login", m.From.ID, survey)
			mid.store.SetPlayer(player)
		}
	}

	if player.Blocked {
		m.Reply("You're blocked. Have a nice life!")
		return
	}

	// ask next question
	if !player.IsFull() {
		survey.Asking = askNext(player, m)
		err = mid.store.SetSurvey("login", m.From.ID, survey)
		if err != nil {
			log.Errorf("can't save survey %s: %q", "login", err)
		}
		return
	}

	// user registered, we're done here
	if player.IsFull() {
		mid.store.Reset(m.ChatID)
		m.Reply("Registered")
		templates.RenderPlayer(m, player)
	}
}

func askNext(player *models.Player, m *tbot.Message) string {
	question := nextQuestion(player)
	if question == nil {
		return ""
	}
	if question.Options != nil {
		m.ReplyKeyboard(question.Prompt, [][]string{question.Options}, tbot.OneTimeKeyboard)
	} else {
		m.Reply(question.Prompt)
	}
	return question.Key
}

func nextQuestion(player *models.Player) *Question {
	switch {
	case player.FirstName == "":
		return questions["first_name"]
	case player.LastName == "":
		return questions["last_name"]
	}
	return nil
}

func setAnswer(player *models.Player, question *Question, answer string) string {
	if comment, ok := question.isValidAnswer(answer); !ok {
		return comment
	}
	switch question.Key {
	case "first_name":
		player.FirstName = answer
	case "last_name":
		player.LastName = answer
	}
	return ""
}
