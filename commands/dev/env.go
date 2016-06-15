package dev

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

var (

	// EnvCmd ...
	EnvCmd = &cobra.Command{
		Use:   "evar",
		Short: "Manages environment variables in your local dev app.",
		Long:  ``,
	}

	// EnvAddCmd ...
	EnvAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds environment variable(s) to your dev app.",
		Long: `
Adds environment variable(s) to your dev app. Multiple key-value
pairs can be added simultaneously using a comma-delimited list.
		`,
		Run: envAddFn,
	}

	// EnvListCmd ...
	EnvListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all environment variables registered in your dev app.",
		Long:  ``,
		Run:   envListFn,
	}

	// EnvRemoveCmd ...
	EnvRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Removes environment variable(s) from your dev app.",
		Long: `
Removes environment variable(s) from your dev app. Multiple keys
can be removed simultaneously using a comma-delimited list.
		`,
		Run: envRemoveFn,
	}
)

//
func init() {
	EnvCmd.AddCommand(EnvAddCmd)
	EnvCmd.AddCommand(EnvRemoveCmd)
	EnvCmd.AddCommand(EnvListCmd)
}

// envAddFn ...
func envAddFn(ccmd *cobra.Command, args []string) {
	evars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &evars)
	for _, arg := range args {
		for _, pair := range strings.Split(arg, ",") {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				evars[strings.ToUpper(parts[0])] = parts[1]
			}
		}
	}

	data.Put(config.AppName()+"_meta", "env", evars)
}

// envListFn ...
func envListFn(ccmd *cobra.Command, args []string) {
	evars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &evars)
	fmt.Println(evars)
}

// envRemoveFn ...
func envRemoveFn(ccmd *cobra.Command, args []string) {
	evars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &evars)
	for _, arg := range args {
		for _, key := range strings.Split(arg, ",") {
			delete(evars, strings.ToUpper(key))
		}
	}
	data.Put(config.AppName()+"_meta", "env", evars)
}