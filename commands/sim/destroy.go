package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

// DestroyCmd ...
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys the docker machines associated with your sim app.",
	Long: ``,
	PreRun: validate.Requires("provider"),
	Run:    destroyFn,
}

// destroyFn ...
func destroyFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("sim_destroy", processor.DefaultControl))
}