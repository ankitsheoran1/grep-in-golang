package command

import (
	"os"
    "fmt"
	"github.com/ankitsheoran1/grep/locate"
	"github.com/fatih/color"
)

type Arguments struct {
	SearchString string
	Directory    string
	Options      locate.OptionConfig
}

type Parser interface {
	ParseArgs() (Arguments, bool)
	ParseOpts(opts []string)locate.OptionConfig
}

// ParseArgs uses os.Args and parses the flags provided. If the flags(n) given are 0=n<2 
// then we print the usage to stdout, else we parse, assert and return the args
func ParserArgs() (Arguments, bool) {
	if len(os.Args) <= 2 {
		Usage()
		os.Exit(0)
	}
	args := os.Args[1:]
	opts := args[2:]
	mainArgs := Arguments{}
	if len(opts) != 0 {
		mainArgs.Options = ParseOpts(opts)
	}

	if args[0] == "usage" {
		Usage()
		os.Exit(0)
		return Arguments{}, false
	}

	if len(args) < 2 {
		Usage()
		return Arguments{}, false
	}
	mainArgs.SearchString, mainArgs.Directory = args[0], args[1]
	ok := validateDir(mainArgs.Directory)
	if !ok {
		fmt.Printf("%s is not a valid directory\n", mainArgs.Directory)
		return Arguments{}, false
	}
	return mainArgs, true
}


func ParseOpts(opts []string) locate.OptionConfig {
	var options locate.OptionConfig
	for _, opt := range opts {
		switch opt {
		case "-h": 
			options.Hidden = true
	    case "-v": 
			options.Verbose = true
		default:
			color.Red("Unrecognized argument %s\n", opt)		
		}
	}
	return options

}


func Usage() {
	fmt.Println(
		`
Usage:

grip <searchString> ( <searchDir> | . ) [-opt]

Arguments:

        <searchString>    The desired text you want to search for

        <searchDir>       The directory in which you'd like to search. Use '.' to search in the current directory

Options:

        -h                                Search hidden folders and files

        `,
	)

}


func validateDir(dir string) bool {
	_, err := os.ReadDir(dir)

	return err == nil
}