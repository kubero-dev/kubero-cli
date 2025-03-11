/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	"github.com/faelmori/kubero-cli/cmd/cli"
	"github.com/faelmori/kubero-cli/internal/config"
)

func init() {

	cfg := config.NewViperConfig("", "config")
	_ = cfg.LoadConfig()
}

func main() {
	cli.Execute()
}

//
