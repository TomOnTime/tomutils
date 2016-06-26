package main

// To run these tests:
//     go test -bench=.

import "testing"

var result int

func benchmarkLs(i int, b *testing.B) {
	var r int
	mylist := MakeList(i)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = LinearSearch(mylist)
	}
	result = r
}

func benchmarkSp(i int, b *testing.B) {
	var r int
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		mylist := MakeList(i)
		b.StartTimer()
		r = SortAndPick(mylist)
	}
	result = r
}

func BenchmarkSp10(b *testing.B)        { benchmarkSp(10, b) }
func BenchmarkSp20(b *testing.B)        { benchmarkSp(20, b) }
func BenchmarkSp30(b *testing.B)        { benchmarkSp(30, b) }
func BenchmarkSp40(b *testing.B)        { benchmarkSp(40, b) }
func BenchmarkSp50(b *testing.B)        { benchmarkSp(50, b) }
func BenchmarkSp60(b *testing.B)        { benchmarkSp(60, b) }
func BenchmarkSp70(b *testing.B)        { benchmarkSp(70, b) }
func BenchmarkSp80(b *testing.B)        { benchmarkSp(90, b) }
func BenchmarkSp90(b *testing.B)        { benchmarkSp(90, b) }
func BenchmarkSp100(b *testing.B)       { benchmarkSp(100, b) }
func BenchmarkSp1000(b *testing.B)      { benchmarkSp(1000, b) }
func BenchmarkSp10000(b *testing.B)     { benchmarkSp(10000, b) }
func BenchmarkSp100000(b *testing.B)    { benchmarkSp(100000, b) }
func BenchmarkSp200000(b *testing.B)    { benchmarkSp(200000, b) }
func BenchmarkSp300000(b *testing.B)    { benchmarkSp(300000, b) }
func BenchmarkSp400000(b *testing.B)    { benchmarkSp(400000, b) }
func BenchmarkSp500000(b *testing.B)    { benchmarkSp(500000, b) }
func BenchmarkSp600000(b *testing.B)    { benchmarkSp(600000, b) }
func BenchmarkSp700000(b *testing.B)    { benchmarkSp(700000, b) }
func BenchmarkSp800000(b *testing.B)    { benchmarkSp(800000, b) }
func BenchmarkSp900000(b *testing.B)    { benchmarkSp(900000, b) }
func BenchmarkSp1000000(b *testing.B)   { benchmarkSp(1000000, b) }
func BenchmarkSp10000000(b *testing.B)  { benchmarkSp(10000000, b) }
func BenchmarkSp20000000(b *testing.B)  { benchmarkSp(20000000, b) }
func BenchmarkSp30000000(b *testing.B)  { benchmarkSp(30000000, b) }
func BenchmarkSp40000000(b *testing.B)  { benchmarkSp(40000000, b) }
func BenchmarkSp50000000(b *testing.B)  { benchmarkSp(50000000, b) }
func BenchmarkSp60000000(b *testing.B)  { benchmarkSp(60000000, b) }
func BenchmarkSp70000000(b *testing.B)  { benchmarkSp(70000000, b) }
func BenchmarkSp80000000(b *testing.B)  { benchmarkSp(80000000, b) }
func BenchmarkSp90000000(b *testing.B)  { benchmarkSp(90000000, b) }
func BenchmarkSp100000000(b *testing.B) { benchmarkSp(100000000, b) }
func BenchmarkSp200000000(b *testing.B) { benchmarkSp(200000000, b) }
func BenchmarkSp300000000(b *testing.B) { benchmarkSp(300000000, b) }

//func BenchmarkSp400000000(b *testing.B)  { benchmarkSp(400000000, b) }
//func BenchmarkSp500000000(b *testing.B)  { benchmarkSp(500000000, b) }
//func BenchmarkSp600000000(b *testing.B)  { benchmarkSp(600000000, b) }
//func BenchmarkSp700000000(b *testing.B)  { benchmarkSp(700000000, b) }
//func BenchmarkSp800000000(b *testing.B)  { benchmarkSp(800000000, b) }
//func BenchmarkSp900000000(b *testing.B)  { benchmarkSp(900000000, b) }
//func BenchmarkSp1000000000(b *testing.B) { benchmarkSp(1000000000, b) }

func BenchmarkLs10(b *testing.B)        { benchmarkLs(10, b) }
func BenchmarkLs20(b *testing.B)        { benchmarkLs(20, b) }
func BenchmarkLs30(b *testing.B)        { benchmarkLs(30, b) }
func BenchmarkLs40(b *testing.B)        { benchmarkLs(40, b) }
func BenchmarkLs50(b *testing.B)        { benchmarkLs(50, b) }
func BenchmarkLs60(b *testing.B)        { benchmarkLs(60, b) }
func BenchmarkLs70(b *testing.B)        { benchmarkLs(70, b) }
func BenchmarkLs80(b *testing.B)        { benchmarkLs(90, b) }
func BenchmarkLs90(b *testing.B)        { benchmarkLs(90, b) }
func BenchmarkLs100(b *testing.B)       { benchmarkLs(100, b) }
func BenchmarkLs1000(b *testing.B)      { benchmarkLs(1000, b) }
func BenchmarkLs10000(b *testing.B)     { benchmarkLs(10000, b) }
func BenchmarkLs100000(b *testing.B)    { benchmarkLs(100000, b) }
func BenchmarkLs1000000(b *testing.B)   { benchmarkLs(1000000, b) }
func BenchmarkLs10000000(b *testing.B)  { benchmarkLs(10000000, b) }
func BenchmarkLs20000000(b *testing.B)  { benchmarkLs(20000000, b) }
func BenchmarkLs30000000(b *testing.B)  { benchmarkLs(30000000, b) }
func BenchmarkLs40000000(b *testing.B)  { benchmarkLs(40000000, b) }
func BenchmarkLs50000000(b *testing.B)  { benchmarkLs(50000000, b) }
func BenchmarkLs60000000(b *testing.B)  { benchmarkLs(60000000, b) }
func BenchmarkLs70000000(b *testing.B)  { benchmarkLs(70000000, b) }
func BenchmarkLs80000000(b *testing.B)  { benchmarkLs(80000000, b) }
func BenchmarkLs90000000(b *testing.B)  { benchmarkLs(90000000, b) }
func BenchmarkLs100000000(b *testing.B) { benchmarkLs(100000000, b) }
func BenchmarkLs200000000(b *testing.B) { benchmarkLs(200000000, b) }

//func BenchmarkLs300000000(b *testing.B) { benchmarkLs(300000000, b) }
//func BenchmarkLs400000000(b *testing.B)  { benchmarkLs(400000000, b) }
//func BenchmarkLs500000000(b *testing.B)  { benchmarkLs(500000000, b) }
//func BenchmarkLs600000000(b *testing.B)  { benchmarkLs(600000000, b) }
//func BenchmarkLs700000000(b *testing.B)  { benchmarkLs(700000000, b) }
//func BenchmarkLs800000000(b *testing.B)  { benchmarkLs(800000000, b) }
//func BenchmarkLs900000000(b *testing.B)  { benchmarkLs(900000000, b) }
//func BenchmarkLs1000000000(b *testing.B) { benchmarkLs(1000000000, b) }
