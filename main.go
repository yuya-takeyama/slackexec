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
var printVersion bool
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
	flag.BoolVar(&printVersion, "version", false, "print version")
	flag.Parse()

	if (printVersion) {
		fmt.Fprintf(os.Stderr, "%s version %s, build %s\n", Name, Version, GitCommit)
		os.Exit(0)
	}

	client = NewSlack(name, icon, channel, os.Getenv("SLACK_API_TOKEN"))

	hostname, err = os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to get hostname: %s\n", Name, err)
		os.Exit(ExitFatal)
	}

	osUser, err = user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to get username: %s\n", Name, err)
		os.Exit(ExitFatal)
	}

	cwd, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to get working directory: %s\n", Name, err)
		os.Exit(ExitFatal)
	}
}

func main() {
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s -channel=CHANNELNAME COMMAND\n", Name)
		os.Exit(ExitFatal)
	}

	command := args[0]

	client.Post(fmt.Sprintf("Running on `%s@%s`:`%s`\n```\n$ %s\n```", osUser.Username, hostname, cwd, command))

	cmd, buf := execCommand(args[0])
	exitStatus, err := posixexec.Run(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to exec command: %s\n", Name, err)
		os.Exit(ExitFatal)
	}

	client.Post("Output:\n```\n" + buf.String() + "```")

	os.Exit(exitStatus)
}

func execCommand(command string) (*exec.Cmd, *bytes.Buffer) {
	cmd := exec.Command(os.Getenv("SHELL"), "-c", command)
	buf := new(bytes.Buffer)
	writer := io.MultiWriter(buf, os.Stdout)
	cmd.Stdout = writer
	cmd.Stderr = writer

	return cmd, buf
}
