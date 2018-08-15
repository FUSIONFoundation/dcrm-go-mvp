package dcrm 

import (
    	crand"crypto/rand"
	"math/rand"
	"time"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var (
    SecureRnd = rand.New(rand.NewSource(time.Now().UnixNano()))
    //paillier
    privKey,_ = GenerateKey(crand.Reader, 1024)
    //zk
    ZKParams = generatePublicParams(secp256k1.S256(), 256, 512, SecureRnd, &privKey.PublicKey)
    //commitment
    MPK = generateMasterPK()
)
