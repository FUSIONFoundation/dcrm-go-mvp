package dcrm

/*
//#cgo linux CFLAGS: -I../../common/math/libtommath-0.41
//#cgo linux LDFLAGS: -L../../common/math/libtommath-0.41 -ltommath
//#cgo linux LDFLAGS: -L./
//#include "../../common/math/libtommath-0.41/tommath.h"
//#include </usr/include/tommath.h>
//#include "../../common/math/libtommath-0.41/my_ecc.c"
#cgo linux CFLAGS: -I../../common/math/gmp-6.1.2
#cgo linux LDFLAGS: -L../../common/math/gmp-6.1.2 -lgmp 
#include "../../common/math/gmp-6.1.2/gmp.h"
#include "../../common/math/gmp-6.1.2/mpz/my_gmp.c"
#cgo linux CFLAGS: -I../../common/math/pbc-0.5.14 -std=gnu99
#cgo linux CFLAGS: -I../../common/math/pbc-0.5.14/include
#cgo linux LDFLAGS: -L../../common/math/pbc-0.5.14 -lpbc 
#include "../../common/math/pbc-0.5.14/ecc/pairing.c"*/
import "C"

import (
	"container/list"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"fmt"
	"unsafe"
	"github.com/ethereum/go-ethereum/common/math/decimal"
	"github.com/ethereum/go-ethereum/common/math"
)

//===============gmp============================
/*func getPrime(length int32,seed int64) *big.Int {
	bitsize := length//(length + 7)/8
	p := make([]byte, bitsize)
	for i := range p {
		p[i] = 0
	}
	pPtr := (*C.char)(unsafe.Pointer(&p[0]))

	plen := new(int32)
	*plen = length
	plenPtr := (*C.int)(unsafe.Pointer(plen))

	pseed := new(int64)
	*pseed = seed
	pseedPtr := (*C.long)(unsafe.Pointer(pseed))

	C.getPrime(pPtr,plenPtr,pseedPtr)
	s := string(p[0:*plen])
	//fmt.Println("p count is: %s\n",s)
	prime,_ := new(big.Int).SetString(s,10)
	//fmt.Println("prime is: %v\n",prime)
	return prime
}*/

func isProbablePrime(num *big.Int) bool {
	
	ret := new(int)
	retPtr := (*C.int)(unsafe.Pointer(ret))
	
	numlen := new(int)
	*numlen = (num.BitLen()+7)/8
	nums := make([]byte,*numlen)
	for i := range nums {
		nums[i] = 0
	}
	math.ReadBits(num,nums[:])

	numPtr := (*C.char)(unsafe.Pointer(&nums[0]))
	numlenPtr := (*C.int)(unsafe.Pointer(numlen))
	C.isProbablePrime(retPtr,numPtr,numlenPtr)

	if (*ret) == 1 {
	    return true
	}

	return false
}

func newBigIntStrimZero(signum int,val []byte) *big.Int {
	off := 0;
	lv := len(val)

	if signum < -1 || signum > 1 {
	    return nil
        }
	
	if (off < 0) || (lv < 0) ||
            (lv > 0 && off >= lv) { 
            return nil
        }

	ret := new(big.Int).SetBytes(val[:])
	dec := decimal.NewFromBigInt(ret,0)
	str := dec.String()
	ret2,_ := new(big.Int).SetString(str,10)

	if ret2.Sign() != signum {
	    ret3 := new(big.Int).Neg(ret2)
	    //test
	    zero,_ := new(big.Int).SetString("0",10)
	    if ret3.Cmp(zero)  < 0 {
		fmt.Println("======newBigIntStrimZero is < 0=====")
	    }
	    //test
	    return ret3
	}
	
	return ret2
}

