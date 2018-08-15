package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/dcrm"
)

func main() {
    userCount := 4
    msgHash := "8f481706b85902c502aef1e4ec0cfe614565ee6f9675d2603a720aaa2c71e65d"

    userList := dcrm.KeyGenerate(int32(userCount))
    encX := ((userList.Front().Value).(*dcrm.User)).GetEncX()
    pkx := ((userList.Front().Value).(*dcrm.User)).GetPk_x()
    pky := ((userList.Front().Value).(*dcrm.User)).GetPk_y()
    signature := dcrm.Sign(userList,encX,msgHash)
    fmt.Println("nECDSA Signature is (r,s,v):\n")
    fmt.Println("r is:\n",signature.GetR())
    fmt.Println("s is:\n",signature.GetS())
    fmt.Println("v is:\n",signature.GetRecoveryParam())

    if signature != nil {
	dcrm.Verify(signature,msgHash,pkx,pky)
	return
    }
    
    fmt.Println("ECDSA Signature Verify NOT Passed!\n")
}

