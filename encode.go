package reedsolomon

// size of sub-matrix
const UnitSize int = 16 * 1024

func (r rsBase) Encode(in, out Matrix) (err error) {
	gen := r.gen
	for i := 0; i < r.in; i++ {
		data := in[i]
		for oi := 0; oi < r.out; oi++ {
			if i == 0 {
				mulBase(gen[oi][i], data, out[oi])
			} else {
				mulXORBase(gen[oi][i], data, out[oi])
			}
		}
	}
	return
}

////////////// internal functions //////////////


func mulBase(c byte, in, out []byte) {
	mt := mulTable[c]
	for i := 0; i < len(in); i++ {
		out[i] = mt[in[i]]
	}
}

func mulXORBase(c byte, in, out []byte) {
	mt := mulTable[c]
	for i := 0; i < len(in); i++ {
		out[i] ^= mt[in[i]]
	}
}
