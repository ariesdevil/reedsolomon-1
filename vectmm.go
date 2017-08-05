package reedsolomon

func (r rsAVX2) Encode(in, out Matrix) (err error) {
	off := 0
	for oi := 0; oi < r.out; oi++ {
		vectMul(r.tables[off:], in, out[oi])
		off += r.in * 32
	}
	return
}

//go:noescape
func vectMul(tbl []byte, in Matrix, out []byte)

//// Size of Shard must be integral multiple of 256B
//func (r rsAVX2) Encode(in, out Matrix) (err error) {
//	size := len(in[0])
//	start, end := 0, 0
//	do := UnitSize
//	if size <= UnitSize {
//		r.matrixMul(start, size, in, out)
//	} else {
//		for start < size {
//			end = start + do
//			if end <= size {
//				r.matrixMul(start, end, in, out)
//				start = end
//			} else {
//				r.matrixMul(start, size, in, out)
//				start = size
//			}
//		}
//	}
//	return
//}
//
//func (r rsAVX2) matrixMul(start, end int, in, out Matrix) {
//	for i := 0; i < r.in; i++ {
//		tmp := i * r.out
//		for oi := 0; oi < r.out; oi++ {
//			offset := (tmp + oi) * 32
//			table := r.tables[offset : offset+32]
//			if i == 0 {
//				mulAVX2(table, in[i][start:end], out[oi][start:end])
//			} else {
//				mulXORAVX2(table, in[i][start:end], out[oi][start:end])
//			}
//		}
//	}
//}
//
////go:noescape
//func mulAVX2(table, in, out []byte)
//
////go:noescape
//func mulXORAVX2(table, in, out []byte)
