package main

import "github.com/mitchellh/cli"

// Commands is the mapping of all the available Terraform commands.
var Commands map[string]cli.CommandFactory

var PrimaryCommands []string

var HiddenCommands map[string]struct{}

var Ui cli.Ui
