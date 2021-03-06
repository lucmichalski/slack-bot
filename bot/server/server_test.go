package server

import (
	"github.com/innogames/slack-bot/bot/config"
	"github.com/innogames/slack-bot/client"
	"github.com/innogames/slack-bot/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// test if the start/stop is working as expected...main server tests are part of handler_test.go
func TestServer(t *testing.T) {
	cfg := config.Server{
		Listen: "127.0.0.1:6545",
	}
	slackClient := &mocks.SlackClient{}

	server := NewServer(
		cfg,
		slackClient,
	)

	go server.StartServer()
	time.Sleep(time.Millisecond * 10)
	defer server.Stop()

	t.Run("health check", func(t *testing.T) {
		res, err := client.HTTPClient.Get("http://127.0.0.1:6545/health")
		assert.Nil(t, err)

		response := make([]byte, 4)
		res.Body.Read(response)
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "pong", string(response))
	})
}
