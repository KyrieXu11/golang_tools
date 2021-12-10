package git_manager

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var (
	command *Command
)

type GitProfile struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

type Command struct {
	Current *GitProfile

	Groups []*GitProfile
}

func GitCommand() *Command {
	if command != nil {
		command = &Command{}
	}
	return command
}

func readConfFile() {

}

func (o *Command) isGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	s := out.String()
	s = strings.ReplaceAll(s, "\n", "")
	if res, err := strconv.ParseBool(s); err != nil {
		return err
	} else if !res {
		return fmt.Errorf("there is no a git repo")
	}
	return nil
}

func (o *Command) List() {
	if err := o.isGitRepo(); err != nil {
		log.Fatalln(err)
	}

}

func (o *Command) flush() {

}
