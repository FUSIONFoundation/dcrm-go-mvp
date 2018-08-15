package dcrm 

import (
    "math/big"
    "math/rand"
    "time"
    crand"crypto/rand"
)

func get_rand_int(bitlen uint) *big.Int {
	one,_ := new(big.Int).SetString("1",10)
	zz := new(big.Int).Lsh(one,bitlen)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	z := new(big.Int).Rand(rnd,zz) //[0,zz)
	return z
}

func randomFromZn(n *big.Int,rnd *rand.Rand) *big.Int {

    var result *big.Int
    for {
	    result = get_rand_int(uint(n.BitLen()))
	    r := result.Cmp(n)
	    if r < 0 {
		break
	    }
    }

    return result
}

func randomFromZnStar(n *big.Int,rnd *rand.Rand) *big.Int {
    result,_ := crand.Prime(crand.Reader,n.BitLen())
    return result
}

func isElementOfZn(element *big.Int,n *big.Int) bool {
    zero,_ := new(big.Int).SetString("0",10)
    r := element.Cmp(zero)
    rr := element.Cmp(n)

    if r >= 0 && rr < 0 {
	return true
    }

    return false
}

