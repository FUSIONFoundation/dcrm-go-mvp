package dcrm 

import (
    "math/big"
    "github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type PublicParameters struct {
    h1 *big.Int
    h2 *big.Int
    nTilde *big.Int
    paillierPubKey *PublicKey
}

func (this *PublicParameters) New(BitCurve *secp256k1.BitCurve,nTilde *big.Int,kPrime int32,h1 *big.Int,h2 *big.Int,paillierPubKey *PublicKey) {
	this.nTilde = nTilde
	this.h1 = h1
	this.h2 = h2
	this.paillierPubKey = paillierPubKey

	if BitCurve == nil {
	    return//test
	}
}