func modInverse(val *big.Int,modulus *big.Int) (*big.Int,error) {
	ret := new(int)
	retPtr := (*C.int)(unsafe.Pointer(ret))
	
	got_count := new(int)
	got_countPtr := (*C.int)(unsafe.Pointer(got_count))
	
	got_datalen := ((modulus.BitLen()+7)/8)+1
	got_data := make([]byte, got_datalen)
	for i := range got_data {
		got_data[i] = 0
	}
	got_dataPtr := (*C.char)(unsafe.Pointer(&got_data[0]))

	vallen := new(int)
	*vallen = ((val.BitLen()+7)/8)+1
	vals := make([]byte,*vallen)
	for i := range vals {
		vals[i] = 0
	}
	math.ReadBits(val, vals[0:(*vallen)-1])

	modlen := new(int)
	*modlen = ((modulus.BitLen()+7)/8)+1
	moduluss := make([]byte,*modlen)
	for i := range moduluss {
		moduluss[i] = 0
	}
	math.ReadBits(modulus, moduluss[0:(*modlen)-1])
	
	valPtr := (*C.char)(unsafe.Pointer(&vals[0]))
	modulussPtr := (*C.char)(unsafe.Pointer(&moduluss[0]))
	vallenPtr := (*C.int)(unsafe.Pointer(vallen))
	modlenPtr := (*C.int)(unsafe.Pointer(modlen))
	C.modInverse(retPtr,got_countPtr,got_dataPtr,valPtr,vallenPtr,modulussPtr,modlenPtr)

	retvalue := newBigIntStrimZero(*ret,got_data[0:*got_count])

	return retvalue,nil
}

func modPowSecure(base *big.Int,exponent *big.Int,modulus *big.Int) *big.Int {

	ret := new(int)
	retPtr := (*C.int)(unsafe.Pointer(ret))
	
	got_count := new(int)
	got_countPtr := (*C.int)(unsafe.Pointer(got_count))
	
	got_datalen := ((modulus.BitLen()+7)/8)+1
	got_data := make([]byte, got_datalen)
	for i := range got_data {
		got_data[i] = 0
	}
	got_dataPtr := (*C.char)(unsafe.Pointer(&got_data[0]))

	baselen := new(int)
	*baselen = ((base.BitLen()+7)/8)+1
	bases := make([]byte,*baselen)
	for i := range bases {
		bases[i] = 0
	}
	math.ReadBits(base, bases[0:(*baselen)-1])

	explen := new(int) 
	*explen = ((exponent.BitLen()+7)/8)+1
	exponents := make([]byte,*explen)
	for i := range exponents {
		exponents[i] = 0
	}
	math.ReadBits(exponent, exponents[0:(*explen)-1])

	modlen := new(int)
	*modlen = ((modulus.BitLen()+7)/8)+1
	moduluss := make([]byte,*modlen)
	for i := range moduluss {
		moduluss[i] = 0
	}
	math.ReadBits(modulus, moduluss[0:(*modlen)-1])
	
	basesPtr := (*C.char)(unsafe.Pointer(&bases[0]))
	exponentsPtr := (*C.char)(unsafe.Pointer(&exponents[0]))
	modulussPtr := (*C.char)(unsafe.Pointer(&moduluss[0]))
	baselenPtr := (*C.int)(unsafe.Pointer(baselen))
	explenPtr := (*C.int)(unsafe.Pointer(explen))
	modlenPtr := (*C.int)(unsafe.Pointer(modlen))
	C.modPowSecure(retPtr,got_countPtr,got_dataPtr,basesPtr,baselenPtr,exponentsPtr,explenPtr,modulussPtr,modlenPtr)

	retvalue := newBigIntStrimZero(*ret,got_data[0:*got_count])

	return retvalue
}

func modPow(base *big.Int,exponent *big.Int,modulus *big.Int) *big.Int {
	zero2,_ := new(big.Int).SetString("0",10)
	if exponent.Cmp(zero2) >= 0 {
	    return new(big.Int).Exp(base,exponent,modulus)
	}

	z := new(big.Int).ModInverse(base,modulus)
	exp := new(big.Int).Neg(exponent)
	return new(big.Int).Exp(z,exp,modulus)
}
//============================================

func KeyGenerate(userCnt int32) *list.List {

    userList := list.New()
    kgRoundOne(userList,userCnt)
    kgRoundTwo(userList)
    kgRoundThree(userList)
    return userList
}

