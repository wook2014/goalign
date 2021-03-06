package protein

import (
	"fmt"
	"math"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/goalign/io"
	"gonum.org/v1/gonum/mat"
)

const (
	BRENT_ITMAX = 10000
	BRENT_ZEPS  = 1.e-10
	BRENT_CGOLD = 0.3819660
	DBL_MIN     = 2.2250738585072014e-308

	BL_MIN = 1.e-08
	BL_MAX = 100.0
)

var DBL_EPSILON float64 = math.Nextafter(1, 2) - 1

func (model *ProtDistModel) MLDist(a align.Alignment, weights []float64) (p, q, dist *mat.Dense, err error) {
	var j, k, l int
	var state0, state1 int
	var len float64
	var init, sum float64
	var warn bool
	var d_max float64
	var Fs *mat.Dense
	var selected []bool

	if a.Alphabet() != align.AMINOACIDS {
		err = fmt.Errorf("Cannot compute protein distance with this alignment: Wrong alphabet")
	}

	if weights == nil {
		weights = make([]float64, a.Length())
		for i := range weights {
			weights[i] = 1.
		}
	}

	_, selected = selectedSites(a, weights, model.removegaps)

	p, q, dist = model.JC69Dist(a, weights, selected)

	warn = false

	// Create F for one thread
	Fs = mat.NewDense(model.Ns(), model.Ns(), nil)

	// Alignment of 2 sequences
	for j = 0; j < a.NbSequences(); j++ { // begin for j->n_otu
		seq1, _ := a.GetSequenceCharById(j)
		for k = j + 1; k < a.NbSequences(); k++ { // begin for k->n_otu
			seq2, _ := a.GetSequenceCharById(k)
			pair := seqpairdist{j, k, seq1, seq2, nil, nil}

			checkAmbiguities(&pair, 1)
			// If sequences are different (avoid ambiguities), compute distance
			// Else distance = 0
			if check2SequencesDiff(&pair) {
				//Hide_Ambiguities(pair)
				init = dist.At(j, k)

				if (init == PROT_DIST_MAX) || (init < .0) {
					init = 0.1
				}
				d_max = init
				Fs.Apply(func(i, j int, v float64) float64 { return .0 }, Fs)
				len = 0.0

				for l = 0; l < a.Length(); l++ {
					if selected[l] {
						w := weights[l]
						if pair.seq1Ambigu[l] || pair.seq2Ambigu[l] {
							w = 0.0
						}
						state0 = a.AlphabetCharToIndex(pair.seq1[l])
						state1 = a.AlphabetCharToIndex(pair.seq2[l])
						if (state0 > -1) && (state1 > -1) {
							Fs.Set(state0, state1, Fs.At(state0, state1)+w)
							len += w
						}
					}
				}

				if len > .0 {
					Fs.Apply(func(i, j int, v float64) float64 { return (v / len) }, Fs)
				}

				sum = mat.Sum(Fs)

				if sum < .001 {
					d_max = -1.
				} else if (sum > 1.-.001) && (sum < 1.+.001) {
					d_max = model.opt_Dist_F(d_max, Fs)
				} else {
					return nil, nil, nil, fmt.Errorf("Invalid value when computing distance. sum = %f.", sum)
				}

				if d_max >= PROT_DIST_MAX {
					warn = true
					d_max = PROT_DIST_MAX
				}
			} else {
				// Do not correct for dist < BL_MIN,
				// otherwise Fill_Missing_Dist will not be called
				d_max = 0.
			}
			dist.Set(j, k, d_max)
			dist.Set(k, j, d_max)
		} // end for k->n_otu
	} // end for j->n_otu
	if warn {
		io.PrintMessage(fmt.Sprintf("Give up this dataset because at least one distance exceeds %.2f.", PROT_DIST_MAX))
	}

	return
}

func (model *ProtDistModel) lk_Dist(F *mat.Dense, dist float64) float64 {
	var i, j int
	var len, lnL float64

	len = -1.

	len = dist // * model.gamma_rr

	if len < BL_MIN {
		len = BL_MIN
	} else if len > BL_MAX {
		len = BL_MAX
	}
	model.pMat(len)

	lnL = .0

	ns := model.Ns()
	for i = 0; i < ns; i++ {
		for j = 0; j < ns; j++ {
			lnL += F.At(i, j) * math.Log(model.partialLK(i, j))
		}
	}

	return lnL
}

