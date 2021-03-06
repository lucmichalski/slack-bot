package cron

import (
	"github.com/innogames/slack-bot/bot"
	"github.com/innogames/slack-bot/bot/config"
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/innogames/slack-bot/bot/util"
	"github.com/innogames/slack-bot/client"
	cronLib "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// NewCronCommand registers cron which are configurable in the yaml config
func NewCronCommand(slackClient client.SlackClient, crons []config.Cron) bot.Command {
	if len(crons) == 0 {
		return nil
	}

	cron := cronLib.New()
	cmd := &command{slackClient, crons, cron}

	for _, cronCommand := range crons {
		_, err := cron.AddFunc(cronCommand.Schedule, cmd.getCallback(cronCommand))
		if err != nil {
			log.Error(err)
		}
	}

	cron.Start()
	log.Infof("Initialized %d crons", len(crons))

	return cmd
}

type command struct {
	slackClient client.SlackClient
	cfg         []config.Cron
	cron        *cronLib.Cron
}

func (c *command) getCallback(cron config.Cron) func() {
	// todo validate template before execution. but this is tricky as some functions gets registered later on...
	return func() {
		for _, commandTemplate := range cron.Commands {
			command, err := util.CompileTemplate(commandTemplate)
			if err != nil {
				log.Error(err)
				continue
			}
			text, err := util.EvalTemplate(command, util.Parameters{})
			if err != nil {
				log.Error(err)
				continue
			}

			newMessage := msg.Message{}
			newMessage.User = "cron"
			newMessage.Channel, _ = client.GetChannel(cron.Channel)
			newMessage.Text = text
			client.InternalMessages <- newMessage
		}
	}
}
