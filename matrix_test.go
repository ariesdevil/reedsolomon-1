package reedsolomon

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

//func TestGenEncodeCauchy(t *testing.T) {
//	c1 := GenEncodeMatrix(4, 2)
//	c2 := genEncMatrixCauchy(4, 2)
//	fmt.Println(c1)
//	fmt.Println(c2)
//}
//
//func TestNewInvert(t *testing.T) {
//	raw1 := GenEncodeMatrix(4, 2)
//	raw2 := genEncMatrixCauchy(4, 2)
//	raw1.swapRows(0, 4)
//	raw1.swapRows(1, 5)
//	raw2.swap(0, 4, 4)
//	raw2.swap(1, 5, 4)
//	fmt.Println(raw1[:4].invert())
//	fmt.Println(raw2[:16].invert(4))
//}

// TODO cmp newInvert invert

func TestNewInvert(t *testing.T) {
	for i := 1; i < 128; i++ {
		m := genCauchyMatrix(i, i)
		i1, err := m.invert()
		if err != nil {
			t.Fatal(err)
		}
		ms := newMatrix(i, i)
		for j := range m {
			copy(ms[j*i:j*i+i], m[j])
		}
		i2, err := ms.invert(i)
		i1s := make([]byte, i*i)
		for j := range i1 {
			copy(i1s[j*i:j*i+i], i1[j])
		}
		if !bytes.Equal(i1s, i2) {
			t.Fatal("invert new fault", i1, i2)
		}
	}
}

func TestMatrixInverse(t *testing.T) {
	testCases := []struct {
		matrixData     [][]byte
		expectedResult string
		shouldPass     bool
		expectedErr    error
	}{
		// Test case validating inverse of the input Matrix.
		{
			// input
			[][]byte{
				[]byte{56, 23, 98},
				[]byte{3, 100, 200},
				[]byte{45, 201, 123},
			},
			// expected
			"[[175, 133, 33], [130, 13, 245], [112, 35, 126]]",
			// expected to pass.
			true,
			nil,
		},
		// Test case Matrix[0][0] == 0
		{
			[][]byte{
				[]byte{0, 23, 98},
				[]byte{3, 100, 200},
				[]byte{45, 201, 123},
			},
			"[[245, 128, 152], [188, 64, 135], [231, 81, 239]]",
			true,
			nil,
		},
		// Test case validating inverse of the input Matrix.
		{
			// input
			[][]byte{
				[]byte{1, 0, 0, 0, 0},
				[]byte{0, 1, 0, 0, 0},
				[]byte{0, 0, 0, 1, 0},
				[]byte{0, 0, 0, 0, 1},
				[]byte{7, 7, 6, 6, 1},
			},
			// expected
			"[[1, 0, 0, 0, 0]," +
				" [0, 1, 0, 0, 0]," +
				" [123, 123, 1, 122, 122]," +
				" [0, 0, 1, 0, 0]," +
				" [0, 0, 0, 1, 0]]",
			true,
			nil,
		},
		// Test case with singular Matrix.
		// expected to fail with error errSingular.
		{

			[][]byte{
				[]byte{4, 2},
				[]byte{12, 6},
			},
			"",
			false,
			ErrSingular,
		},
	}

	for i, testCase := range testCases {
		m := newMatrixData(testCase.matrixData)
		actualResult, actualErr := m.invert()
		if actualErr != nil && testCase.shouldPass {
			t.Errorf("Test %r: Expected to pass, but failed with: <ERROR> %s", i+1, actualErr.Error())
		}
		if actualErr == nil && !testCase.shouldPass {
			t.Errorf("Test %r: Expected to fail with <ERROR> \"%s\", but passed instead.", i+1, testCase.expectedErr)
		}
		// Failed as expected, but does it fail for the expected reason.
		if actualErr != nil && !testCase.shouldPass {
			if testCase.expectedErr != actualErr {
				t.Errorf("Test %r: Expected to fail with error \"%s\", but instead failed with error \"%s\" instead.", i+1, testCase.expectedErr, actualErr)
			}
		}
		// Test passes as expected, but the output values
		// are verified for correctness here.
		if actualErr == nil && testCase.shouldPass {
			if testCase.expectedResult != actualResult.string() {
				t.Errorf("Test %r: The inverse Matrix doesnt't match the expected result", i+1)
			}
		}
	}
}

func BenchmarkInvert10x10(b *testing.B) {
	benchmarkInvert(b, 10)
}

func BenchmarkInvert28x28(b *testing.B) {
	benchmarkInvert(b, 28)
}

func BenchmarkNewInvert10x10(b *testing.B) {
	benchmarkInvertNew(b, 10)
}

func BenchmarkNewInvert28x28(b *testing.B) {
	benchmarkInvertNew(b, 28)
}

func benchmarkInvert(b *testing.B, size int) {
	m := GenEncodeMatrix(size, 2)
	m.swapRows(0, size)
	m.swapRows(1, size+1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m[:size].invert()
		if err != nil {
			b.Fatal(b)
		}
	}
}

func benchmarkInvertNew(b *testing.B, size int) {
	m := genEncMatrixCauchy(size, 2)
	m.swap(0, size, size)
	m.swap(1, size+1, size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//copy(m1, m[:size*size])
		_, err := matrix(m[:size*size]).invert(size)
		if err != nil {
			b.Fatal(b)
		}
	}
}

// new a Matrix with Data
func newMatrixData(data [][]byte) Matrix {
	m := Matrix(data)
	return m
}

func (m Matrix) string() string {
	rowOut := make([]string, 0, len(m))
	for _, row := range m {
		colOut := make([]string, 0, len(row))
		for _, col := range row {
			colOut = append(colOut, strconv.Itoa(int(col)))
		}
		rowOut = append(rowOut, "["+strings.Join(colOut, ", ")+"]")
	}
	return "[" + strings.Join(rowOut, ", ") + "]"
}
