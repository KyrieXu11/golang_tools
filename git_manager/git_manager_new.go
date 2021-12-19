package git_manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/tablewriter"
	"github.com/lensesio/tableprinter"
	"github.com/posener/cmd"
	"golang_tools/git_manager/common/constant"
	"golang_tools/git_manager/common/utils"
	"io/fs"
	"io/ioutil"
	"os"
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
	use(id int) error
	// set 命令
	set() error
	// get 命令
	get() error
	// delete 命令
	delete(id int) error
}

type GitProfile struct {
	Id          int    `json:"id" header:"id"`
	Username    string `json:"username" header:"username"`
	Email       string `json:"email" header:"email"`
	Description string `json:"description" header:"description"`
	Status      int    `json:"status" header:"status"`
}

func printGitProfileList(profiles []*GitProfile) {
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.HeaderBgColor = tablewriter.BgBlackColor
	printer.HeaderFgColor = tablewriter.FgGreenColor

	printer.Print(profiles)
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
	savePath   string
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
		if err := o.checkId(); err != nil {
			return err
		}
		return o.command.use(o.Id)
	case constant.CommandList:
		return o.command.list()
	case constant.CommandDelete:
		if err := o.checkId(); err != nil {
			return err
		}
		return o.command.delete(o.Id)
	case constant.CommandGet:
		return o.command.get()
	case constant.CommandSet:
		return o.command.set()
	}
	return nil
}

func (o *Args) checkId() error {
	if o.Id == constant.InvalidId {
		return fmt.Errorf("please input valid id")
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

func newFlagParser() *FlagParser {
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
		var groups []*GitProfile
		var current = new(GitProfile)
		groups = append(groups, current)
		command = &Command{
			savePath:   `./profiles.json`,
			flagParser: newFlagParser(),
			Current:    current,
			Groups:     groups,
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

// readFromFile 读取文件成对象
func (o *Command) readFromFile() error {
	file, err := utils.ReadFile(o.savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tempObj = new(Command)
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(bs) == 0 {
		*tempObj = *o
		return tempObj.flush()
	}
	if err = json.Unmarshal(bs, &tempObj); err != nil {
		return err
	}
	o.Current = tempObj.Current
	o.Groups = tempObj.Groups
	return nil
}

// flush 刷新内存配置至文件
func (o *Command) flush() error {
	f, err := utils.ReadFile(o.savePath)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(o.savePath, data, fs.ModeExclusive); err != nil {
		return err
	}
	return nil
}

func (o *Command) list() error {
	printGitProfileList(o.Groups)
	return nil
}

func (o *Command) use(id int) error {
	var err = fmt.Errorf("please input valid id")
	if o.Current.Id == id {
		return nil
	}
	for _, group := range o.Groups {
		if group.Id == id {
			*o.Current = *group
			return o.flush()
		}
	}
	return err
}

func (o *Command) delete(id int) error {
	if o.Current.Id == id {
		return fmt.Errorf("chose profile id:%d is in use", id)
	}
	length := len(o.Groups)
	for i, group := range o.Groups {
		if group.Id == id {
			for j := i; j < length-1; j++ {
				o.Groups[j] = o.Groups[j+1]
			}
		}
	}
	return o.flush()
}

func (o *Command) set() error {
	return o.flush()
}

func (o *Command) get() error {
	return nil
}
