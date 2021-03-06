# Goalign: toolkit and api for alignment manipulation

## Commands

### diff
Takes an input alignment, and compares all the sequences to the first one.
Any character that is identical to the reference sequence is replaced with ".".
Prints only characters that are different from the first sequence of the alignment.

If option `--counts` is given, then the output is not an alignment but a file containing for each sequence
the number of occurences of each difference with the reference sequence.
The format is tab separated, with following columns:

- Sequence name (reference sequence is not included)
- For each unique difference, its number of occurences

#### Usage
```
Usage:
  goalign diff [flags]

Flags:
  -h, --help            help for diff
  -o, --output string   Diff output file (default "stdout")

Global Flags:
  -i, --align string    Alignment input file (default "stdin")
      --auto-detect     Auto detects input format (overrides -p, -x and -u)
  -u, --clustal         Alignment is in clustal? default fasta
      --input-strict    Strict phylip input format (only used with -p)
  -x, --nexus           Alignment is in nexus? default fasta
      --no-block        Write Phylip sequences without space separated blocks (only used with -p)
      --one-line        Write Phylip sequences on 1 line (only used with -p)
      --output-strict   Strict phylip output format (only used with -p)
  -p, --phylip          Alignment is in phylip? default fasta
```

#### Examples

* Printing only differences seen in an input alignment:

input alignment:
```
   10   100
Seq0000  GATTAATTTG CCGTAGGCCA GAATCTGAAG ATCGAACACT TTAAGTTTTC ACTTCTAATG GAGAGGACTA GTTCATACTT TTTAAACACT TTTACATCGA
Seq0001  TGTCGGACCT AAGTATTGAG TACAACGGTG TATTCCAGCG GTGGAGAGGT CTATTTTTCC GGTTGAAGGA CTCTAGAGCT GTAAAGGGTA TGGCCATGTG
Seq0002  CTAAGCGCGG GCGGATTGCT GTTGGAGCAA GGTTAAATAC TCGGCAATGC CCCATGATCC CCCAAGGACA ATAAGAGCGA AGTTAGAACA AATGAACCCC
Seq0003  GAGTGGAGGC TTTATGGCAC AAGGTATTAG AGACTGAGGG GCACCCCGGC ATGGTAAGCA GGAGCCATCG CGAAGGCTTC AGGTATCTTC CTGTGTTACC
Seq0004  CATAGCCCCT GATGCCCTGA CCCGTGTCGC GGCAACGTCT ACATTTCACG ATAAATACTC CGCTGCTAGT CGGCTCTAGA TGCTTTTCTT CCAGATCTGG
Seq0005  AGTTTGACTA TGAGCGCCGG CTTAGTGCTG ACAGTGATGC TCCGTTGTAA GGGTCCTGAT GTTCTTGTGC TCGCGCATAT TAGAGCTGAG TTTCCCAAAG
Seq0006  TCGCCACGGT GTGGAATGTA CGTTATGGCA GTAATCAGCG GCTTTCACCG ACATGCCCCC TCCGTGGCTC CTTGCGACCA TCGGCGGACC TGCGGTGTCG
Seq0007  CTGGTAATAC CTGCGCTATT TCGTCAGTTC GTGTACGGGT AACGATAGCG GTTAATGCTT ATTCCGATCA GCTCACACCC ATGAAGGTGG CTCTGGAGCC
Seq0008  TCGTTAACCC ACTCTAACCA CCTCCTGTAG CGACATCGGG TGCTCGGCTT GGATACCTTC GTCATATTGG ACCCCAGGTC TCAACCTCGT GAGCTCTCTG
Seq0009  ACCTACGGCT CTAGACAGCT GAAGTCCGGT TCCGAGCACT GTACGGAAAC TTGAAAAGGC TCGACGGAGG CTTGTTCCGC AGAGTGGGAC TATAACATAC
```

```
gotree diff -i input.phy -p --one-line
```

Should print:

