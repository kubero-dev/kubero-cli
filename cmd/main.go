/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/

package main

import (
	"os"
)

func main() {
	kbr := NewKuberoClient()
	if err := kbr.Command().Execute(); err != nil {
		kbr.log.Error(err.Error(), map[string]interface{}{
			"context": "kubero-cli",
			"action":  "Execute",
			"error":   err.Error(),
		})
		os.Exit(1)
	}
}
