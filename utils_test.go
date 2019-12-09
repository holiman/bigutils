package bigutils

import (
	"bytes"
	"math/big"
	"testing"
)

func TestMarshallBigint(t *testing.T) {
	t.Run("generic", func(t *testing.T) {
		testMarshallBigint(t, marshallBigintGeneric)
	})
	t.Run("64bit", func(t *testing.T) {
		testMarshallBigint(t, marshallBigint64bit)
	})
}
func testMarshallBigint(t *testing.T, marshaller func(*big.Int, []byte) int) {
	bigint := big.NewInt(0)
	for i := 0; i < 512; i++ {
		buf := make([]byte, 100)
		len := marshaller(bigint, buf)
		got := buf[:len]
		exp := bigint.Bytes()
		if !bytes.Equal(got, exp) {
			t.Errorf("marshalling error:"+
				" value: %v\n"+
				" bits: %x\n"+
				" exp: %x\n"+
				" got: %x\n",
				bigint, bigint.Bits(), exp, got)
		}
		if bigint.Uint64() == 0 {
			bigint.SetUint64(1)
		} else {
			bigint.Mul(bigint, big.NewInt(2))
		}
	}
}

//
//BenchmarkMarshallBigint/0-256bits-prealloc-generic-6         	  129177	      8259 ns/op	       0 B/op	       0 allocs/op
//BenchmarkMarshallBigint/0-256bits-prealloc-64bit-6           	  345541	      3471 ns/op	       0 B/op	       0 allocs/op
//BenchmarkMarshallBigint/0-256bits-prealloc-platform-6        	  302342	      4040 ns/op	       0 B/op	       0 allocs/op
//BenchmarkMarshallBigint/0-256bits-noprealloc-bigint-6        	  100200	     11277 ns/op	    5680 B/op	     257 allocs/op
//BenchmarkMarshallBigint/0-256bits-noprealloc-generic-6       	   72660	     14502 ns/op	    5680 B/op	     257 allocs/op
//BenchmarkMarshallBigint/0-256bits-noprealloc-64bit-6         	  123219	      9611 ns/op	    5680 B/op	     257 allocs/op
//
func BenchmarkMarshallBigint(b *testing.B) {
	var bigints = []*big.Int{
		new(big.Int),
	}
	bigint := big.NewInt(1)
	bigints = append(bigints, bigint)
	for i := 0; i < 256; i++ {
		bigint = new(big.Int).Mul(bigint, big.NewInt(2))
		bigints = append(bigints, bigint)
	}
	// Set bit i
	buf := make([]byte, 100)
	b.Run("0-256bits-prealloc-generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				marshallBigintGeneric(bigint, buf)
			}
		}
	})
	b.Run("0-256bits-prealloc-64bit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				marshallBigint64bit(bigint, buf)
			}
		}
	})
	b.Run("0-256bits-prealloc-platform", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				MarshallBigInt(bigint, buf)
			}
		}
	})
	b.Run("0-256bits-noprealloc-bigint", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				bigint.Bytes()
			}
		}
	})
	b.Run("0-256bits-noprealloc-generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				buf := make([]byte, MaxByteLen(bigint))
				marshallBigintGeneric(bigint, buf)
			}
		}
	})
	b.Run("0-256bits-noprealloc-64bit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				buf := make([]byte, MaxByteLen(bigint))
				marshallBigint64bit(bigint, buf)
			}
		}
	})
}

func bigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

var (
	tt256   = bigPow(2, 256)
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1))
)

func u256_ReferenceImpl(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

func TestUint256(t *testing.T) {
	bigint := big.NewInt(0)
	for i := 0; i < 512; i++ {
		got := new(big.Int).SetBytes(bigint.Bytes())
		exp := new(big.Int).SetBytes(bigint.Bytes())
		U256(got)
		u256_ReferenceImpl(exp)
		if got.Cmp(exp) != 0 {
			t.Errorf("error: value: %v\n"+
				" exp: %x\n"+
				" got: %x\n",
				bigint, exp, got)
		}
		if bigint.Uint64() == 0 {
			bigint.SetUint64(1)
		} else {
			bigint.Mul(bigint, big.NewInt(2))
		}
	}
}

//BenchmarkU256/0-256bits-6         	 1536889	       715 ns/op
//BenchmarkU256/0-256bits-bigint-6  	  314708	      3850 ns/op
func BenchmarkU256(b *testing.B) {
	var mkVectors = func() []*big.Int {
		var bigints = []*big.Int{
			new(big.Int),
		}
		bigint := big.NewInt(1)
		bigints = append(bigints, bigint)
		for i := 0; i < 256; i++ {
			bigint = new(big.Int).Mul(bigint, big.NewInt(2))
			bigints = append(bigints, bigint)
		}
		return bigints
	}
	bigints := mkVectors()
	b.Run("0-256bits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				U256(bigint)
			}
		}
	})
	b.Run("0-256bits-bigint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, bigint := range bigints {
				u256_ReferenceImpl(bigint)
			}
		}
	})
}
