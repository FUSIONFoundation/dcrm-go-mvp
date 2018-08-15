package dcrm 

import (
	"math/big"
	"github.com/ethereum/go-ethereum/crypto/pbc"
)

type Open struct {
    secrets []*big.Int
    randomness *pbc.Element
}

func (open *Open) New(randomness *pbc.Element,secrets []*big.Int) {
    open.randomness = randomness
    open.secrets = secrets //test
}

func (open *Open) getSecrets() []*big.Int {
    return open.secrets
}

func (open *Open) getRandomness() *pbc.Element {
    return open.randomness
}

