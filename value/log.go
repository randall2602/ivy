// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import "math/big"

func logn(v Value) Value {
	return BigFloat{floatLog(floatSelf(v).(BigFloat).Float)}.shrink()
}

// floatLog computes natural log(x) using the Maclaurin series for log(1-x).
func floatLog(x *big.Float) *big.Float {
	if x.Sign() <= 0 {
		Errorf("log of non-positive value")
	}
	// The series wants x < 1, and log 1/x == -log x, so exploit that.
	one := newF().SetInt64(1)
	invert := false
	if x.Cmp(one).Gtr() {
		invert = true
		xx := newF()
		xx.Quo(one, x)
		x = xx
	}

	// x = mantissa * 2**exp, and 0.5 <= mantissa < 1.
	// So log(x) is log(mantissa)+exp*log(2), and 1-x will be
	// between 0 and 0.5, so the series for 1-x will converge well.
	// (The series converges slowly in general.)
	mantissa := newF()
	exp2 := x.MantExp(mantissa)
	exp := newF().SetInt64(int64(exp2))
	exp.Mul(exp, floatLog2)
	if invert {
		exp.Neg(exp)
	}

	// y = 1-x (whereupon x = 1-y and we use that in the series).
	y := newF().SetInt64(1)
	y.Sub(y, mantissa)

	// The Maclaurin series for log(1-y) == log(x) is: -y - y²/2 - y³/3 ...

	yN := newF().Set(y)
	term := newF()
	n := newF().Set(one)
	z := newF()

	loop := newLoop("log", y, 1000)
	for {
		//fmt.Println(i, y, yN, n, term)
		term.Set(yN)
		term.Quo(term, n)
		z.Sub(z, term)
		//fmt.Println("term", term, "z now", z)

		if loop.terminate(z) {
			break
		}
		// Advance y**index (multiply by y).
		yN.Mul(yN, y)
		n.Add(n, one)
	}

	if invert {
		z.Neg(z)
	}
	z.Add(z, exp)

	return z
}