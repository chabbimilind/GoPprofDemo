package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func J_expect_18_18(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func I_expect_16_36(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func H_expect_14_546(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func G_expect_12_73(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func F_expect_10_91(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func E_expect_9_09(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func D_expect_7_27(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func C_expect_5_46(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func B_expect_3_64(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func A_expect_1_82(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func main() {
	var q uint64
	multiplier := flag.Uint64("m", 100000, "multiplier")
	flag.Parse()
	//fmt.Println("multiplier=", *multiplier)
	file, err := os.Create("serial_prof")
	if err != nil {
		log.Fatal(err)
	}
	if err = pprof.StartCPUProfile(file); err != nil {
		log.Fatal(err)
	}
	mult := *multiplier

	for i := uint64(0); i < 100; i++ {
		f := i + A_expect_1_82(0xebabefac23, 1*mult)
		g := i + B_expect_3_64(f, 2*mult)
		h := i + C_expect_5_46(g, 3*mult)
		k := i + D_expect_7_27(h, 4*mult)
		l := i + E_expect_9_09(k, 5*mult)
		m := i + F_expect_10_91(l, 6*mult)
		n := i + G_expect_12_73(m, 7*mult)
		o := i + H_expect_14_546(n, 8*mult)
		p := i + I_expect_16_36(o, 9*mult)
		q = i + J_expect_18_18(p, 10*mult)
	}
	pprof.StopCPUProfile()
	file.Close()
	fmt.Println(q)
}