func kgRoundOne(userList *list.List,userCnt int32) {
    var i int32
    for i =0;i < userCnt;i++ {
	temUser := new(User)
	xShare := randomFromZn(secp256k1.S256().N, SecureRnd)
	fmt.Println("======kgRoundOne=====xShare is \n",xShare)
	if xShare.Sign() == -1 {
		xShare.Add(xShare,secp256k1.S256().P)
	}
	kg := make([]byte, 32)
	math.ReadBits(xShare, kg[:])
	kgx0,kgy0 := secp256k1.KMulG(kg[:])

	xShareRnd := randomFromZnStar((&privKey.PublicKey).N,SecureRnd)

	encXShare := encrypt((&privKey.PublicKey),xShare, xShareRnd)

	yShares := secp256k1.S256().Marshal(kgx0,kgy0)
	
	var nums = []*big.Int{encXShare,new(big.Int).SetBytes(yShares[:])}
	mpkEncXiYi := multiLinnearCommit(SecureRnd,MPK,nums)
	openEncXiYi := mpkEncXiYi.cmtOpen()
	cmtEncXiYi := mpkEncXiYi.cmtCommitment()
	fmt.Println("----Info: User i calculate Commitment in KG round ONE---\n",i)

	temUser.setxShare(xShare)
	temUser.setyShare_x(kgx0)
	temUser.setyShare_y(kgy0)
	temUser.setxShareRnd(xShareRnd)
	temUser.setEncXShare(encXShare)
	temUser.setMpkEncXiYi(mpkEncXiYi)
	temUser.setOpenEncXiYi(openEncXiYi)
	temUser.setCmtEncXiYi(cmtEncXiYi)
	userList.PushBack(temUser)
    }
}

func kgRoundTwo(userList *list.List) {
    	a := 0
	for e := userList.Front(); e != nil; e = e.Next() {
	    zkpKG := new(ZkpKG)
	    zkpKG.New(ZKParams,((e.Value).(*User)).getxShare(),SecureRnd,secp256k1.S256().Gx,secp256k1.S256().Gy,((e.Value).(*User)).getEncXShare(), ((e.Value).(*User)).getxShareRnd())

	    ((e.Value).(*User)).setZkpKG(zkpKG)

	fmt.Println("----Info: User i calculate Zero-Knowledge in KG round TWO---\n",a)
	a = a+1
	}
}

func kgRoundThree(userList *list.List) {

    a := 0

    //1 check commitment
    for e := userList.Front(); e != nil; e = e.Next() {
	if (checkcommitment(((e.Value).(*User)).getCmtEncXiYi(), ((e.Value).(*User)).getOpenEncXiYi(), MPK) == false) {
	    fmt.Println("##Error####################: KG Round 3, User i does not pass checking Commitment!\n",a)
	}

	a = a+1
    }

    //verify zk
    a = 0
    for e := userList.Front(); e != nil; e = e.Next() {

	rr := ((e.Value).(*User)).getOpenEncXiYi().getSecrets()[1]
	rrlen := ((rr.BitLen()+7)/8)
	rrs := make([]byte,rrlen)
	math.ReadBits(rr,rrs[:])
	rx,ry := secp256k1.S256().Unmarshal(rrs[:])

	if (((e.Value).(*User)).getZkpKG().verify(ZKParams,secp256k1.S256(),rx,ry,((e.Value).(*User)).getOpenEncXiYi().getSecrets()[0]) == false) {
	    fmt.Println("##Error####################: KG Round 3, User i does not pass verifying Zero-Knowledge!\n",a)
	}
	a = a+1
    }
    
    for e := userList.Front(); e != nil; e = e.Next() {
	((e.Value).(*User)).setEncX(calculateEncPrivateKey(userList));
	pkx,pky := calculatePubKey(userList)
	((e.Value).(*User)).setPk_x(pkx);
	((e.Value).(*User)).setPk_y(pky);
    }
}

func calculatePubKey(userList *list.List) (*big.Int,*big.Int) {
    yShare_x0 := ((userList.Front().Value).(*User)).getyShare_x()
    yShare_y0 := ((userList.Front().Value).(*User)).getyShare_y()

    e := userList.Front()
    for e = e.Next(); e != nil; e = e.Next() {
	yShare_xi := ((e.Value).(*User)).getyShare_x()
	yShare_yi := ((e.Value).(*User)).getyShare_y()

	yShare_x0,yShare_y0 = secp256k1.S256().Add(yShare_x0,yShare_y0,yShare_xi,yShare_yi) 
    }

    return yShare_x0,yShare_y0
}

