package reedsolomon

import (
	"errors"
)

type Matrix [][]byte // byte[row][col]

type matrix []byte

func newMatrix(rows, cols int) matrix {
	m := make([]byte, rows*cols)
	return m
}
func NewMatrix(rows, cols int) Matrix {
	m := Matrix(make([][]byte, rows))
	for i := range m {
		m[i] = make([]byte, cols)
	}
	return m
}

// return identity Matrix(upper) cauchy Matrix(lower)
func GenEncodeMatrix(d, p int) Matrix {
	rows := d + p
	cols := d
	m := NewMatrix(rows, cols)
	// identity Matrix
	for j := 0; j < cols; j++ {
		m[j][j] = byte(1)
	}
	// cauchy Matrix
	c := genCauchyMatrix(d, p)
	for i, v := range c {
		copy(m[d+i], v)
	}
	return m
}

func genEncMatrixCauchy(data, parity int) matrix {
	rows := data + parity
	cols := data
	m := newMatrix(rows, cols)
	// identity matrix
	for j := 0; j < cols; j++ {
		m[j*data+j] = byte(1)
	}
	// cauchy matrix
	p := data * data
	for i := cols; i < rows; i++ {
		for j := 0; j < cols; j++ {
			d := i ^ j
			a := inverseTable[d]
			m[p] = byte(a)
			p++
		}
	}
	return m
}

func genCauchyMatrix(d, p int) Matrix {
	rows := d + p
	cols := d
	m := NewMatrix(p, cols)
	start := 0
	for i := cols; i < rows; i++ {
		for j := 0; j < cols; j++ {
			d := i ^ j
			a := inverseTable[d]
			m[start][j] = byte(a)
		}
		start++
	}
	return m
}

// TODO need test
func (m matrix) swap(i, j, n int) {
	for k := 0; k < n; k++ {
		m[i*n+k], m[j*n+k] = m[j*n+k], m[i*n+k]
	}
}

//// TODO need copy m for invert
//func (m matrix) invert(n int) (out matrix, err error) {
//
//	// Set out_mat[] to the identity matrix
//	out = newMatrix(n, n)
//	for i := 0; i < n; i++ {
//		out[i*n+i] = byte(1)
//	}
//	// Inverse
//	for i := 0; i < n; i++ {
//		if m[i*n+i] == 0 {
//			for j := i + 1; j < n; j++ {
//				if m[j*n+i] == 0 {
//					break
//				}
//
//				// TODO how could j == n?
//				if j == n {
//					err = ErrSingular
//					return
//				}
//				m.swap(i, j, n)
//				out.swap(i, j, n)
//			}
//		}
//
//		tmp := inverseTable[m[i*n+i]]
//		for j := 0; j < n; j++ {
//			m[i*n+j] = gfMul(m[i*n+j], tmp)
//			out[i*n+j] = gfMul(out[i*n+j], tmp)
//		}
//
//		for j := 0; j < n; j++ {
//			if j == i {
//				continue
//			}
//			tmp = m[j*n+i]
//			for k := 0; k < n; k++ {
//				out[j*n+k] ^= gfMul(tmp, out[i*n+k])
//				m[j*n+k] ^= gfMul(tmp, m[i*n+k])
//			}
//		}
//	}
//	return
//}

func (m Matrix) invert() (Matrix, error) {
	size := len(m)
	iM := identityMatrix(size)
	mIM, _ := m.augIM(iM)

	err := mIM.gaussJordan()
	if err != nil {
		return nil, err
	}
	return mIM.subMatrix(size), nil
}

func (m matrix) invert(n int) (matrix, error) {
	raw := newMatrix(n, 2*n)
	for i := 0; i < n; i++ {
		t := i * n
		copy(raw[2*t:2*t+n], m[t:t+n])
		raw[2*t+i+n] = byte(1)
	}
	err := raw.gaussJordan(n, 2*n)
	if err != nil {
		return nil, err
	}
	return raw.subMatrix(n), nil
}

// IN -> (IN|I)
func (m Matrix) augIM(iM Matrix) (Matrix, error) {
	result := NewMatrix(len(m), len(m[0])+len(iM[0]))
	for r, row := range m {
		for c := range row {
			result[r][c] = m[r][c]
		}
		cols := len(m[0])
		for c := range iM[0] {
			result[r][cols+c] = iM[r][c]
		}
	}
	return result, nil
}

