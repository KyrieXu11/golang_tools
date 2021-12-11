package git_manager

import (
	"bytes"
	"fmt"
	"github.com/posener/cmd"
	"golang_tools/git_manager/constant"
	"os/exec"
	"strconv"
	"strings"
)

var (
	command    *Command
	flagParser *FlagParser
)

type ICommand interface {
	// list 命令
	list() error
	// use 命令
	use() error
	// set 命令
	set() error
	// get 命令
	get() error
	// delete 命令
	delete() error
}

type GitProfile struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

func printGitProfileList(profiles []*GitProfile) {

}

func (g *GitProfile) print() {

}

func (g *GitProfile) string() string {
	if g.Status == constant.Deleted {
		return ""
	}
	return fmt.Sprintf("|%15d\t|%15s\t|%15s\t|%15s\t|", g.Id, g.Username, g.Email, g.Description)
}

type Command struct {
	flagParser *FlagParser

	Current *GitProfile   `json:"current"`
	Groups  []*GitProfile `json:"groups"`
}

type Args struct {
	command     ICommand
	CommandType int
	Id          int
	Username    string
	Email       string
}

func (o *Args) use() error {
	switch o.CommandType {
	case constant.CommandUse:
		return o.command.use()
	case constant.CommandList:
		return o.command.list()
	case constant.CommandDelete:
		return o.command.delete()
	case constant.CommandGet:
		return o.command.get()
	case constant.CommandSet:
		return o.command.set()
	}
	return nil
}

func (o *Args) copy(args *Args) {
	*o = *args
}

func (o *Args) setCommand(c ICommand) *Args {
	o.command = c
	return o
}

type FlagParser struct{}

func argsParser() *FlagParser {
	if flagParser == nil {
		flagParser = new(FlagParser)
	}
	return flagParser
}

func (o *FlagParser) parse() (*Args, error) {
	var res = new(Args)
	root := cmd.New()
	listSubCmd := root.SubCommand("list", "abc")

	useCommand := o.useCommand(root)
	id := useCommand.Int("id", constant.InvalidId, "")

	setCommand := o.setCommand(root)
	username_p := setCommand.String("username", "", "")
	email_p := setCommand.String("email", "", "")

	// 设置解析规则之后，解析参数
	if err := root.Parse(); err != nil {
		return nil, err
	}

	switch {
	case listSubCmd.Parsed():
		res.CommandType = constant.CommandList
	case useCommand.Parsed():
		res.CommandType = constant.CommandUse
		res.Id = *id
	case setCommand.Parsed():
		res.CommandType = constant.CommandUse
		res.Username = *username_p
		res.Email = *email_p
	}
	return res, nil
}

func (o *FlagParser) useCommand(root *cmd.Cmd) *cmd.SubCmd {
	useSubCmd := root.SubCommand("use", "")
	return useSubCmd
}

func (o *FlagParser) setCommand(root *cmd.Cmd) *cmd.SubCmd {
	return root.SubCommand("set", "")
}

func GitCommand() *Command {
	if command == nil {
		command = &Command{
			flagParser: argsParser(),
		}
	}
	return command
}

func (o *Command) Main() error {
	// 判断是否为git目录
	if err := o.isGitRepo(); err != nil {
		return err
	}

	// 初始化
	if err := o.readFromFile(); err != nil {
		return err
	}
	if args, err := o.flagParser.parse(); err != nil {
		return err
	} else {
		return args.setCommand(o).use()
	}
}

func (o *Command) isGitRepo() error {
	c := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	var out bytes.Buffer
	c.Stdout = &out
	if err := c.Start(); err != nil {
		return err
	}
	if err := c.Wait(); err != nil {
		return err
	}
	s := out.String()
	s = strings.ReplaceAll(s, "\n", "")
	if res, err := strconv.ParseBool(s); err != nil {
		return err
	} else if !res {
		return fmt.Errorf("there is not a git repo")
	}
	return nil
}

func (o *Command) list() error {
	fmt.Println("list")
	return nil
}

// readFromFile 读取文件成对象
func (o *Command) readFromFile() error {
	return nil
}

// flush 刷新内存配置至文件
func (o *Command) flush() {

}

func (o *Command) use() error {
	return nil
}

func (o *Command) delete() error {
	return nil
}

func (o *Command) set() error {
	return nil
}

func (o *Command) get() error {
	return nil
}
