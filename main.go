package main

import (
	command "github.com/ankitsheoran1/grep/command"
	"github.com/ankitsheoran1/grep/locate"
)

func main() {

	args, ok := command.ParserArgs()

	if !ok {
		return
	}
	locator := locate.NewLocator(args.Directory)
	locator.Options = args.Options
	locator.Dig(args.SearchString)
}


