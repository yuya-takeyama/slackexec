package main

import (
	"github.com/nlopes/slack"
	"strings"
)

type Slack struct {
	Name    string
	Icon    string
	Channel string
	Client  *slack.Slack
}

func NewSlack(name, icon, channel, token string) *Slack {
	return &Slack{
		Name:    name,
		Icon:    icon,
		Channel: channel,
		Client:  slack.New(token),
	}
}

func (s *Slack) Post(message string) error {
	_, _, err := s.Client.PostMessage(
		s.Channel,
		message,
		s.getPostMessageParameters(),
	)

	return err
}

func (s *Slack) getPostMessageParameters() slack.PostMessageParameters {
	if strings.HasPrefix(s.Icon, "http") {
		return slack.PostMessageParameters{
			Username: s.Name,
			IconURL:  s.Icon,
		}
	} else {
		return slack.PostMessageParameters{
			Username:  s.Name,
			IconEmoji: s.Icon,
		}
	}
}
