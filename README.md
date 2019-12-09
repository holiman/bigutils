## Big Utils

This library contains some helper functionality for `big.Int`, the golang-native
library for arbitrary precision integers. 

In particular, the `big.Int` library lacks a couple of features:

- Fast writing to a pre-allocated buffer, 
- Fast `and` to a `256` bit space

## Benchmarks

### Marshalling

The `0-256bits-prealloc-64bit-6` writes `big.Int`s ranging from zero to 256 bit in size, 
to a preallocated buffer in `3471ns`. 
The `0-256bits-noprealloc-bigint-6` is the `big.Int` native `Bytes` implementation, which 
clocks in at `11277ns`, with quite a few allocs aswell. 

```
BenchmarkMarshallBigint/0-256bits-prealloc-generic-6         	  129177	      8259 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshallBigint/0-256bits-prealloc-64bit-6           	  345541	      3471 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshallBigint/0-256bits-prealloc-platform-6        	  302342	      4040 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshallBigint/0-256bits-noprealloc-bigint-6        	  100200	     11277 ns/op	    5680 B/op	     257 allocs/op
BenchmarkMarshallBigint/0-256bits-noprealloc-generic-6       	   72660	     14502 ns/op	    5680 B/op	     257 allocs/op
BenchmarkMarshallBigint/0-256bits-noprealloc-64bit-6         	  123219	      9611 ns/op	    5680 B/op	     257 allocs/op
```

### Masking

The `0-256bits-64bit-6` masks `big.Int`s ranging from zero to 256 bits in size, on a
`64bit` architecture. It does so in `695ns`. 
The reference implementation uses `big.And`, and clocks in at `3646ns`
```
BenchmarkU256/0-256bits-generic-6         	 1329020	       770 ns/op
BenchmarkU256/0-256bits-64bit-6           	 1704949	       695 ns/op
BenchmarkU256/0-256bits-bigint-6          	  336811	      3646 ns/op
```