package align

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestRandomAlignment(t *testing.T) {
	length := 3000
	nbseqs := 500
	a, err := RandomAlignment(AMINOACIDS, length, nbseqs)
	if err != nil {
		t.Error(err)
	}

	if a.Length() != length {
		t.Error(fmt.Sprintf("Length should be %d and is %d", length, a.Length()))
	}
	if a.NbSequences() != nbseqs {
		t.Error(fmt.Sprintf("Nb sequences should be %d and is %d", nbseqs, a.NbSequences()))
	}
}

func TestAppendIdentifier(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 50)
	if err != nil {
		t.Error(err)

	}
	a.AppendSeqIdentifier("IDENT", false)

	a.IterateChar(func(name string, sequence []rune) {
		if !strings.HasPrefix(name, "IDENT") {
			t.Error("Sequence name does not start with expected id: IDENT")
		}
	})

	a.AppendSeqIdentifier("IDENT", true)
	a.IterateChar(func(name string, sequence []rune) {
		if !strings.HasSuffix(name, "IDENT") {
			t.Error("Sequence name does not end with expected id: IDENT")
		}
	})
}

func TestRemoveOneGaps(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 300)
	if err != nil {
		t.Error(err)

	}

	/* We add 1 gap per site */
	pos := 0
	a.IterateChar(func(name string, sequence []rune) {
		sequence[pos] = GAP
		pos++
	})

	a.RemoveGaps(0.0)

	if a.Length() != 0 {
		t.Error("We should have removed all positions")
	}
	a.IterateChar(func(name string, sequence []rune) {
		if len(sequence) != 0 {
			t.Error(fmt.Sprintf("Sequence length after removing gaps should be 0 and is : %d", len(sequence)))
		}
	})
}

func TestRemoveAllGaps(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 300)
	if err != nil {
		t.Error(err)

	}

	backupseq := make([]rune, 0, 300)
	seq0, found := a.GetSequenceChar(0)
	if !found {
		t.Error("Problem finding first sequence")
	}

	/* We add all gaps on 1 site */
	/* And one gap at all sites */
	pos1 := 20
	pos2 := 0
	a.IterateChar(func(name string, sequence []rune) {
		sequence[pos1] = GAP
		sequence[pos2] = GAP
		pos2++
	})
	backupseq = append(backupseq, seq0...)
	/* Remove position 20 */
	backupseq = append(backupseq[:20], backupseq[21:]...)

	a.RemoveGaps(1.0)

	if a.Length() != 299 {
		t.Error("We should have removed only one position")
	}

	a.IterateChar(func(name string, sequence []rune) {
		if len(sequence) != 299 {
			t.Error(fmt.Sprintf("Sequence length after removing gaps should be 299 and is : %d", len(sequence)))
		}
	})

	newseq, found2 := a.GetSequenceChar(0)
	if !found2 {
		t.Error("Problem finding first seqence")
	}

	for i, c := range newseq {
		if c != backupseq[i] {
			t.Error(fmt.Sprintf("Char at position %d should be %c and is %c", i, backupseq[i], c))
		}
	}
}

func TestClone(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 300)
	if err != nil {
		t.Error(err)

	}

	/* We add 1 gap per site */
	pos := 0
	a.IterateChar(func(name string, sequence []rune) {
		sequence[pos] = GAP
		pos++
	})

	a2, err2 := a.Clone()
	if err2 != nil {
		t.Error(err2)
	}

	a.RemoveGaps(0.0)

	a2.IterateChar(func(name string, sequence []rune) {
		if len(sequence) != 300 {
			t.Error(fmt.Sprintf("Clone lenght should be 300 and is : %d", len(sequence)))
		}
	})
}

