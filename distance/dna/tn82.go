//+ build ignore

package dna

import (
	"fmt"
	"math"

	"github.com/evolbioinfo/goalign/align"
)

type TN82Model struct {
	/* Vector of nt proba */
	pi            []float64 // proba of each nt
	numSites      float64   // Number of selected sites (no gaps)
	selectedSites []bool    // true for selected sites
	removegaps    bool      // If true, we will remove posision with >=1 gaps
	sequenceCodes [][]uint8 // Sequences converted into int codes
}

func NewTN82Model(removegaps bool) *TN82Model {
	return &TN82Model{
		nil,
		0,
		nil,
		removegaps,
		nil,
	}
}

// Distance computes TN82 distance between 2 sequences
func (m *TN82Model) Distance(seq1 []uint8, seq2 []uint8, weights []float64) (float64, error) {
	diff, total := countDiffs(seq1, seq2, m.selectedSites, weights, false)
	diff = diff / total

	psi := init2DFloat(4, 4)
	totalPairs, err := countNtPairs2Seq(seq1, seq2, m.selectedSites, weights, psi)
	if err != nil {
		return 0.0, err
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			psi[i][j] = psi[i][j] / totalPairs
		}
	}
	denom := 0.0
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			denom += psi[i][j] * psi[i][j] / (2 * m.pi[i] * m.pi[j])
		}
	}
	b1 := diff * diff / denom
	dist := -1.0 * b1 * math.Log(1.0-diff/b1)
	if dist > 0 {
		return dist, nil
	} else {
		return 0, nil
	}
}

func (m *TN82Model) InitModel(al align.Alignment, weights []float64, gamma bool, alpha float64) (err error) {
	m.numSites, m.selectedSites = selectedSites(al, weights, m.removegaps)
	if m.sequenceCodes, err = alignmentToCodes(al); err != nil {
		return
	}
	m.pi, err = probaNt(m.sequenceCodes, m.selectedSites, weights)
	return
}

// Sequence returns the ith sequence of the alignment
// encoded in int
func (m *TN82Model) Sequence(i int) (seq []uint8, err error) {
	if i < 0 || i >= len(m.sequenceCodes) {
		err = fmt.Errorf("This sequence does not exist: %d", i)
		return
	}
	seq = m.sequenceCodes[i]
	return
}