func calculateEncPrivateKey(userList *list.List) *big.Int {

    encX := ((userList.Front().Value).(*User)).getEncXShare()
    e := userList.Front()
    for e = e.Next(); e != nil; e = e.Next() {
	encXi := ((e.Value).(*User)).getEncXShare()
	encX = cipherAdd((&privKey.PublicKey),encX,encXi);
    }
    
    fmt.Println("---Info: Calculate the Encrypted Private Key,EncPrivateKey: ---\n",encX)

    return encX
}

func signRoudOne(userList *list.List,encX *big.Int)  {
	var rhoI, rhoIRnd, uI, vI *big.Int
	var mpkUiVi *MTDCommitment
	var openUiVi *Open
	var cmtUiVi *Commitment
	a := 0 

	for e := userList.Front(); e != nil; e = e.Next() {
	    rhoI = randomFromZn(secp256k1.S256().N, SecureRnd)
	    rhoIRnd = randomFromZnStar((&privKey.PublicKey).N,SecureRnd)
	    uI = encrypt((&privKey.PublicKey),rhoI, rhoIRnd)
	    vI = cipherMultiply((&privKey.PublicKey),encX, rhoI)
	    
	    var nums = []*big.Int{uI,vI}
	    mpkUiVi = multiLinnearCommit(SecureRnd,MPK,nums)
	    openUiVi = mpkUiVi.cmtOpen()
	    cmtUiVi = mpkUiVi.cmtCommitment()

	    ((e.Value).(*User)).setRhoI(rhoI)
	    ((e.Value).(*User)).setRhoIRnd(rhoIRnd)
	    ((e.Value).(*User)).setuI(uI)
	    ((e.Value).(*User)).setvI(vI)
	    ((e.Value).(*User)).setMpkUiVi(mpkUiVi)
	    ((e.Value).(*User)).setOpenUiVi(openUiVi)
	    ((e.Value).(*User)).setCmtUiVi(cmtUiVi)
	    
	    fmt.Println("---Info: User i calculate Commitment in round ONE---\n",a)
	    a = a+1

    }

}

func signRoudTwo(userList *list.List,encX *big.Int)  {
	a := 0 
	for e := userList.Front(); e != nil; e = e.Next() {
	    zkp1 := new(ZkpSignOne)
	    zkp1.New(ZKParams,((e.Value).(*User)).getRhoI(),SecureRnd,((e.Value).(*User)).getRhoIRnd(),((e.Value).(*User)).getvI(), encX,((e.Value).(*User)).getuI())

	    ((e.Value).(*User)).setZkp1(zkp1)
	    fmt.Println("---Info: User i calculate Zero-Knowledge in round TWO---\n",a)
	    a = a+1
	}
}

func calculateU(userList *list.List) *big.Int {
    u := ((userList.Front().Value).(*User)).getOpenUiVi().getSecrets()[0]
    e := userList.Front()
    for e = e.Next(); e != nil; e = e.Next() {
	ui := ((e.Value).(*User)).getOpenUiVi().getSecrets()[0]
	u = cipherAdd((&privKey.PublicKey),u,ui)
    }
    
    fmt.Println("---Info: Calculate the Encrypted Inner-Data u, U:---\n",u)

    return u
}

func calculateV(userList *list.List) *big.Int {
    v := ((userList.Front().Value).(*User)).getOpenUiVi().getSecrets()[1]
    e := userList.Front()
    for e = e.Next(); e != nil; e = e.Next() {
	vi := ((e.Value).(*User)).getOpenUiVi().getSecrets()[1]
	v = cipherAdd((&privKey.PublicKey),v,vi)
    }
    
    fmt.Println("---Info: Calculate the Encrypted Inner-Data v, V:---\n",v)

    return v
}