func TestClone2(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 300)
	if err != nil {
		t.Error(err)

	}

	a2, err2 := a.Clone()
	if err2 != nil {
		t.Error(err2)
	}

	i := 0
	a2.IterateChar(func(name string, sequence []rune) {
		s2, ok := a.GetSequenceChar(i)
		n2, ok2 := a.GetSequenceName(i)

		if !ok || !ok2 {
			t.Error(fmt.Sprintf("Sequence not found in clone alignment: %s", name))
		}

		if len(sequence) != len(s2) {
			t.Error(fmt.Sprintf("Clone length is different from original length : %d != %d", len(sequence), len(s2)))
		}
		if name != n2 {
			t.Error(fmt.Sprintf("Clone and original sequences at position %d have different names : %s != %s", name, n2))
		}
		for j, c := range sequence {
			if c != s2[j] {
				t.Error(fmt.Sprintf("Clone sequence is different from original at position %d : %c != %c", j, c, s2[j]))
			}
		}
		i++
	})
}

func TestAvgAlleles(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 300)
	if err != nil {
		t.Error(err)

	}

	a.IterateChar(func(name string, sequence []rune) {
		for j, _ := range sequence {
			sequence[j] = 'A'
		}
	})

	if a.AvgAllelesPerSite() != 1 {
		t.Error("There should be 1 allele per site in this alignment")
	}
}

func TestAvgAlleles2(t *testing.T) {
	a, err := RandomAlignment(AMINOACIDS, 300, 300)
	if err != nil {
		t.Error(err)

	}

	i := 0
	a.IterateChar(func(name string, sequence []rune) {
		for j, _ := range sequence {
			/* One only gap sequence */
			if i == 10 {
				sequence[j] = GAP
			} else if i <= 75 {
				sequence[j] = 'A'
			} else if i <= 150 {
				sequence[j] = 'C'
			} else if i <= 225 {
				sequence[j] = 'G'
			} else {
				sequence[j] = 'T'
			}
		}
		// Add a gap at a whole position
		sequence[50] = GAP
		i++
	})
	fmt.Println(a.AvgAllelesPerSite())
	if a.AvgAllelesPerSite() != 4 {
		t.Error("There should be 4 allele per site in this alignment")
	}
}

func TestRename(t *testing.T) {
	length := 3000
	nbseqs := 500
	a, err := RandomAlignment(AMINOACIDS, length, nbseqs)

	if err != nil {
		t.Error(err)
	}

	// Map File
	namemap1 := make(map[string]string)
	namemap2 := make(map[string]string)
	for i := 0; i < 500; i++ {
		namemap1[fmt.Sprintf("Seq%04d", i)] = fmt.Sprintf("New%04d", i)
		namemap2[fmt.Sprintf("New%04d", i)] = fmt.Sprintf("Seq%04d", i)
	}

	a.Rename(namemap1)
	i := 0
	a.IterateChar(func(name string, sequence []rune) {
		expected := fmt.Sprintf("New%04d", i)
		if name != expected {
			t.Error(fmt.Sprintf("Sequence name should be %s and is %s", expected, name))
		}
		i++
	})

	a.Rename(namemap2)
	i = 0
	a.IterateChar(func(name string, sequence []rune) {
		expected := fmt.Sprintf("Seq%04d", i)
		if name != expected {
			t.Error(fmt.Sprintf("Sequence name should be %s and is %s", expected, name))
		}
		i++
	})
}

func TestRogue(t *testing.T) {
	length := 3000
	nbseqs := 500
	nrogue := 0.5
	a, err := RandomAlignment(AMINOACIDS, length, nbseqs)

	if err != nil {
		t.Error(err)
	}

	a2, err2 := a.Clone()
	if err2 != nil {
		t.Error(err2)
	}

	rogues, intacts := a2.SimulateRogue(nrogue, 1.0)

	if (len(rogues) + len(intacts)) != a.NbSequences() {
		t.Error("Number of intact + rogue sequences is not the same than the total number of sequences")
	}

	for _, s := range rogues {
		seq2, ok2 := a2.GetSequence(s)
		seq, ok := a.GetSequence(s)
		if !ok || !ok2 {
			t.Error("Rogue name does not exist in alignment...")
		}
		if seq == seq2 {
			t.Error("Rogue sequence is the same after simulation (should be shuffled)...")
		}
	}
	for _, s := range intacts {
		seq2, ok2 := a2.GetSequence(s)
		seq, ok := a.GetSequence(s)
		if !ok || !ok2 {
			t.Error("Intact name does not exist in alignment...")
		}
		if seq != seq2 {
			t.Error("Intact sequences should be he same before and after simulation...")
		}
	}
}

