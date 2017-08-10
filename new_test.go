package reedsolomon

import "testing"

func BenchmarkGenTBL(b *testing.B) {
	g := genCauchyMatrix(10, 4)
	genTables(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genTables(g)
	}
}

func BenchmarkGenTBLOld(b *testing.B) {
	g := genCauchyMatrix(10, 4)
	genTablesOld(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genTablesOld(g)
	}
}