```
   10   100
Seq0000  GATTAATTTG CCGTAGGCCA GAATCTGAAG ATCGAACACT TTAAGTTTTC ACTTCTAATG
Seq0001  TG.CGGACCT AA...TTGAG T.CAAC.GT. TATTCCAG.G G.GGAGAGGT CTA.T.TTCC
Seq0002  CTAAGCGCG. G..G.TTG.T .TTGGA.C.A GGTT..ATAC .CGGCAA.G. C.CATG.TCC
Seq0003  ..G.GGAGGC TTTAT...AC A.GGTATT.. .GACTGAGGG GC.CCCCGG. .TGGTA.GCA
Seq0004  C..AGCCCCT GATGCCCTG. CCCGTGTCGC GG.A.CGT.. AC.TT.CACG .TAAA..C.C
Seq0005  AG..TGAC.A TGAGC.C.GG CTTAG..CT. .CA.TGATGC .CCGT.G.AA GGG..CTGAT
Seq0006  TCGCC.CGGT GT.G.ATGT. CGT.A..GCA G.AATCAG.G GCTTTCACCG ..A.GCCCCC
Seq0007  CTGGT.A.AC .T.CGCTATT TCG..A.TTC G.GT.CGGG. AACGA.AGCG GT.AA.GC.T
Seq0008  TCG.T.ACCC A.TCTAA... CCTC...T.. CGAC.T.GGG .GCTCGGC.T GGA.ACCT.C
Seq0009  ACC..CGGCT .TAG.CAG.T ...GTCCGGT TC...G.... G..C.GAAA. TTGAAA.GGC

   GAGAGGACTA GTTCATACTT TTTAAACACT TTTACATCGA
   .GTT.A.GG. C.CT.G.GC. G.A..GGGTA .GGC...GTG
   CCC.A.GAC. A.AAGAG.GA AG.T.GA..A AA.GA.C.CC
   .GAGCC.TCG CGAAGGCT.C AGGT.T.TTC C.GTGT.ACC
   CGCT.CTAGT CGG.TCTAGA .GCTTTTCT. CCAGATCT.G
   .TTCTTGTGC TCG.GC.TA. .AG.GCTGAG ...C.CAAAG
   TCCGT.G..C C..GCG..CA .CGGCGG..C .GCGGTGTCG
   ATTCC..TC. .C...C..CC A.G..GGTGG C.CTGGAGCC
   .TC.TATTGG ACC.CAGG.C .CA.CCTCG. GAGCTC..TG
   TC..C.GAGG C..GT.C.GC AGAGTGGGAC .A..ACATAC
```

* Counting differences seen in an input alignment:

input alignment:
```
   10   100
Seq0000  GATTAATTTG CCGTAGGCCA GAATCTGAAG ATCGAACACT TTAAGTTTTC ACTTCTAATG GAGAGGACTA GTTCATACTT TTTAAACACT TTTACATCGA
Seq0001  TGTCGGACCT AAGTATTGAG TACAACGGTG TATTCCAGCG GTGGAGAGGT CTATTTTTCC GGTTGAAGGA CTCTAGAGCT GTAAAGGGTA TGGCCATGTG
Seq0002  CTAAGCGCGG GCGGATTGCT GTTGGAGCAA GGTTAAATAC TCGGCAATGC CCCATGATCC CCCAAGGACA ATAAGAGCGA AGTTAGAACA AATGAACCCC
Seq0003  GAGTGGAGGC TTTATGGCAC AAGGTATTAG AGACTGAGGG GCACCCCGGC ATGGTAAGCA GGAGCCATCG CGAAGGCTTC AGGTATCTTC CTGTGTTACC
Seq0004  CATAGCCCCT GATGCCCTGA CCCGTGTCGC GGCAACGTCT ACATTTCACG ATAAATACTC CGCTGCTAGT CGGCTCTAGA TGCTTTTCTT CCAGATCTGG
Seq0005  AGTTTGACTA TGAGCGCCGG CTTAGTGCTG ACAGTGATGC TCCGTTGTAA GGGTCCTGAT GTTCTTGTGC TCGCGCATAT TAGAGCTGAG TTTCCCAAAG
Seq0006  TCGCCACGGT GTGGAATGTA CGTTATGGCA GTAATCAGCG GCTTTCACCG ACATGCCCCC TCCGTGGCTC CTTGCGACCA TCGGCGGACC TGCGGTGTCG
Seq0007  CTGGTAATAC CTGCGCTATT TCGTCAGTTC GTGTACGGGT AACGATAGCG GTTAATGCTT ATTCCGATCA GCTCACACCC ATGAAGGTGG CTCTGGAGCC
Seq0008  TCGTTAACCC ACTCTAACCA CCTCCTGTAG CGACATCGGG TGCTCGGCTT GGATACCTTC GTCATATTGG ACCCCAGGTC TCAACCTCGT GAGCTCTCTG
Seq0009  ACCTACGGCT CTAGACAGCT GAAGTCCGGT TCCGAGCACT GTACGGAAAC TTGAAAAGGC TCGACGGAGG CTTGTTCCGC AGAGTGGGAC TATAACATAC
```

```
gotree diff -i input.phy -p --one-line
```

Should print:
```
	AC	AG	AT	CA	CG	CT	GA	GC	GT	TA	TC	TG
Seq0001	5	12	5	5	5	6	2	2	8	7	7	10
Seq0002	5	9	7	6	3	2	3	6	3	13	7	9
Seq0003	4	10	8	5	2	8	3	7	2	6	8	16
Seq0004	8	6	11	5	4	6	1	10	4	7	11	9
Seq0005	8	12	8	5	5	4	4	2	6	7	7	7
Seq0006	10	10	5	3	7	3	3	5	6	3	12	10
Seq0007	6	9	8	2	8	4	2	6	5	9	9	5
Seq0008	11	5	9	3	4	4	4	6	4	5	10	9
Seq0009	7	9	5	4	3	4	4	5	3	9	6	11
```