func TestRogue2(t *testing.T) {
	length := 3000
	nbseqs := 500
	nrogue := 0.0
	a, err := RandomAlignment(AMINOACIDS, length, nbseqs)

	if err != nil {
		t.Error(err)
	}

	a2, err2 := a.Clone()
	if err2 != nil {
		t.Error(err2)
	}

	rogues, intacts := a2.SimulateRogue(nrogue, 1.0)

	if (len(rogues) + len(intacts)) != a.NbSequences() {
		t.Error("Number of intact + rogue sequences is not the same than the total number of sequences")
	}

	if len(rogues) > 0 {
		t.Error("There should be no rogue taxa in output")
	}

	for _, s := range intacts {
		seq2, ok2 := a2.GetSequence(s)
		seq, ok := a.GetSequence(s)
		if !ok || !ok2 {
			t.Error("Intact name does not exist in alignment...")
		}
		if seq != seq2 {
			t.Error("Intact sequences should be he same before and after simulation...")
		}
	}
}

func TestRogue3(t *testing.T) {
	length := 3000
	nbseqs := 500
	nrogue := 1.0
	a, err := RandomAlignment(AMINOACIDS, length, nbseqs)

	if err != nil {
		t.Error(err)
	}

	a2, err2 := a.Clone()
	if err2 != nil {
		t.Error(err2)
	}

	rogues, intacts := a2.SimulateRogue(nrogue, 1.0)

	if (len(rogues) + len(intacts)) != a.NbSequences() {
		t.Error(fmt.Sprintf("Number of intact (%d) + rogue (%d) sequences is not the same than the total number of sequences (%d)", len(intacts), len(rogues), a.NbSequences()))
	}

	if len(rogues) < a.NbSequences() {
		t.Error("All sequences should be rogues")
	}

	if len(intacts) > 0 {
		t.Error("There should be no intact sequences")
	}
	for _, s := range rogues {
		seq2, ok2 := a2.GetSequence(s)
		seq, ok := a.GetSequence(s)
		if !ok || !ok2 {
			t.Error("Rogue name does not exist in alignment...")
		}
		if seq == seq2 {
			t.Error("Rogue sequence is the same after simulation (should be shuffled)...")
		}
	}
}

func TestEntropy(t *testing.T) {
	length := 3
	nbseqs := 5
	a, err := RandomAlignment(AMINOACIDS, length, nbseqs)

	alldifferent := []rune{'A', 'R', 'N', 'D', 'C'}
	// First site: only 'R' => Entropy 0.0
	// Second site: Only different aminoacids => Entropy 1.0
	for i := 0; i < 5; i++ {
		a.SetSequenceChar(i, 0, 'R')
		a.SetSequenceChar(i, 1, alldifferent[i])
	}
	e, err := a.Entropy(0, false)
	if err != nil {
		t.Error(err)
	}
	if e != 0.0 {
		t.Error(fmt.Sprintf("Entropy should be 0.0 and is %f", e))
	}

	e, err = a.Entropy(1, false)
	if err != nil {
		t.Error(err)
	}
	expected := -1.0 * float64(nbseqs) * 1.0 / float64(nbseqs) * math.Log(1.0/float64(nbseqs))
	if int(e*10000000) != int(expected*10000000) {
		t.Error(fmt.Sprintf("Entropy should be %.7f and is %.7f", expected, e))
	}
}