package main

import (
	"bytes"
	"flag"
	"github.com/nlopes/slack"
	"io"
	"os"
	"os/exec"
	"strings"
)

var channel string
var name string
var icon string

func initFlags() {
	flag.StringVar(&channel, "channel", "#general", "channel to post message")
	flag.StringVar(&name, "name", "slackexec", "username of the bot")
	flag.StringVar(&icon, "icon", ":computer:", "icon of the bot")
	flag.Parse()
}

func main() {
	initFlags()
	client := slack.New(os.Getenv("SLACK_API_TOKEN"))

	flagArgs := flag.Args()
	executable := flagArgs[0]
	args := flagArgs[1:]

	client.PostMessage(
		channel,
		"Running below command..."+"```\n$ "+executable+" "+strings.Join(args, " ")+"```",
		slack.PostMessageParameters{
			Username:  name,
			IconEmoji: icon,
		},
	)

	cmd := exec.Command(executable, args...)

	buf := new(bytes.Buffer)

	writer := io.MultiWriter(buf, os.Stdout)

	cmd.Stdout = writer
	cmd.Stderr = writer

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	client.PostMessage(
		channel,
		"Result:\n```\n"+buf.String()+"```",
		slack.PostMessageParameters{
			Username:  name,
			IconEmoji: icon,
		},
	)
}
