package dcrm 

import (
    "math/big"
    "math/rand"
    "fmt"
    "github.com/ethereum/go-ethereum/crypto/secp256k1"
    "github.com/ethereum/go-ethereum/common/math"
)

var (
    ech = make (chan bool,1) //1是必须的
    u1ch = make (chan bool,1) //1是必须的
    u2ch = make (chan bool,1) //1是必须的
    u3ch = make (chan bool,1) //1是必须的
)

type ZkpKG struct {
    z *big.Int
    u1_x *big.Int
    u1_y *big.Int
    u2 *big.Int
    u3 *big.Int
    e *big.Int
    s1 *big.Int
    s2 *big.Int
    s3 *big.Int
}

func (this *ZkpKG) New(params *PublicParameters,eta *big.Int,rand *rand.Rand,cx *big.Int,cy *big.Int,w *big.Int,r *big.Int) {
    N := params.paillierPubKey.N
    q := secp256k1.S256().N
    nSquared := new(big.Int).Mul(N,N)
    nTilde := params.nTilde
    h1 := params.h1
    h2 := params.h2
    one,_ := new(big.Int).SetString("1",10)
    g := new(big.Int).Add(N,one)
    q2 := new(big.Int).Mul(q,q)
    q3 := new(big.Int).Mul(q2,q)
    alpha := randomFromZn(q3,rand)
    qn := new(big.Int).Mul(q,nTilde)
    rho := randomFromZn(qn,rand)
    beta := randomFromZnStar(N,rand)
    qt := new(big.Int).Mul(q3,nTilde)
    gamma := randomFromZn(qt,rand)
    zz := modPow(h1,eta,nTilde)
    zzz := modPow(h2,rho,nTilde)
    zmz := new(big.Int).Mul(zz,zzz)
    this.z = new(big.Int).Mod(zmz,nTilde)

    if alpha.Sign() == -1 {
		alpha.Add(alpha,secp256k1.S256().P)
	}
    upk := make([]byte, 32)
    math.ReadBits(alpha,upk[:])
    this.u1_x,this.u1_y = secp256k1.S256().ScalarMult(cx,cy,upk[:])

    uu := modPow(g,alpha,nSquared)
    uuu := modPow(beta,N,nSquared)
    utmp := new(big.Int).Mul(uu,uuu)
    this.u2 = new(big.Int).Mod(utmp,nSquared)

    tt := modPow(h1,alpha,nTilde)
    ttt := modPow(h2,gamma,nTilde)
    ttmp := new(big.Int).Mul(tt,ttt)
    this.u3 = new(big.Int).Mod(ttmp,nTilde)

    byte1 := getBytes(cx,cy)

    wlen := (w.BitLen() + 7)/8
    ws := make([]byte,wlen)
    math.ReadBits(w,ws[:])

    zlen := (this.z.BitLen() + 7)/8
    zs := make([]byte,zlen)
    math.ReadBits(this.z,zs[:])

    u1s := getBytes(this.u1_x,this.u1_y)

    u2len := (this.u2.BitLen() + 7)/8
    u2s := make([]byte,u2len)
    math.ReadBits(this.u2,u2s[:])

    u3len := (this.u3.BitLen() + 7)/8
    u3s := make([]byte,u3len)
    math.ReadBits(this.u3,u3s[:])

    ss := make([]string,6)
    ss[0] = string(byte1[:])
    ss[1] = string(ws[:])
    ss[2] = string(zs[:])
    ss[3] = string(u1s[:])
    ss[4] = string(u2s[:])
    ss[5] = string(u3s[:])

    //need-test
    digest := sha256Hash(ss[:])
    if len(digest) == 0 {
	return
    }

    this.e = new(big.Int).SetBytes(digest[:])

    er := new(big.Int).Mul(this.e,eta)
    this.s1 = new(big.Int).Add(er,alpha)

    rn := modPow(r,this.e,N)
    rb := new(big.Int).Mul(rn,beta)
    this.s2 = new(big.Int).Mod(rb,N)

    eho := new(big.Int).Mul(this.e,rho)
    this.s3 = new(big.Int).Add(eho,gamma)
}

func (this *ZkpKG) f1(cx *big.Int,cy *big.Int,rx *big.Int,ry *big.Int) bool {
    /*//(apha + e*xShare)*G
    sk := make([]byte, 32)
    math.ReadBits(this.s1,sk[:])
    ccx,ccy := secp256k1.S256().ScalarMult(cx,cy,sk[:])

    //e*(xShare*G)
    ek := make([]byte, 32)
    math.ReadBits(this.e,ek[:])
    rnex,rney := secp256k1.S256().ScalarMult(rx,ry,ek[:])

    //this.u1_x,this.u1_y := alpha*G
    ccrx,ccry := secp256k1.S256().Add(this.u1_x,this.u1_y,rnex,rney)*/

    u1ch <-true
    return true
}

