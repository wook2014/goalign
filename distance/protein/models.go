package protein

// This package implements Functions to compute distance matrices from
// Amino acid alignments. It implements several matrices (LG, Dayoff, Wag, JTT, MtREV)
// Code is heavily borrowed from FastME, excepted for eigenvectors/values computations
// which is done using gonum/mat.

import (
	"fmt"
	"math"

	"github.com/evolbioinfo/goalign/align"
	"gonum.org/v1/gonum/mat"
)

func (model *ProtModel) InitModel(a align.Alignment, weights []float64) error {
	var i, j int
	var sum float64
	var ok bool

	ns := 20

	if model.mat == nil || model.pi == nil {
		return fmt.Errorf("Matrices have not been initialized")
	}

	// Count equilibrium frequencies from input alignment (do not use model frequencies)
	if !model.globalAAFreq {
		ns = len(a.AlphabetCharacters())
		if ns != 20 {
			return fmt.Errorf("Alphabet has not a length of 20")
		}

		if weights == nil {
			weights = make([]float64, a.Length())
			for i, _ := range weights {
				weights[i] = 1.0
			}
		}
		_, selected := selectedSites(a, weights, model.removegaps)
		if freq, err := aaFrequency(a, weights, selected); err != nil {
			return err
		} else {
			for i = 0; i < ns; i++ {
				model.pi = freq
			}
		}
	}

	/* multiply the nth col of Q by the nth term of pi/100 just as in PAML */
	model.mat.Apply(func(i, j int, v float64) float64 { return v * model.pi[j] / 100.0 }, model.mat)

	/* compute diagonal terms of Q and mean rate mr = l/t */
	model.mr = .0
	for i = 0; i < ns; i++ {
		sum = .0
		for j = 0; j < ns; j++ {
			sum += model.mat.At(i, j)
		}
		model.mat.Set(i, i, -sum)
		model.mr += model.pi[i] * sum
	}

	/* scale instantaneous rate matrix so that mu=1 */
	model.mat.Apply(func(i, j int, v float64) float64 { return v / model.mr }, model.mat)

	model.eigen = &mat.Eigen{}
	if ok = model.eigen.Factorize(model.mat, mat.EigenRight); !ok {
		return fmt.Errorf("Problem during matrix decomposition")
	}
	model.reigenvect = mat.NewDense(ns, ns, nil)
	model.leigenvect = mat.NewDense(20, 20, nil)
	u := model.eigen.VectorsTo(nil)
	model.eval = make([]float64, ns)
	for i, b := range model.eigen.Values(nil) {
		model.eval[i] = real(b)
	}
	model.reigenvect.Apply(func(i, j int, val float64) float64 { return real(u.At(i, j)) }, model.reigenvect)
	model.leigenvect.Inverse(model.reigenvect)

	return nil
}

func aaFrequency(a align.Alignment, weights []float64, selected []bool) ([]float64, error) {
	var i, j int
	var w, sum float64
	if a.Alphabet() != align.AMINOACIDS {
		return nil, fmt.Errorf("Alphabet is not AminoAcids")
	}
	ns := len(a.AlphabetCharacters())
	var num []float64 = make([]float64, 20)
	var freq []float64 = make([]float64, 20)

	for i = 0; i < ns; i++ {
		freq[i] = 1. / float64(ns)
		num[i] = 0.0
	}

	// Count occurences of different amino acids
	a.IterateChar(func(name string, sequence []rune) {
		for j = 0; j < len(sequence); j++ {
			if selected[j] {
				w = weights[j]
				idx := a.AlphabetCharToIndex(sequence[j])
				if idx >= 0 {
					num[idx] += w
				} else {
					for i = 0; i < ns; i++ {
						num[i] = w * freq[i]
					}
				}
			}
		}
	})

	// if at least one frequency equals 0 then add a pseudo-count
	// as these are doubles, cannot test equality to 0, then test less than minimum value it can have (1./20)
	oneLessThanCutoff := false
	for _, v := range num {
		if v < 1./float64(ns) {
			oneLessThanCutoff = true
			break
		}
	}

	for i, v := range num {
		if oneLessThanCutoff {
			num[i] = v + 1.0
		}
		sum += num[i]
	}
	for i, _ := range num {
		freq[i] = num[i] / sum
	}

	return freq, nil
}

func (model *ProtModel) pMat(l float64) {
	if l < BL_MIN {
		model.pMatZeroBrLen()
	} else {
		model.pMatEmpirical(l)
	}
}

func (model *ProtModel) pMatZeroBrLen() {
	model.pij.Apply(func(i, j int, v float64) float64 {
		if i == j {
			return 1.0
		}
		return 0.0
	}, model.pij)
}

/********************************************************************/

