package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/yuya-takeyama/posixexec"
	"io"
	"os"
	"os/exec"
	"os/user"
)

var channel string
var name string
var icon string

const (
	ExitFatal = 111
)

func init() {
	flag.StringVar(&channel, "channel", "#general", "channel to post message")
	flag.StringVar(&name, "name", "slackexec", "username of the bot")
	flag.StringVar(&icon, "icon", ":computer:", "icon of the bot")
	flag.Parse()
}

func main() {
	client := NewSlack(name, icon, channel, os.Getenv("SLACK_API_TOKEN"))

	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: slackexec -channel=CHANNELNAME COMMAND")
		os.Exit(ExitFatal)
	}

	command := args[0]

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to get hostname: %s\n", err)
		os.Exit(ExitFatal)
	}

	osUser, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to get username: %s\n", err)
		os.Exit(ExitFatal)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to get working directory: %s\n", err)
		os.Exit(ExitFatal)
	}

	client.Post(fmt.Sprintf("Running on `%s@%s`:`%s`\n```\n$ %s\n```", osUser.Username, hostname, wd, command))

	cmd := exec.Command("/bin/sh", "-c", args[0])
	buf := new(bytes.Buffer)
	writer := io.MultiWriter(buf, os.Stdout)
	cmd.Stdout = writer
	cmd.Stderr = writer

	exitStatus, err := posixexec.Run(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to exec command: %s\n", err)
		os.Exit(ExitFatal)
	}

	client.Post("Output:\n```\n" + buf.String() + "```")

	os.Exit(exitStatus)
}
