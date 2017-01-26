package fasta

import (
	"bytes"
	"github.com/fredericlemoine/goalign/align"
)

const (
	FASTA_LINE = 80
)

func min_int(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func WriteAlignment(al align.Alignment) string {
	var buf bytes.Buffer
	al.IterateChar(func(name string, seq []rune) {
		buf.WriteString(">")
		buf.WriteString(name)
		buf.WriteString("\n")
		for i := 0; i < len(seq); i++ {
			if i%FASTA_LINE == 0 && i > 0 {
				buf.WriteString("\n")
			}
			buf.WriteRune(seq[i])
		}
		buf.WriteRune('\n')
	})
	return buf.String()
}

/* This function writes the input alignment as standard fasta sequences
It removes "-" characters.
*/
func WriteSequences(al align.Alignment) string {
	var buf bytes.Buffer

	al.IterateChar(func(name string, seq []rune) {
		buf.WriteString(">")
		buf.WriteString(name)
		buf.WriteString("\n")
		nbchar := 0
		for i := 0; i < len(seq); i++ {
			if seq[i] != '-' {
				buf.WriteRune(seq[i])
				nbchar++
				if nbchar == FASTA_LINE {
					buf.WriteString("\n")
					nbchar = 0
				}
			}
		}
		if nbchar != 0 {
			buf.WriteString("\n")
		}
	})
	return buf.String()
}
