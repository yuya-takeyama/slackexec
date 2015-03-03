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
var hostname string
var osUser *user.User
var cwd string
var client *Slack

const (
	ExitFatal = 111
)

func init() {
	var err error

	flag.StringVar(&channel, "channel", "#general", "channel to post message")
	flag.StringVar(&name, "name", "slackexec", "username of the bot")
	flag.StringVar(&icon, "icon", ":computer:", "icon of the bot")
	flag.Parse()

	client = NewSlack(name, icon, channel, os.Getenv("SLACK_API_TOKEN"))

	hostname, err = os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to get hostname: %s\n", err)
		os.Exit(ExitFatal)
	}

	osUser, err = user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to get username: %s\n", err)
		os.Exit(ExitFatal)
	}

	cwd, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to get working directory: %s\n", err)
		os.Exit(ExitFatal)
	}
}

func main() {
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: slackexec -channel=CHANNELNAME COMMAND")
		os.Exit(ExitFatal)
	}

	command := args[0]

	client.Post(fmt.Sprintf("Running on `%s@%s`:`%s`\n```\n$ %s\n```", osUser.Username, hostname, cwd, command))

	cmd, buf := execCommand(args[0])
	exitStatus, err := posixexec.Run(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "slackexec: failed to exec command: %s\n", err)
		os.Exit(ExitFatal)
	}

	client.Post("Output:\n```\n" + buf.String() + "```")

	os.Exit(exitStatus)
}

func execCommand(command string) (*exec.Cmd, *bytes.Buffer) {
	cmd := exec.Command("/bin/sh", "-c", command)
	buf := new(bytes.Buffer)
	writer := io.MultiWriter(buf, os.Stdout)
	cmd.Stdout = writer
	cmd.Stderr = writer

	return cmd, buf
}
