package main

import (
	"golang_tools/git_manager"
	"log"
)

func main() {
	if err := git_manager.GitCommand().Main(); err != nil {
		log.Fatalln(err)
	}
}