var ErrSingular = errors.New("reedsolomon: Matrix is singular")

func (m matrix) gaussJordan(rows, columns int) error {
	for r := 0; r < rows; r++ {
		// If the element on the diagonal is 0, find a row below
		// that has a non-zero and swap them.
		if m[2*r*rows+r] == 0 {
			for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
				if m[2*rowBelow*rows+r] != 0 {
					m.swap(r, rowBelow, 2*rows)
					break
				}
			}
		}
		// After swap, if we find all elements in this column is 0, it means the Matrix's det is 0
		if m[2*r*rows+r] == 0 {
			return ErrSingular
		}
		// Scale to 1.
		if m[2*r*rows+r] != 1 {
			d := m[2*r*rows+r]
			scale := inverseTable[d]
			// every element(this column) * m[r][r]'s inverse
			for c := 0; c < columns; c++ {
				m[2*r*rows+c] = gfMul(m[2*r*rows+c], scale)
			}
		}
		//Make everything below the 1 be a 0 by subtracting a multiple of it
		for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
			if m[2*rowBelow*rows+r] != 0 {
				// scale * m[r][r] = scale, scale + scale = 0
				// makes m[r][r+1] = 0 , then calc left elements
				scale := m[2*rowBelow*rows+r]
				for c := 0; c < columns; c++ {
					m[2*rowBelow*rows+c] ^= gfMul(scale, m[2*r*rows+c])
				}
			}
		}
	}
	// Now clear the part above the main diagonal.
	// same logic with clean upper
	for d := 0; d < rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m[2*rowAbove*rows+d] != 0 {
				scale := m[2*rowAbove*rows+d]
				for c := 0; c < columns; c++ {
					m[2*rowAbove*rows+c] ^= gfMul(scale, m[2*d*rows+c])
				}
			}
		}
	}
	return nil
}

// (IN|I) -> (I|OUT)
func (m Matrix) gaussJordan() error {
	rows := len(m)
	columns := len(m[0])
	// Clear out the part below the main diagonal and scale the main
	// diagonal to be 1.
	for r := 0; r < rows; r++ {
		// If the element on the diagonal is 0, find a row below
		// that has a non-zero and swap them.
		if m[r][r] == 0 {
			for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
				if m[rowBelow][r] != 0 {
					m.swapRows(r, rowBelow)
					break
				}
			}
		}
		// After swap, if we find all elements in this column is 0, it means the Matrix's det is 0
		if m[r][r] == 0 {
			return ErrSingular
		}
		// Scale to 1.
		if m[r][r] != 1 {
			d := m[r][r]
			scale := inverseTable[d]
			// every element(this column) * m[r][r]'s inverse
			for c := 0; c < columns; c++ {
				m[r][c] = gfMul(m[r][c], scale)
			}
		}
		//Make everything below the 1 be a 0 by subtracting a multiple of it
		for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
			if m[rowBelow][r] != 0 {
				// scale * m[r][r] = scale, scale + scale = 0
				// makes m[r][r+1] = 0 , then calc left elements
				scale := m[rowBelow][r]
				for c := 0; c < columns; c++ {
					m[rowBelow][c] ^= gfMul(scale, m[r][c])
				}
			}
		}
	}
	// Now clear the part above the main diagonal.
	// same logic with clean upper
	for d := 0; d < rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m[rowAbove][d] != 0 {
				scale := m[rowAbove][d]
				for c := 0; c < columns; c++ {
					m[rowAbove][c] ^= gfMul(scale, m[d][c])
				}
			}
		}
	}
	return nil
}

func identityMatrix(n int) Matrix {
	m := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		m[i][i] = byte(1)
	}
	return m
}

// (I|OUT) -> OUT
func (m Matrix) subMatrix(size int) Matrix {
	result := NewMatrix(size, size)
	for r := 0; r < size; r++ {
		for c := size; c < size*2; c++ {
			result[r][c-size] = m[r][c]
		}
	}
	return result
}

func (m matrix) subMatrix(size int) matrix {
	ret := newMatrix(size, size)
	for i := 0; i < size; i++ {
		copy(ret[i*size:i*size+size], m[2*i*size+size:2*i*size+2*size])
	}
	return ret
}

// SwapRows Exchanges two rows in the Matrix.
func (m Matrix) swapRows(r1, r2 int) {
	m[r2], m[r1] = m[r1], m[r2]
}

func gfMul(a, b byte) byte {
	return mulTable[a][b]
}
