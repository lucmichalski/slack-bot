package command

import (
	"github.com/innogames/slack-bot/bot"
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/innogames/slack-bot/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReaction(t *testing.T) {
	slackClient := mocks.SlackClient{}
	reaction := NewReactionCommand(&slackClient)

	command := bot.Commands{}
	command.AddCommand(reaction)

	t.Run("invalid command", func(t *testing.T) {
		message := msg.Message{}
		message.Text = "i need a reaction"

		actual := command.Run(message)
		assert.Equal(t, false, actual)
	})

	t.Run("add reaction", func(t *testing.T) {
		message := msg.Message{}
		message.Text = "add reaction :test:"
		message.Channel = "chan"
		message.Timestamp = "time"

		slackClient.On("AddReaction", "test", message)
		actual := command.Run(message)
		assert.Equal(t, true, actual)
	})

	t.Run("remove reaction", func(t *testing.T) {
		message := msg.Message{}
		message.Text = "remove reaction :test:"
		message.Channel = "chan"
		message.Timestamp = "time"

		slackClient.On("RemoveReaction", "test", message)
		actual := command.Run(message)
		assert.Equal(t, true, actual)
	})
}
