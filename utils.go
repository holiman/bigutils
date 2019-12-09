package bigutils

import (
	"math/big"
	"math/bits"
)

const _S = bits.UintSize / 8 // word size in bytes

// MarshallBigInt writes 'b' to the buffer 'buf'. It's equivalent to
// `copy(buf, b.Bytes())`, plus it also returns the number of bytes written.
// If you want to know an _upper bound_ on the space needed in 'buf', you need to
// use `MaxByteLen(b)`.
func MarshallBigInt(b *big.Int, buf []byte) int {
	switch bits.UintSize {
	case 32:
		return marshallBigint32bit(b, buf)
	case 64:
		return marshallBigint64bit(b, buf)
	default:
		return marshallBigintGeneric(b, buf)
	}
}

// marshallBigintGeneric is a platform-independent implementation of MarshallBigInt
func marshallBigintGeneric(b *big.Int, buf []byte) int {
	z := b.Bits()
	if len(z) == 0 {
		return 0
	}
	// We want to skip any empty MSB bytes in the first word
	var (
		i      int
		j      = (_S - 1) * 8
		zindex = len(z) - 1
		d      = z[zindex]
	)
	// z consists of Words (either 64 or 32 bit ints). We need to start by taking
	// the most significant word, and skip the empty bytes of that word
	for ; j >= 0; j -= 8 {
		if data := byte(d >> j); data != 0 {
			buf[i] = data
			i++
			j -= 8
			break
		}
	}
collect:
	for ; j >= 0; j -= 8 {
		buf[i] = byte(d >> j)
		i++
	}
	// At this point, first Word is done
	zindex--
	if zindex >= 0 {
		d = z[zindex]
		j = (_S - 1) * 8
		goto collect
	}
	return i
}

func MaxByteLen(b *big.Int) int {
	return len(b.Bits()) * _S
}

// marshallBigint64bit is an implementation of MarshallBigInt for 64-bit platforms
func marshallBigint64bit(b *big.Int, buf []byte) int {
	z := b.Bits()
	if len(z) == 0 {
		return 0
	}
	// z consists of Words (either 64 or 32 bit ints). We need to start by taking
	// the most significant word, and skip the empty bytes of that word
	var (
		i      int          // index if where we are in buffer
		zindex = len(z) - 1 // index to read from
		d      = z[zindex]  // the Word currently being read
	)
	// Figure out how much we can sikip
	switch {
	case byte(d>>56) != 0:
		goto a
	case byte(d>>48) != 0:
		goto b
	case byte(d>>40) != 0:
		goto c
	case byte(d>>32) != 0:
		goto d
	case byte(d>>24) != 0:
		goto e
	case byte(d>>16) != 0:
		goto f
	case byte(d>>8) != 0:
		goto g
	default:
		goto h
	}

a:
	buf[i], i = byte(d>>56), i+1
b:
	buf[i], i = byte(d>>48), i+1
c:
	buf[i], i = byte(d>>40), i+1
d:
	buf[i], i = byte(d>>32), i+1
e:
	buf[i], i = byte(d>>24), i+1
f:
	buf[i], i = byte(d>>16), i+1
g:
	buf[i], i = byte(d>>8), i+1
h:
	buf[i], i = byte(d>>0), i+1
	// At this point, a Word is done, continue with the next
	if zindex > 0 {
		zindex--
		d = z[zindex]
		goto a
	}
	return i
}

// marshallBigint32bit is an implementation of MarshallBigInt for 32-bit platforms
func marshallBigint32bit(b *big.Int, buf []byte) int {
	z := b.Bits()
	if len(z) == 0 {
		return 0
	}
	// z consists of Words (either 64 or 32 bit ints). We need to start by taking
	// the most significant word, and skip the empty bytes of that word
	var (
		i      int          // index if where we are in buffer
		zindex = len(z) - 1 // index to read from
		d      = z[zindex]  // the Word currently being read
	)
	// Figure out how much we can sikip
	switch {
	case byte(d>>24) != 0:
		goto e
	case byte(d>>16) != 0:
		goto f
	case byte(d>>8) != 0:
		goto g
	default:
		goto h
	}

e:
	buf[i], i = byte(d>>24), i+1
f:
	buf[i], i = byte(d>>16), i+1
g:
	buf[i], i = byte(d>>8), i+1
h:
	buf[i], i = byte(d>>0), i+1
	// At this point, a Word is done, continue with the next
	if zindex > 0 {
		zindex--
		d = z[zindex]
		goto e
	}
	return i
}

// U256 encodes b as a 256 bit two's complement number. This operation is destructive.
// It is semantically equivalent to `b = b && (2^256-1)`
func U256(b *big.Int) {
	if bits.UintSize == 64 {
		u256_64bit(b)
		return
	}
	if bits.UintSize == 32 {
		u256_32bit(b)
		return
	}
	u256_generic(b)
}

func u256_64bit(x *big.Int) {
	z := x.Bits()
	// z consists of Words, which are uint64 on this platforms. We just have to
	// ensure that there are not more than 4 Words
	if len := len(z); len > 4 {
		x.SetBits(z[len-4:])
	}
}

func u256_32bit(x *big.Int) {
	z := x.Bits()
	// z consists of Words, which are uint32 on this platforms. We just have to
	// ensure that there are not more than 8 Words
	if len := len(z); len > 8 {
		x.SetBits(z[len-8:])
	}
}
func u256_generic(x *big.Int) {
	z := x.Bits()
	// z consists of Words, which are size _S  on this platforms. We just have to
	// ensure that there are not more than 256/_S Words
	nWords := 256 / _S
	if len := len(z); len > nWords {
		x.SetBits(z[len-nWords:])
	}
}
