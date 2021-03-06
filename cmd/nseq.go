package cmd

import (
	"fmt"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/goalign/io"
	"github.com/spf13/cobra"
)

// nseqCmd represents the nseq command
var nseqCmd = &cobra.Command{
	Use:   "nseq",
	Short: "Prints the number of sequences in the alignment",
	Long: `Prints the number of sequences in the alignment. 
May take a Phylip of Fasta input alignment.

If the input alignment contains several alignments, will process all of them

Example of usages:

goalign stats nseq -i align.phylip -p
goalign stats nseq -i align.fasta
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if unaligned {
			var seqs align.SeqBag

			if seqs, err = readsequences(infile); err != nil {
				io.LogError(err)
				return
			}
			fmt.Println(seqs.NbSequences())
		} else {
			var aligns *align.AlignChannel

			if aligns, err = readalign(infile); err != nil {
				io.LogError(err)
				return
			}

			for al := range aligns.Achan {
				fmt.Println(al.NbSequences())
			}

			if aligns.Err != nil {
				err = aligns.Err
				io.LogError(err)
			}
		}
		return
	},
}

func init() {
	nseqCmd.PersistentFlags().BoolVar(&unaligned, "unaligned", false, "Considers sequences as unaligned and format fasta (phylip, nexus,... options are ignored)")
	statsCmd.AddCommand(nseqCmd)
}