func (model *ProtDistModel) partialLK(i, j int) float64 {
	var lk float64
	lk = .0

	lk += model.model.Pi(i) * model.pij.At(i, j) // * model.gamma_r_proba
	return lk
}

func (model *ProtDistModel) opt_Dist_F(dist float64, F *mat.Dense) float64 {
	var ax, bx, cx float64
	var optdist float64

	if dist < BL_MIN {
		dist = BL_MIN
	}

	ax = BL_MIN
	bx = dist
	cx = BL_MAX

	optdist = dist
	model.dist_F_Brent(ax, bx, cx, 1.E-10, 1000, &optdist, F)
	return optdist
}

func (model *ProtDistModel) dist_F_Brent(ax, bx, cx, tol float64, n_iter_max int, param *float64, F *mat.Dense) float64 {
	var iter int
	var a, b, d, etemp, fu, fv, fw, fx, p, q, r, tol1, tol2, u, v, w, x, xm float64
	var curr_lnL float64
	var e float64
	e = 0.0

	//optimize distance, not likelihood
	var old_param, cur_param float64

	d = 0.0
	if ax < cx {
		a = ax
		b = cx
	} else {
		a = cx
		b = ax
	}

	x = bx
	w = bx
	v = bx
	fw = -model.lk_Dist(F, math.Abs(bx))
	fv = fw
	fx = fw
	curr_lnL = -fw

	old_param = math.Abs(bx)
	cur_param = math.Abs(bx)

	for iter = 1; iter <= BRENT_ITMAX; iter++ {
		xm = 0.5 * (a + b)

		tol1 = tol*math.Abs(x) + BRENT_ZEPS
		tol2 = 2.0 * tol1

		if (iter > 1) && math.Abs(old_param-cur_param) < 1.E-06 {
			*param = x
			curr_lnL = model.lk_Dist(F, *param)
			return -curr_lnL
		}

		if math.Abs(e) > tol1 {
			r = (x - w) * (fx - fv)
			q = (x - v) * (fx - fw)
			p = (x-v)*q - (x-w)*r
			q = 2.0 * (q - r)
			if q > 0.0 {
				p = -p
			}

			q = math.Abs(q)
			etemp = e
			e = d

			if math.Abs(p) >= math.Abs(0.5*q*etemp) || p <= q*(a-x) || p >= q*(b-x) {
				if x >= xm {
					e = a - x
				} else {
					e = b - x
				}
				d = BRENT_CGOLD * e
			} else {
				d = p / q
				u = x + d
				if u-a < tol2 || b-u < tol2 {
					d = sign(tol1, xm-x)
				}
			}
		} else {
			if x >= xm {
				e = a - x
			} else {
				e = b - x
			}
			d = BRENT_CGOLD * e
		}

		if math.Abs(d) >= tol1 {
			u = x + d
		} else {
			u = x + sign(tol1, d)
		}
		if u < BL_MIN {
			u = BL_MIN
		}

		(*param) = math.Abs(u)
		fu = -model.lk_Dist(F, math.Abs(u))
		curr_lnL = -fu

		if fu <= fx {
			if iter > n_iter_max {
				return -fu
			}
			if u >= x {
				a = x
			} else {
				b = x
			}

			shift(&v, &w, &x, &u)
			shift(&fv, &fw, &fx, &fu)
		} else {
			if u < x {
				a = u
			} else {
				b = u
			}

			if fu <= fw || (math.Abs(w-x) < DBL_EPSILON) {
				v = w
				w = u
				fv = fw
				fw = fu
			} else if fu <= fv || (math.Abs(v-x) < DBL_EPSILON) || (math.Abs(v-w) < DBL_EPSILON) {
				v = u
				fv = fu
			}
		}
		old_param = cur_param
		cur_param = *param
	}

	panic("Too many iterations in BRENT.")

	return (-1)
}
