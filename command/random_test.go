package command

import (
	"github.com/innogames/slack-bot/bot"
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/innogames/slack-bot/mocks"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRandom(t *testing.T) {
	slackClient := mocks.SlackClient{}

	randomCommand := NewRandomCommand(&slackClient).(*randomCommand)
	randomCommand.random.Seed(1) // we want always the same random

	command := bot.Commands{}
	command.AddCommand(randomCommand)

	t.Run("invalid command", func(t *testing.T) {
		message := msg.Message{}
		message.Text = "randomness"

		actual := command.Run(message)
		assert.Equal(t, false, actual)
	})

	t.Run("no options should not match", func(t *testing.T) {
		message := msg.Message{}
		message.Text = "random"

		slackClient.On("SendMessage", message, "You have to pass more arguments").Return("")

		actual := command.Run(message)
		assert.Equal(t, true, actual)
	})

	t.Run("pick random with one entry", func(t *testing.T) {
		message := msg.Message{}
		message.Text = "random 1"

		slackClient.On("SendMessage", message, "1").Return("")

		actual := command.Run(message)
		assert.Equal(t, true, actual)
	})

	t.Run("pick random entry", func(t *testing.T) {
		rand.Seed(1) // we want always the same random

		message := msg.Message{}
		message.Text = "random 4 5 6"

		// seed was chosen to pick the "3" every time
		slackClient.On("SendMessage", message, "4").Return("")

		actual := command.Run(message)
		assert.Equal(t, true, actual)
	})
}
