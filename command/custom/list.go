package custom

import (
	"fmt"
	"github.com/innogames/slack-bot/bot/matcher"
	"github.com/innogames/slack-bot/bot/msg"
)

func (c command) List(match matcher.Result, message msg.Message) {
	list := loadList(message)
	if len(list) == 0 {
		c.slackClient.SendMessage(message, "No commands define yet. Use `add command 'your alias' 'command to execute'`")
		return
	}

	responseText := fmt.Sprintf("You defined %d commands:", len(list))
	for alias, command := range list {
		responseText += fmt.Sprintf("\n - %s: `%s`", alias, command)
	}

	c.slackClient.SendMessage(message, responseText)
}
