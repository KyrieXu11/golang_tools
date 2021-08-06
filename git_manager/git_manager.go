package git_manager

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Proxy *NetWorkProxy
	User  *User
}

type NetWorkProxy struct {
	HttpProxy  string
	HttpsProxy string
}

type User struct {
	Name  string
	Email string
}

type ConfigGroup struct {
	GroupName string
	Config    *Config
}

func GetCurrentPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

func IsGitRepo() bool {
	res, err := RunCmd("git", "rev-parse", "--is-inside-work-tree")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if res != "true" {
		fmt.Println(res)
		return false
	}
	return true
}

func RunCmd(name string, arg ...string) (string, error) {
	var err error
	var curPath string
	cmd := exec.Command(name, arg...)
	curPath, err = GetCurrentPath()
	if err != nil {
		return "", err
	}
	cmd.Dir = curPath
	msg, err := cmd.CombinedOutput() // 混合输出stdout+stderr
	if err != nil {
		return "", err
	}
	// 报错时 exit status 1
	return string(msg), err
}