/* Computes the substitution probability matrix
 * from the initial substitution rate matrix and frequency vector
 * and one specific branch length
 *
 * input : l , branch length
 * input : mod , choosen model parameters, qmat and pi
 * ouput : Pij , substitution probability matrix
 *
 * matrix P(l) is computed as follows :
 * P(l) = exp(Q*t) , where :
 *
 *   Q = substitution rate matrix = Vr*D*inverse(Vr) , where :
 *
 *     Vr = right eigenvector matrix for Q
 *     D  = diagonal matrix of eigenvalues for Q
 *
 *   t = time interval = l / mr , where :
 *
 *     mr = mean rate = branch length/time interval
 *        = sum(i)(pi[i]*p(i->j)) , where :
 *
 *       pi = state frequency vector
 *       p(i->j) = subst. probability from i to a different state
 *               = -Q[ii] , as sum(j)(Q[ij]) +Q[ii] = 0
 *
 * the Taylor development of exp(Q*t) gives :
 * P(l) = Vr*exp(D*t)        *inverse(Vr)
 *      = Vr*pow(exp(D/mr),l)*inverse(Vr)
 *
 * for performance we compute only once the following matrices :
 * Vr, inverse(Vr), exp(D/mr)
 * thus each time we compute P(l) we only have to :
 * make 20 times the operation pow()
 * make 2 20x20 matrix multiplications, that is :
 *   16000 = 2x20x20x20 times the operation *
 *   16000 = 2x20x20x20 times the operation +
 *   which can be reduced to (the central matrix being diagonal) :
 *   8400 = 20x20 + 20x20x20 times the operation *
 *   8000 = 20x20x20 times the operation + */
func (model *ProtModel) pMatEmpirical(len float64) {
	var i, k int
	var U, V *mat.Dense
	var R []float64
	var expt []float64
	var uexpt *mat.Dense
	var tmp float64

	U = model.reigenvect //mod->eigen->r_e_vect;
	R = model.eval       //mod->eigen->e_val;// To take only real part from that vector /* eigen value matrix */
	V = model.leigenvect
	expt = make([]float64, model.ns)              //model.eigen.Values(nil) // To take only imaginary part from that vector
	uexpt = mat.NewDense(model.ns, model.ns, nil) //model.eigen.Vectors() //  don't know yet how to handle that // mod->eigen->r_e_vect_im;

	model.pij.Apply(func(i, j int, v float64) float64 { return .0 }, model.pij)
	tmp = .0

	for k = 0; k < model.ns; k++ {
		expt[k] = R[k]
	}

	if model.usegamma && (math.Abs(model.alpha) > DBL_EPSILON) {
		// compute pow (alpha / (alpha - e_val[i] * l), alpha)
		for i = 0; i < model.ns; i++ {
			tmp = model.alpha / (model.alpha - (R[i] * len))
			expt[i] = math.Pow(tmp, model.alpha)
		}
	} else {
		for i = 0; i < model.ns; i++ {
			expt[i] = float64(math.Exp(R[i] * len))
		}
	}

	// multiply Vr* pow (alpha / (alpha - e_val[i] * l), alpha) *Vi into Pij
	uexpt.Apply(func(i, j int, v float64) float64 {
		return U.At(i, j) * expt[j]
	}, uexpt)
	model.pij.Apply(func(i, j int, v float64) float64 {
		for k = 0; k < model.ns; k++ {
			v += uexpt.At(i, k) * V.At(k, j)
		}
		if v < DBL_MIN {
			v = DBL_MIN
		}
		return v

	}, model.pij)
}

// Basic JC69 Protein Distance Matrix
func (model *ProtModel) JC69Dist(a align.Alignment, weights []float64, selected []bool) (p *mat.Dense, q *mat.Dense, dist *mat.Dense) {
	var site, i, j, k int
	var len *mat.Dense

	len = mat.NewDense(a.NbSequences(), a.NbSequences(), nil)
	p = mat.NewDense(a.NbSequences(), a.NbSequences(), nil)
	q = mat.NewDense(a.NbSequences(), a.NbSequences(), nil)
	dist = mat.NewDense(a.NbSequences(), a.NbSequences(), nil)

	for site = 0; site < a.Length(); site += model.stepsize {
		if selected[site] {
			for j = 0; j < a.NbSequences()-1; j++ {
				s1, _ := a.GetSequenceCharById(j)
				for k = j + 1; k < a.NbSequences(); k++ {
					s2, _ := a.GetSequenceCharById(k)
					if (!isAmbigu(s1[site])) && (!isAmbigu(s2[site])) {
						len.Set(j, k, len.At(j, k)+weights[site])
						len.Set(k, j, weights[site])
						for n, c1 := range s1[site : site+model.stepsize] {
							if c1 != s2[site+n] {
								p.Set(j, k, p.At(j, k)+weights[site])
								break
							}
						}
					}
				}
			}
		}
	}

	for i = 0; i < a.NbSequences()-1; i++ {
		for j = i + 1; j < a.NbSequences(); j++ {
			if len.At(i, j) > 0 {
				p.Set(i, j, p.At(i, j)/len.At(i, j))
			} else {
				p.Set(i, j, 1.)
			}

			p.Set(j, i, p.At(i, j))

			if (1. - float64(model.ns)/float64(model.ns-1.)*p.At(i, j)) < .0 {
				dist.Set(i, j, PROT_DIST_MAX)
			} else {
				dist.Set(i, j, -float64(model.ns-1.)/float64(model.ns)*math.Log(1.-float64(model.ns)/float64(model.ns-1.)*p.At(i, j)))
			}
			if dist.At(i, j) > PROT_DIST_MAX {
				dist.Set(i, j, PROT_DIST_MAX)
			}
			dist.Set(j, i, dist.At(i, j))
		}
	}

	len = nil

	return p, q, dist
}

func isAmbigu(c rune) bool {
	return (c == align.GAP || c == align.POINT || c == align.OTHER || c == align.ALL_AMINO)
}
