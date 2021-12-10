package git_manager

//
// import (
// 	"flag"
// 	"fmt"
// 	"gopkg.in/yaml.v2"
// 	"io/ioutil"
// 	"os"
// 	"os/exec"
// 	"strings"
// )
//
// /*
// [https]
//        proxy = https://127.0.0.1:7890
// [http]
//        proxy = http://127.0.0.1:7890
// */
//
// var (
// 	ConfigList []*ConfigGroup
// 	conf       *GMConfig
// )
//
// type GMConfig struct {
// 	path string `yaml:"file-path"`
// }
//
// type Config struct {
// 	Http  *HTTP
// 	Https *HTTPS
// 	User  *User
// }
//
// type NetWork struct {
// 	Proxy string
// }
//
// type HTTP struct {
// 	NetWork
// }
//
// type HTTPS struct {
// 	NetWork
// }
//
// type User struct {
// 	Name  string
// 	Email string
// }
//
// type ConfigGroup struct {
// 	GroupName string
// 	Config    *Config
// }
//
// func init() {
// 	bs, err := ioutil.ReadFile("etc/gm_config.yml")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	if err = yaml.Unmarshal(bs, conf); err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// }
//
// func Add(groupName string, config *Config) {
// 	set(groupName, config)
// }
//
// func set(groupName string, config *Config) {
// 	validation := IsGitRepo()
// 	if !validation {
// 		return
// 	}
// 	for _, group := range ConfigList {
// 		if group.GroupName == groupName {
// 			group.set(groupName, config)
// 			return
// 		}
// 	}
// 	group := &ConfigGroup{
// 		GroupName: groupName,
// 		Config:    config,
// 	}
// 	ConfigList = append(ConfigList, group)
// }
//
// func Sync() {
//
// }
//
// func (c *ConfigGroup) set(groupName string, config *Config) {
// 	if c.GroupName != "" {
// 		if c.GroupName != groupName {
// 			fmt.Println("group name is not exist")
// 			return
// 		}
// 	} else {
// 		c.GroupName = groupName
// 	}
// 	*(c.Config) = *config
// }
//
// func GetCurrentPath() (string, error) {
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.Replace(dir, "\\", "/", -1), nil
// }
//
// func IsGitRepo() bool {
// 	res, err := RunCmd("git", "rev-parse", "--is-inside-work-tree")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return false
// 	}
// 	if res != "true" {
// 		fmt.Println(res)
// 		return false
// 	}
// 	return true
// }
//
// func RunCmd(name string, arg ...string) (string, error) {
// 	var err error
// 	var curPath string
// 	cmd := exec.Command(name, arg...)
// 	curPath, err = GetCurrentPath()
// 	if err != nil {
// 		return "", err
// 	}
// 	cmd.Dir = curPath
// 	msg, err := cmd.CombinedOutput() // 混合输出stdout+stderr
// 	if err != nil {
// 		return "", err
// 	}
// 	// 报错时 exit status 1
// 	return string(msg), err
// }
//
// func Main(f *flag.Flag) {
//
// 	defer Sync()
// }
