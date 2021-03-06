package jenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/innogames/slack-bot/bot"
	"github.com/innogames/slack-bot/bot/matcher"
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/innogames/slack-bot/bot/util"
	"github.com/innogames/slack-bot/client/jenkins"
	"github.com/innogames/slack-bot/command/queue"
	"github.com/slack-go/slack"
	"time"
)

type buildWatcherCommand struct {
	jenkinsCommand
}

// newBuildWatcherCommand watches the status of an already running jenkins build
func newBuildWatcherCommand(base jenkinsCommand) bot.Command {
	return &buildWatcherCommand{base}
}

func (c *buildWatcherCommand) GetMatcher() matcher.Matcher {
	return matcher.NewRegexpMatcher(`(notify|inform)( me about)? (job|build) ?(?P<job>[\w\-_]*)( #?(?P<build>\d+))?`, c.Run)
}

func (c *buildWatcherCommand) IsEnabled() bool {
	return c.jenkins != nil
}

func (c *buildWatcherCommand) Run(match matcher.Result, message msg.Message) {
	jobName := match.GetString("job")
	buildNumber := match.GetInt("build")

	job, err := c.jenkins.GetJob(jobName)
	if err != nil {
		c.SendMessage(message, fmt.Sprintf("Job *%s* does not exist", jobName))
		return
	}

	build, err := getBuild(job, buildNumber)
	if err != nil {
		c.ReplyError(message, err)
		return
	}

	if !build.Raw.Building {
		c.SendMessage(message, fmt.Sprintf("No job for *%s* is running right now", jobName))
		return
	}

	text := fmt.Sprintf(
		"Okay, I'll inform you when the job %s #%s is done",
		jobName,
		build.Info().ID,
	)
	attachment := jenkins.GetAttachment(build, text)
	msgTimestamp := c.SendMessage(message, "", attachment)

	done := queue.AddRunningCommand(
		message,
		fmt.Sprintf("inform job %s #%d", jobName, build.GetBuildNumber()),
	)

	go func() {
		<-jenkins.WatchBuild(build)
		done <- true

		c.SendMessage(
			message,
			"",
			slack.MsgOptionUpdate(msgTimestamp),
			jenkins.GetAttachment(build, text),
		)

		c.RemoveReaction(jenkins.IconRunning, message)
		if build.IsGood() {
			c.AddReaction(jenkins.IconSuccess, message)
		} else {
			c.AddReaction(jenkins.IconFailed, message)
		}

		duration := time.Duration(build.GetDuration()) * time.Millisecond
		c.SendMessage(message, fmt.Sprintf(
			"<@%s> *%s*: %s #%s: %s in %s",
			message.User,
			build.GetResult(),
			jobName,
			build.Info().ID,
			build.GetUrl(),
			util.FormatDuration(duration),
		))
	}()
}

func getBuild(job jenkins.Job, buildNumber int) (*gojenkins.Build, error) {
	if buildNumber == 0 {
		_, err := job.Poll()
		if err != nil {
			return nil, err
		}

		return job.GetLastBuild()
	}
	return job.GetBuild(int64(buildNumber))
}

func (c *buildWatcherCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "notify build",
			Description: "inform you when a running jenkins build finishes",
			Examples: []string{
				"inform me about build AtcBrowser #1233",
				"inform me about build AtcMobile",
				"notify build AtcMobile",
			},
			Category: category,
		},
		{
			Command:     "inform job",
			Description: "inform you when a running jenkins build finishes",
			Examples: []string{
				"inform me about build AtcBrowser #1233",
				"inform me about build AtcMobile",
				"notify build AtcMobile",
			},
			Category: category,
		},
	}
}
