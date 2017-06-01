# Goalign: toolkit and api for alignment manipulation

## Commands

### random
This command generate a random alignment with uniform distribution of nucleotides or amino acids. It is intended for testing purpose, as no evolutionary information is taken into account.

#### Usage
```
Usage:
  goalign random [flags]

Flags:
  -a, --amino-acids        Aminoacid sequences (otherwise, nucleotides)
  -l, --length int         Length of sequences to generate (default 100)
  -n, --nb-seqs int        Number of sequences to generate (default 10)
  -o, --out-align string   Random alignment output file (default "stdout")
  -s, --seed int           Initial Random Seed (default 1495468802031837781)

Global Flags:
  -i, --align string   Alignment input file (default "stdin")
  -p, --phylip         Alignment is in phylip? False=Fasta
```

#### Examples

* Generating a random alignment with 100 sequences and 1000 nucleotides:
```
goalign random -n 100 -l 1000 -s 10 | goalign stats
```

Should give the following statistics:
```
length	1000
nseqs	100
avgalleles	4.0000
char	nb	freq
A	24899	0.248990
C	25032	0.250320
G	24888	0.248880
T	25181	0.251810
```