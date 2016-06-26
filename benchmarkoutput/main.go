package main

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

var result int

func MakeList(n int) []int {
	hirange := 1000000
	r := make([]int, n)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i, _ := range r {
		r[i] = r1.Intn(hirange)
	}
	return r
}

func LinearSearch(l []int) int {
	max := l[0]
	for n, _ := range l {
		if n > max {
			max = n
		}
	}
	return max
}

func SortAndPick(l []int) int {
	sort.Ints(l)
	return l[len(l)-1]
}

//func main() {
//	x := MakeList(10000)
//	fmt.Println(x)
//	fmt.Println(SortAndPick(x))
//}

func benchmarkLs(i int) {
	fn := func(b *testing.B) {
		var r int
		mylist := MakeList(i)
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			r = LinearSearch(mylist)
		}
		result = r
	}
	r := testing.Benchmark(fn)

	fmt.Printf("i=%d\t", i)
	fmt.Printf("%d ns/op\t", int(r.T)/r.N)
	fmt.Printf("%d ns/op/i\n", int(r.T)/r.N/i)
}

func benchmarkSp(i int) {
	fn := func(b *testing.B) {
		var r int
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			mylist := MakeList(i)
			b.StartTimer()
			r = SortAndPick(mylist)
		}
		result = r
	}
	r := testing.Benchmark(fn)

	fmt.Printf("i=%d\t", i)
	fmt.Printf("%d ns/op\t", int(r.T)/r.N)
	fmt.Printf("%d ns/op/i\n", int(r.T)/r.N/i)
}

func main() {
	benchmarkLs(10)
	benchmarkLs(100)
	benchmarkLs(1000)
	benchmarkLs(10000)
	benchmarkLs(100000)
	benchmarkLs(1000000)
	benchmarkSp(10)
	benchmarkSp(100)
	benchmarkSp(1000)
	benchmarkSp(10000)
	benchmarkSp(100000)
	benchmarkSp(1000000)
}
