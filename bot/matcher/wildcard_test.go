package matcher

import (
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWildcard(t *testing.T) {
	t.Run("Match", func(t *testing.T) {
		runner := func(ref msg.Ref, text string) bool {
			return true
		}
		subject := WildcardMatcher(runner)

		message := msg.Message{}
		message.Text = "any"
		_, match := subject.Match(message)

		assert.True(t, match.Matched())
	})

	t.Run("NoMatch", func(t *testing.T) {
		runner := func(ref msg.Ref, text string) bool {
			return false
		}
		subject := WildcardMatcher(runner)

		message := msg.Message{}
		message.Text = "any"
		_, match := subject.Match(message)

		assert.False(t, match.Matched())
	})
}
