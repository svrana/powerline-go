package main

import (
	"fmt"
	"os"
)

func segmentUser(p *powerline) {
	var userPrompt string
	if *p.args.Shell == "bash" {
		userPrompt = fmt.Sprintf("\\e[1m\\u")
	} else if *p.args.Shell == "zsh" {
		userPrompt = "%n"
	} else {
		user, _ := os.LookupEnv("USER")
		userPrompt = user
	}

	var background uint8
	if os.Getuid() == 0 {
		background = p.theme.UsernameRootBg
	} else {
		background = p.theme.UsernameBg
	}

	p.appendSegment("user", segment{
		content:    userPrompt,
		foreground: p.theme.UsernameFg,
		background: background,
	})
}
