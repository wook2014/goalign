package cmd

import (
	"github.com/spf13/cobra"
)

// mutateCmd represents the mutate command
var gapsCmd = &cobra.Command{
	Use:   "snvs",
	Short: "Adds substitutions uniformly in an input alignment",
	Long: `Adds substitutions uniformly in an input alignment.

      if rate <= 0 : does nothing
      if rate > 1  : then rate = 1
`,
	Run: func(cmd *cobra.Command, args []string) {
		aligns := readalign(infile)
		f := openWriteFile(mutateOutput)
		for al := range aligns.Achan {
			al.Mutate(mutateRate)
			writeAlign(al, f)
		}
		f.Close()
	},
}

func init() {
	mutateCmd.AddCommand(gapsCmd)
}
