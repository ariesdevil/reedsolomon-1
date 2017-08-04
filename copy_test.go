package reedsolomon

import (
	"bytes"
	"fmt"
	"testing"
)

// TODO: gen cauchy matrix
func TestVerifyTblgfw(t *testing.T) {
	g := genCauchyMatrix(10, 4)
	t1 := genTables(g)
	t2 := genTablesgfw(g)
	for i := 0; i < 10; i++ {
		for oi := 0; oi < 4; oi++ {
			offset1 := (i*4 + oi) * 32
			tmp1 := t1[offset1 : offset1+32]
			offset2 := (oi*10 + i) * 32
			tmp2 := t2[offset2 : offset2+32]
			if !bytes.Equal(tmp1, tmp2) {
				t.Fatal("error")
			}

		}
	}
}
func TestVerifyComb(t *testing.T) {
	g := genCauchyMatrix(10, 4)
	t1 := genTablesComb(g)
	t2 := genTables(g)
	fmt.Println(len(t2), 40*32)
	if !bytes.Equal(t1, t2) {
		t.Fatal("shit")
	}
}

func BenchmarkCombGenTables(b *testing.B) {
	g := genCauchyMatrix(10, 4)
	genTablesComb(g)
	b.SetBytes(10 * 4 * 32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genTablesComb(g)
	}
}

func genTablesComb(gen Matrix) []byte {
	rows := len(gen)
	cols := len(gen[0])
	tables := make([]byte, 32*rows*cols)
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			c := gen[j][i]
			ct := combTable[c][:]
			offset := (i*rows + j) * 32
			copy32B(tables[offset:offset+32], ct)
		}
	}
	return tables
}

func TestVeryNewCopy(t *testing.T) {
	g := genCauchyMatrix(10, 4)
	t1 := genTables(g)
	t2 := genTablesNew(g)
	if !bytes.Equal(t1, t2) {
		t.Fatal("tables new failed")
	}

}

func BenchmarkCopy(b *testing.B) {
	dst := make([]byte, 16)
	src := make([]byte, 16)
	copy(dst, src)
	b.SetBytes(int64(16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(dst, src)
	}
}

func BenchmarkCopySSE2(b *testing.B) {
	dst := make([]byte, 16)
	src := make([]byte, 16)
	copySSE2(dst, src)
	b.SetBytes(int64(16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copySSE2(dst, src)
	}
}

func BenchmarkGenTables(b *testing.B) {
	g := genCauchyMatrix(10, 4)
	genTables(g)
	b.SetBytes(10 * 4 * 32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genTables(g)
	}
}

func BenchmarkGenTablesNew(b *testing.B) {
	g := genCauchyMatrix(10, 4)
	genTablesNew(g)
	b.SetBytes(10 * 4 * 32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genTablesNew(g)
	}
}
