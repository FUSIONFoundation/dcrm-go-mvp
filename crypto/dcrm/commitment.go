package dcrm 

import (
	"github.com/ethereum/go-ethereum/crypto/pbc"
)

type Commitment struct {
	committment *pbc.Element
	pubkey *pbc.Element
}

func (ct *Commitment) New(pubkey *pbc.Element,a *pbc.Element) {
    ct.pubkey = pubkey
    ct.committment = a
}