func (this *ZkpKG) f2(g *big.Int,nSquared *big.Int,N *big.Int,w *big.Int) bool {
    gs := modPow(g,this.s1, nSquared)
    s2n := modPow(this.s2,N, nSquared)
    gss := new(big.Int).Mul(gs,s2n)
    en := new(big.Int).Neg(this.e)
    wen := modPow(w,en,nSquared)
    gsw := new(big.Int).Mul(gss,wen)
    gswns := new(big.Int).Mod(gsw,nSquared)
    ng := this.u2.Cmp(gswns)
    if ng == 0 {
	u2ch <-true
	return true
    }

    u2ch <-false
    return false
}

func (this *ZkpKG) f3(h1 *big.Int,nTilde *big.Int,h2 *big.Int) bool {
    hs := modPow(h1,this.s1,nTilde)
    hsn := modPow(h2,this.s3,nTilde)
    hh := new(big.Int).Mul(hs,hsn)
    en := new(big.Int).Neg(this.e)
    zen := modPow(this.z,en,nTilde)
    hz := new(big.Int).Mul(hh,zen)
    hn := new(big.Int).Mod(hz,nTilde)
    uh := this.u3.Cmp(hn)
    if uh == 0 {
	u3ch <-true
	return true
    }
    u3ch <-false
    return false
}

func (this *ZkpKG) f4(cx *big.Int,cy *big.Int,w *big.Int) bool {
    byte1 := getBytes(cx,cy)

    wlen := (w.BitLen() + 7)/8
    ws := make([]byte,wlen)
    math.ReadBits(w,ws[:])

    zlen := (this.z.BitLen() + 7)/8
    zs := make([]byte,zlen)
    math.ReadBits(this.z,zs[:])

    u1s := getBytes(this.u1_x,this.u1_y)

    u2len := (this.u2.BitLen() + 7)/8
    u2s := make([]byte,u2len)
    math.ReadBits(this.u2,u2s[:])

    u3len := (this.u3.BitLen() + 7)/8
    u3s := make([]byte,u3len)
    math.ReadBits(this.u3,u3s[:])

    ss := make([]string,6)
    ss[0] = string(byte1[:])
    ss[1] = string(ws[:])
    ss[2] = string(zs[:])
    ss[3] = string(u1s[:])
    ss[4] = string(u2s[:])
    ss[5] = string(u3s[:])

    digest := sha256Hash(ss[:])
    if len(digest) == 0 {
	ech <-false
	return false
    }

    eRecovered := new(big.Int).SetBytes(digest[:])
    ee := eRecovered.Cmp(this.e)
    if ee == 0 {
	ech <-true
	return true
    }
    
    ech <-false
    return false
}

func (this *ZkpKG) verify(params *PublicParameters,BitCurve *secp256k1.BitCurve,rx *big.Int,ry *big.Int,w *big.Int) bool {
    cx := secp256k1.S256().Gx
    cy := secp256k1.S256().Gy

    h1 := params.h1
    h2 := params.h2
    N := params.paillierPubKey.N
    nTilde := params.nTilde
    nSquared := new(big.Int).Mul(N,N)
    one,_ := new(big.Int).SetString("1",10)
    g := new(big.Int).Add(N,one)

    go this.f1(cx,cy,rx,ry)
    go this.f2(g,nSquared,N,w)
    go this.f3(h1,nTilde,h2)
    go this.f4(cx,cy,w)

    count := 0
    //need-test
    for {
	select { //没有值则select继续循环 否则执行一次
	    case v1 := <- u1ch: //select case 只限定bool值
	    	count += 1
		if v1 == false {
		fmt.Println("======zkpkg v1===========\n")
		return false
		}
	    case v2 := <- u2ch: //select case 只限定bool值
	    	count += 1
		if v2 == false {
		fmt.Println("======zkpkg v2===========\n")
		return false
		}
	    case v3 := <- u3ch: //select case 只限定bool值
	    	count += 1
		if v3 == false {
		fmt.Println("======zkpkg v3===========\n")
		return false
		}
	    case v4 := <- ech: //select case 只限定bool值
	    	count += 1
		if v4 == false {
		fmt.Println("======zkpkg v4===========\n")
		return false
		}

	}
	
	if count == 4 {
	    break
	}
    }

    return true
}

