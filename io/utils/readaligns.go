package utils

import (
	"bufio"
	"io"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/goalign/io/clustal"
	"github.com/evolbioinfo/goalign/io/fasta"
	"github.com/evolbioinfo/goalign/io/nexus"
	"github.com/evolbioinfo/goalign/io/phylip"
)

// Parses the input buffer while automatically
// detecting the format between Newick, Phylip, and Nexus
//
// If several alignments are present in the onput file, only the first will be
// parsed.
//
// Returned format may be align.FORMAT_PHYLIP, align.FORMAT_FASTA, or align.FORMAT_NEXUS
//
// rootinpustrict: In the case of phylip detected format: should we consider it as strict or not?
//
// There is no new go routine here because only 1 alignment is parsed. No need to give a closer.
// If the reader comes from a file, the file must be closed in the calling function.
func ParseAlignmentAuto(r *bufio.Reader, rootinputstrict bool) (al align.Alignment, format int, err error) {
	var firstbyte byte

	if firstbyte, err = r.ReadByte(); err != nil {
		return
	}

	if err = r.UnreadByte(); err != nil {
		return
	}

	// First test Fasta format
	if firstbyte == '>' {
		format = align.FORMAT_FASTA
		al, err = fasta.NewParser(r).Parse()
	} else if firstbyte == '#' {
		if al, err = nexus.NewParser(r).Parse(); err != nil {
			return
		}
		format = align.FORMAT_NEXUS
	} else if firstbyte == 'C' {
		if al, err = clustal.NewParser(r).Parse(); err != nil {
			return
		}
		format = align.FORMAT_CLUSTAL
	} else {
		// Finally test Phylip
		format = align.FORMAT_PHYLIP
		al, err = phylip.NewParser(r, rootinputstrict).Parse()
	}
	return
}

// Parses the input buffer while automatically
// detecting the format between Newick, Phylip, and Nexus
//
// If several alignments are present in the input file, they are queued in the channel
//
// rootinpustrict: In the case of phylip detected format: should we consider it as strict or not?
//
// If there is something to close ( f!=nil) after the parsing (like input file, etc.), f will be closed
// after parsing is finished (even in the go routine in the case of several input alignments).
// If the alignment comes from a file for exemple, the file will be closed by this function, so no need to
// do it in the calling function
func ParseMultiAlignmentsAuto(f io.Closer, r *bufio.Reader, rootinputstrict bool) (alchan *align.AlignChannel, format int, err error) {
	var al align.Alignment
	var firstbyte byte

	alchan = &align.AlignChannel{}

	if firstbyte, err = r.ReadByte(); err != nil {
		return
	}

	if err = r.UnreadByte(); err != nil {
		return
	}
	// First test Fasta format
	if firstbyte == '>' {
		if al, err = fasta.NewParser(r).Parse(); err != nil {
			return
		}
		format = align.FORMAT_FASTA
		alchan.Achan = make(chan align.Alignment, 1)
		alchan.Achan <- al
		if f != nil {
			f.Close()
		}
		close(alchan.Achan)
	} else if firstbyte == '#' {
		if al, err = nexus.NewParser(r).Parse(); err != nil {
			return
		}
		format = align.FORMAT_NEXUS
		alchan.Achan = make(chan align.Alignment, 1)
		alchan.Achan <- al
		if f != nil {
			f.Close()
		}
		close(alchan.Achan)
	} else if firstbyte == 'C' {
		if al, err = clustal.NewParser(r).Parse(); err != nil {
			return
		}
		format = align.FORMAT_CLUSTAL
		alchan.Achan = make(chan align.Alignment, 1)
		alchan.Achan <- al
		if f != nil {
			f.Close()
		}
		close(alchan.Achan)
	} else {
		format = align.FORMAT_PHYLIP
		// Finally test Phylip
		alchan.Achan = make(chan align.Alignment, 15)
		go func() {
			phylip.NewParser(r, rootinputstrict).ParseMultiple(alchan)
			if f != nil {
				f.Close()
			}
		}()
	}
	return
}
