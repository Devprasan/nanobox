package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts your dev platform.",
		Long: `
Starts your dev platform from its previous state. If starting for
the first time, you should also generate a build (nanobox build)
and deploy it into your dev platform (nanobox dev deploy).
		`,
		PreRun: validate.Requires("provider"),
		Run:    devStart,
	}
)

//
// devStart ...
func devStart(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_start", processor.DefaultControl))
}