func signRoundThree(userList *list.List,encX *big.Int) bool  {
    
    aborted := false

    //1 check commitment
    a := 0
    for e := userList.Front(); e != nil; e = e.Next() {
    if checkcommitment(((e.Value).(*User)).getCmtUiVi(), ((e.Value).(*User)).getOpenUiVi(),MPK) == false {
	    aborted = true
	    fmt.Println("##Error####################: SignRound 3, User i does not pass checking Commitment! \n",a);
	    return aborted
	    }
	    
	    a = a+1
    }

    //2 verify zk
    b := 0
    for e := userList.Front(); e != nil; e = e.Next() {
	if ((e.Value).(*User)).getZkp1().verify(ZKParams,secp256k1.S256(), ((e.Value).(*User)).getOpenUiVi().getSecrets()[1],encX,((e.Value).(*User)).getOpenUiVi().getSecrets()[0]) == false {
		aborted = true				
		fmt.Println("##Error####################: SignRound 3, User i does not pass verifying Zero-Knowledge!",b);
		return aborted
	    }

	b = b+1
    }

    u := calculateU(userList)
    v := calculateV(userList)
    tttt,_ := new(big.Int).SetString("0",10)
    if v.Cmp(tttt) != 0 {//test
    }

    //3
    a = 0
    for e := userList.Front(); e != nil; e = e.Next() {
	kI := randomFromZn(secp256k1.S256().N, SecureRnd)
	if kI.Sign() == -1 {
		kI.Add(kI,secp256k1.S256().P)
	}
	rI := make([]byte, 32)
	math.ReadBits(kI, rI[:])
	rIx,rIy := secp256k1.KMulG(rI[:])
	cI := randomFromZn(secp256k1.S256().N, SecureRnd)
	cIRnd := randomFromZnStar((&privKey.PublicKey).N,SecureRnd)
	mask := encrypt((&privKey.PublicKey),new(big.Int).Mul(secp256k1.S256().N, cI),cIRnd)
	wI := cipherAdd((&privKey.PublicKey),cipherMultiply((&privKey.PublicKey),u, kI), mask)
	///
	rIs := secp256k1.S256().Marshal(rIx,rIy)
	
	var nums = []*big.Int{new(big.Int).SetBytes(rIs[:]),wI}
	mpkRiWi := multiLinnearCommit(SecureRnd,MPK,nums)

	openRiWi := mpkRiWi.cmtOpen()
	cmtRiWi := mpkRiWi.cmtCommitment()
	((e.Value).(*User)).setkI(kI)
	((e.Value).(*User)).setcI(cI)
	((e.Value).(*User)).setcIRnd(cIRnd)
	((e.Value).(*User)).setrI_x(rIx)
	((e.Value).(*User)).setrI_y(rIy)
	((e.Value).(*User)).setMask(mask)
	((e.Value).(*User)).setwI(wI)
	((e.Value).(*User)).setMpkRiWi(mpkRiWi)
	((e.Value).(*User)).setOpenRiWi(openRiWi)
	((e.Value).(*User)).setCmtRiWi(cmtRiWi)
	fmt.Println("--Info: User i calculate Commitment in Signning round THREE!",a);

	a = a+1
    }

    return aborted
}

func signRoundFour(userList *list.List,u *big.Int) {
    a := 0 
    for e := userList.Front(); e != nil; e = e.Next() {
	zkp2 := new(ZkpSignTwo)
	zkp2.New(ZKParams,((e.Value).(*User)).getkI(),((e.Value).(*User)).getcI(),SecureRnd,secp256k1.S256().Gx,secp256k1.S256().Gy,((e.Value).(*User)).getwI(),u,((e.Value).(*User)).getcIRnd())

	((e.Value).(*User)).setZkp_i2(zkp2)
	fmt.Println("--Info: User i calculate Zero-Knowledge in Signning round FOUR\n",a)
	a = a+1
    }
}

func calculateW(userList *list.List) *big.Int {
    w := ((userList.Front().Value).(*User)).getOpenRiWi().getSecrets()[1]
    e := userList.Front()
    for e = e.Next(); e != nil; e = e.Next() {
	wi := ((e.Value).(*User)).getOpenRiWi().getSecrets()[1]
	w = cipherAdd((&privKey.PublicKey),w,wi);
    }
    
    fmt.Println("---Info: Calculate the Encrypted Inner-Data w: ---\n",w)

    return w
}

func calculateR(userList *list.List) (*big.Int,*big.Int) {

	rr := ((userList.Front().Value).(*User)).getOpenRiWi().getSecrets()[0]
	rrlen := ((rr.BitLen()+7)/8)
	rrs := make([]byte,rrlen)
	math.ReadBits(rr,rrs[:])
	rx,ry := secp256k1.S256().Unmarshal(rrs[:])

	e := userList.Front()
	for e = e.Next(); e != nil; e = e.Next() {

	rri := ((e.Value).(*User)).getOpenRiWi().getSecrets()[0]
	rrilen := ((rri.BitLen()+7)/8)
	rris := make([]byte,rrilen)
	math.ReadBits(rri,rris[:])
	rrix,rriy := secp256k1.S256().Unmarshal(rris[:])

	rx,ry = secp256k1.S256().Add(rx,ry,rrix,rriy)
	}

	fmt.Println("---Info: Calculate the Encrypted Inner-Data R(rx,ry): ---\n",rx,ry)

	return rx,ry
}

