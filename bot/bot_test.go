package bot

import (
	"github.com/innogames/slack-bot/bot/config"
	"github.com/innogames/slack-bot/bot/matcher"
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/innogames/slack-bot/client"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBot(t *testing.T) {
	cfg := config.Config{}

	rawSlackClient := &slack.Client{}
	slackClient := &client.Slack{Client: rawSlackClient, RTM: rawSlackClient.NewRTM()}

	commands := &Commands{}
	commands.AddCommand(testCommand2{})

	bot := NewBot(cfg, slackClient, commands)
	bot.auth = &slack.AuthTestResponse{
		UserID: "BOT",
	}
	bot.allowedUsers = map[string]string{
		"U123": "User123",
	}

	t.Run("handle empty message", func(t *testing.T) {
		message := msg.Message{}
		message.Text = ""
		message.Channel = "C123"
		bot.handleMessage(message, true)
	})

	t.Run("handle unauthenticated message", func(t *testing.T) {
		message := msg.Message{}

		message.Text = "test"
		message.User = "U888"
		message.Channel = "C123"
		bot.handleMessage(message, false)
	})

	t.Run("handle valid message", func(t *testing.T) {
		message := msg.Message{}

		message.Text = "test"
		message.User = "U123"
		message.Channel = "C123"
		bot.handleMessage(message, true)
	})
}

func TestIsBotMessage(t *testing.T) {
	bot := Bot{}
	bot.auth = &slack.AuthTestResponse{
		UserID: "BOT",
	}

	t.Run("Is random message", func(t *testing.T) {
		event := &slack.MessageEvent{}
		event.Text = "random text"
		actual := bot.canHandleMessage(event)
		assert.Equal(t, false, actual)
	})

	t.Run("Is random message to other user", func(t *testing.T) {
		event := &slack.MessageEvent{}
		event.Text = "<@USER2> random test"

		actual := bot.canHandleMessage(event)
		assert.Equal(t, false, actual)
	})

	t.Run("Is Bot mentioned", func(t *testing.T) {
		event := &slack.MessageEvent{}
		event.User = "U1234"
		event.Text = "<@BOT> random test"

		actual := bot.canHandleMessage(event)
		assert.Equal(t, true, actual)
	})

	t.Run("Pass internal events", func(t *testing.T) {
		event := &slack.MessageEvent{}
		event.User = "U1234"
		event.Text = "<@BOT> random test"

		actual := bot.canHandleMessage(event)
		assert.Equal(t, true, actual)
	})

	t.Run("Is private channel", func(t *testing.T) {
		event := &slack.MessageEvent{
			Msg: slack.Msg{
				Channel: "DRANDOM",
				User:    "U1234",
			},
		}
		event.Text = "random test"
		actual := bot.canHandleMessage(event)
		assert.Equal(t, true, actual)
	})

	t.Run("Is random channel", func(t *testing.T) {
		event := &slack.MessageEvent{
			Msg: slack.Msg{
				Channel: "GRANDOM",
				User:    "U1234",
			},
		}
		event.Text = "random test"
		actual := bot.canHandleMessage(event)
		assert.Equal(t, false, actual)
	})

	t.Run("Trim", func(t *testing.T) {
		assert.Equal(t, bot.cleanMessage(" ", true), "")
		assert.Equal(t, bot.cleanMessage("<@BOT> random ’test’", true), "random 'test'")
		assert.Equal(t, bot.cleanMessage("<https://test.com|TEST> <https://example.com|example>", false), "<https://test.com|TEST> <https://example.com|example>")
		assert.Equal(t, bot.cleanMessage("<https://test.com|TEST> <https://example.com|example>", true), "TEST example")
	})
}

func BenchmarkTrimMessage(b *testing.B) {
	bot := Bot{}
	bot.auth = &slack.AuthTestResponse{}
	bot.auth.User = "botId"

	message := "<@botId> hallo how are `you’?"

	b.Run("trim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bot.cleanMessage(message, false)
		}
	})
}

func BenchmarkShouldHandle(b *testing.B) {
	bot := Bot{}
	bot.auth = &slack.AuthTestResponse{}
	bot.auth.User = "botId"

	var result bool

	b.Run("match", func(b *testing.B) {
		event := &slack.MessageEvent{}
		event.Channel = "D123"
		event.User = "U123"

		for i := 0; i < b.N; i++ {
			result = bot.canHandleMessage(event)
		}
		assert.True(b, result)
	})
}

type testCommand2 struct {
}

func (c testCommand2) GetMatcher() matcher.Matcher {
	return matcher.NewTextMatcher("test", func(match matcher.Result, message msg.Message) {
	})
}
