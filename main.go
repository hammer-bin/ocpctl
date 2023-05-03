package main

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	"github.com/mitchellh/cli"
	"log"
	"ocpctl/internal/terminal"
	"ocpctl/version"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvCLI = "TF_CLI_ARGS"
)

type ui struct {
	cli.Ui
}

func (u *ui) Warn(msg string) {
	u.Ui.Output(msg)
}

func init() {
	Ui = &ui{
		&cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
			Reader:      os.Stdin,
		},
	}
}

func main() {
	os.Exit(realMain())
}

func realMain() int {

	fmt.Println("Start ocpctl")

	log.Printf("[INFO] Container Platform Version %s", version.Version)

	streams, err := terminal.Init()
	if err != nil {
		Ui.Error(fmt.Sprintf("Failed to configure the terminal: %s", err))
		return 1
	}
	if streams.Stdout.IsTerminal() {
		log.Printf("[TRACE] Stdout is a terminal of width %d", streams.Stdout.Columns())
	} else {
		log.Printf("[TRACE] Stdout is not a terminal")
	}
	if streams.Stderr.IsTerminal() {
		log.Printf("[TRACE] Stderr is a terminal of width %d", streams.Stderr.Columns())
	} else {
		log.Printf("[TRACE] Stderr is not a terminal")
	}
	if streams.Stdin.IsTerminal() {
		log.Printf("[TRACE] Stdin is a terminal")
	} else {
		log.Printf("[TRACE] Stdin is not a terminal")
	}

	// Get the command line args.
	binName := filepath.Base(os.Args[0])
	args := os.Args[1:]
	fmt.Println(binName)

	originalWd, err := os.Getwd()
	if err != nil {
		Ui.Error(fmt.Sprintf("Failed to determine current working directory: %s", err))
		return 1
	}
	fmt.Println(originalWd)

	if Commands == nil {
		//initCommand()
	}

	cliRunner := &cli.CLI{
		Name:       binName,
		Args:       args,
		Commands:   Commands,
		HelpFunc:   helpFunc,
		HelpWriter: os.Stdout,
	}

	args, err = mergeEnvArgs(EnvCLI, cliRunner.Subcommand(), args)
	if err != nil {
		Ui.Error(err.Error())
		return 1
	}

	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	exitCode, err := cliRunner.Run()
	if err != nil {
		Ui.Error(fmt.Sprintf("Error executing CLI: %s", err.Error()))
		return 1
	}

	return exitCode
}

func mergeEnvArgs(envName string, cmd string, args []string) ([]string, error) {
	v := os.Getenv(envName)
	if v == "" {
		return args, nil
	}

	log.Printf("[INFO] %s value: %q", envName, v)
	extra, err := shellwords.Parse(v)
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing extra CLI args from %s: %s", envName, err)
	}

	search := cmd
	if idx := strings.LastIndex(search, " "); idx >= 0 {
		search = cmd[idx+1:]
	}

	idx := -1
	for i, v := range args {
		if v == search {
			idx = i
			break
		}
	}

	idx++

	newArgs := make([]string, len(args)+len(extra))
	copy(newArgs, args[:idx])
	copy(newArgs[idx:], extra)
	copy(newArgs[len(extra)+idx:], args[idx:])
	return newArgs, nil
}