func signRoundFive(userList *list.List,u *big.Int,v *big.Int,message string) *ECDSASignature {
    signature := new(ECDSASignature)
    signature.New()

    aborted := false

    //1 check commitment
    a := 0
    for e := userList.Front(); e != nil; e = e.Next() {
    if checkcommitment(((e.Value).(*User)).getCmtRiWi(), ((e.Value).(*User)).getOpenRiWi(),MPK) == false {
	    aborted = true
	    fmt.Println("##Error####################: SignRound 5, User i does not pass checking Commitment! \n",a);
	    signature.setRoudFiveAborted(aborted)
	    }
	    
	    a = a+1
    }

    //2 verify zk
    b := 0
    for e := userList.Front(); e != nil; e = e.Next() {

	rr := ((e.Value).(*User)).getOpenRiWi().getSecrets()[0]
	rrlen := ((rr.BitLen()+7)/8)
	rrs := make([]byte,rrlen)
	math.ReadBits(rr,rrs[:])
	rx,ry := secp256k1.S256().Unmarshal(rrs[:])

	if ((e.Value).(*User)).getZkp_i2().verify(ZKParams,secp256k1.S256(), rx,ry,u,((e.Value).(*User)).getOpenRiWi().getSecrets()[1]) == false {
		aborted = true				
		fmt.Println("##Error####################: SignRound 5, User i does not pass verifying Zero-Knowledge!",b);
		signature.setRoudFiveAborted(aborted);
	    }

	b = b+1
    }

    w := calculateW(userList)
    rx,ry := calculateR(userList)
    
    //3 calculate the signature (r,s)
    r := new(big.Int).Mod(rx,secp256k1.S256().N)
    mu := decrypt(privKey,w)
    mu.Mod(mu,secp256k1.S256().N)
    muInverse := new(big.Int).ModInverse(mu,secp256k1.S256().N)//need-test
    msgDigest,_ := new(big.Int).SetString(message,16)
    mMultiU := cipherMultiply((&privKey.PublicKey),u, msgDigest)
    rMultiV := cipherMultiply((&privKey.PublicKey),v, r)
    sEnc := cipherMultiply((&privKey.PublicKey),cipherAdd((&privKey.PublicKey),mMultiU, rMultiV), muInverse)

    s := decrypt(privKey,sEnc)
    s.Mod(s,secp256k1.S256().N)

    signature.setRoudFiveAborted(aborted)
    signature.setR(r)
    signature.setS(s)

    two,_ := new(big.Int).SetString("2",10)
    ryy := new(big.Int).Mod(ry,two)
    zero,_ := new(big.Int).SetString("0",10)
    cmp := ryy.Cmp(zero)
    recoveryParam := 1
    if cmp == 0 {
	recoveryParam = 0
    }

    tt := new(big.Int).Rsh(secp256k1.S256().N,1)
    comp := s.Cmp(tt)
    if comp > 0 {
	recoveryParam = 1
	s = new(big.Int).Sub(secp256k1.S256().N,s)
	signature.setS(s);
    }

    //need-test
    signature.setRecoveryParam(int32(recoveryParam))

    return signature

}

func Sign(userList *list.List,encX *big.Int,message string) *ECDSASignature {
	signRoudOne(userList,encX)
	signRoudTwo(userList,encX)

	roudThreeAborted := signRoundThree(userList,encX)
	if roudThreeAborted == true {
	    return nil
	}

	u := calculateU(userList)
	v := calculateV(userList)
	signRoundFour(userList,u)
	
	signature := signRoundFive(userList,u,v,message)
	if signature.getRoudFiveAborted() == true {
		return nil
	}

	return signature
}

func Verify(signature *ECDSASignature,message string,pkx *big.Int,pky *big.Int) bool {
    return signature.verify(message,pkx,pky)
}
