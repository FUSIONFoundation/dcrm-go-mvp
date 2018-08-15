package dcrm 

import (
	"math/big"
	"math/rand"
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/pbc"
)

type MTDCommitment struct {
		commitment *Commitment
		open *Open
}

func (mtdct *MTDCommitment) New(commitment *Commitment,open *Open) {
			mtdct.commitment = commitment
			mtdct.open = open
}

func multiLinnearCommit(rnd *rand.Rand,mpk *CmtMasterPublicKey,secrets []*big.Int) *MTDCommitment {
	e := mpk.pairing.NewZr()
	e.Rand()
	r := mpk.pairing.NewZr()
	r.Rand()

	h := func(target *pbc.Element,megs []string) {
		hash := sha256.New()
		for j := range megs {
		    hash.Write([]byte(megs[j]))
		}
		i := &big.Int{}
		target.SetBig(i.SetBytes(hash.Sum([]byte{})))
	}
    
    secretsBytes := make([]string,len(secrets))
    for i := range secrets {
	    count := ((secrets[i].BitLen()+7)/8)
	    se := make([]byte,count)
	    math.ReadBits(secrets[i], se[:])
	    secretsBytes[i] = string(se[:])
    }

    digest := mpk.pairing.NewZr()
    h(digest,secretsBytes[:])

    ge := mpk.pairing.NewG1()
    ge.MulZn(mpk.g,e)
    
    //he = mpk.h + ge
    he := mpk.pairing.NewG1()
    he.Add(mpk.h,ge)
    
    //he = r*he
    rhe := mpk.pairing.NewG1()
    rhe.MulZn(he,r)
    
    //dg = digest*mpk.g
    dg := mpk.pairing.NewG1()
    dg.MulZn(mpk.g,digest)
    
    //a = mpk.g + he
    a := mpk.pairing.NewG1()
    a.Add(dg,rhe)

    open := new(Open)
    open.New(r,secrets)
    commitment := new(Commitment)
    commitment.New(e,a)

    mtdct := new(MTDCommitment)
    mtdct.New(commitment,open)

    return mtdct
}

func checkcommitment(commitment *Commitment,open *Open,mpk *CmtMasterPublicKey) bool {
    g := mpk.g
    h := mpk.h
    
	f := func(target *pbc.Element,megs []string) {
		hash := sha256.New()
		for j := range megs {
		    hash.Write([]byte(megs[j]))
		}
		i := &big.Int{}
		target.SetBig(i.SetBytes(hash.Sum([]byte{})))
	}
    
    secrets := open.getSecrets();
    secretsBytes := make([]string,len(secrets))
    for i := range secrets {
	    count := ((secrets[i].BitLen()+7)/8)
	    se := make([]byte,count)
	    math.ReadBits(secrets[i], se[:])
	    secretsBytes[i] = string(se[:])
    }

    digest := mpk.pairing.NewZr()
    f(digest,secretsBytes[:])
    
    rg := mpk.pairing.NewG1()
    rg.MulZn(g,open.getRandomness())

    d1 := mpk.pairing.NewG1()
    d1.MulZn(g,commitment.pubkey)

    dh := mpk.pairing.NewG1()
    dh.Add(h,d1)

    gdn := mpk.pairing.NewG1()
    digest.Neg(digest)
    gdn.MulZn(g,digest)

    comd := mpk.pairing.NewG1()
    comd.Add(commitment.committment,gdn)
    b := pbc.DDH(rg,dh,comd,g,mpk.pairing)
    if b == false {
	fmt.Println("==========checkcommitment==== ")
}
    return b
}

/*func multiLinnearCommit(rnd *rand.Rand,mpk *CmtMasterPublicKey,secrets []*big.Int) *MTDCommitment {
    e := randomFromZn(mpk.q, rnd)
    r := randomFromZn(mpk.q, rnd)

    fmt.Println("=====multiLinnearCommit=====secrets[0] is: %v\n",secrets[0])
    fmt.Println("=====multiLinnearCommit=====secrets[1] is: %v\n",secrets[1])
    secretsBytes := make([]byte,(secrets[0].BitLen()+7)/8)
    math.ReadBits(secrets[0], secretsBytes[:])
    for j := range secrets {
	if j != 0 {
	    count := ((secrets[j].BitLen()+7)/8)
	    se := make([]byte,count)
	    math.ReadBits(secrets[j], se[:])
	    secretsBytes = append(secretsBytes[:],se[:]...) 
	}
    }

    //==============================
    //test := [...]byte{101, 200, 155, 190, 125, 249, 44, 116, 147, 5, 141, 58, 33, 140, 120, 56, 124, 209, 255, 122, 9, 214, 68, 13, 140, 141, 101, 167, 126, 26, 30, 90}
    //==============================
    test := sha256Hash(secretsBytes[:])
    digest := new(big.Int).SetBytes(test[:])
    digest.Mod(digest,mpk.q)
    fmt.Println("=====multiLinnearCommit====e is: %v\n",e)
    fmt.Println("=====multiLinnearCommit====r is: %v\n",r)
    fmt.Println("=====multiLinnearCommit====digest is: %v\n",digest)

    
    fmt.Println("=====multiLinnearCommit====start=====e is: %v\n",e)
    ge := mpk.pairing.NewG1()
    ge.MulBig(mpk.g,e)
    fmt.Println("=====multiLinnearCommit======end======\n")
    
    //he = mpk.h + ge
    he := mpk.pairing.NewG1()
    he.Add(mpk.h,ge)
    
    //he = r*he
    rhe := mpk.pairing.NewG1()
    rhe.MulBig(he,r)
    
    //dg = digest*mpk.g
    dg := mpk.pairing.NewG1()
    dg.MulBig(mpk.g,digest)
    
    //a = mpk.g + he
    a := mpk.pairing.NewG1()
    a.Add(dg,rhe)

    open := new(Open)
    open.New(r,secrets)
    commitment := new(Commitment)
    commitment.New(e,a)

    mtdct := new(MTDCommitment)
    mtdct.New(commitment,open)

    return mtdct
}

func checkcommitment(commitment *Commitment,open *Open,mpk *CmtMasterPublicKey) bool {
    g := mpk.g
    h := mpk.h
    secrets := open.getSecrets();
    fmt.Println("==========checkcommitment==========secrets 0 is: %v\n",secrets[0])
    fmt.Println("==========checkcommitment==========secrets 1 is: %v\n",secrets[1])
    fmt.Println("==========checkcommitment==========0 len is: %v\n",(secrets[0].BitLen()+7)/8)
    fmt.Println("==========checkcommitment==========1 len is: %v\n",(secrets[1].BitLen()+7)/8)
    secretsBytes := make([]byte,(secrets[0].BitLen()+7)/8)
    math.ReadBits(secrets[0], secretsBytes[:])
    for j := range secrets {
	if j != 0 {
	    count := ((secrets[j].BitLen()+7)/8)
	    se := make([]byte,count)
	    math.ReadBits(secrets[j], se[:])
	    secretsBytes = append(secretsBytes[:],se[:]...) 
	}
    }
    
    test := sha256Hash(secretsBytes[:])

    //==============================
    //test := [...]byte{101, 200, 155, 190, 125, 249, 44, 116, 147, 5, 141, 58, 33, 140, 120, 56, 124, 209, 255, 122, 9, 214, 68, 13, 140, 141, 101, 167, 126, 26, 30, 90}
    //==============================
    digest := new(big.Int).SetBytes(test[:])

    digest.Mod(digest,mpk.q)
    fmt.Println("==========checkcommitment====digest is %v======",digest)

    fmt.Println("==========checkcommitment====randomness is %v======",open.getRandomness())
    rg := mpk.pairing.NewG1()
    rg.MulBig(g,open.getRandomness())
    //tmp
//    rom,_:= new(big.Int).SetString("397481438335516493009041452860536291833732338080",10)
//    rg.MulBig(g,rom)
    //tmp

    fmt.Println("==========checkcommitment====commitment.pubkey is %v======",commitment.pubkey)
    d1 := mpk.pairing.NewG1()
    d1.MulBig(g,commitment.pubkey)
    //tmp
//    rom2,_:= new(big.Int).SetString("127651860473363855924259002507045823534834509003",10)
//    d1.MulBig(g,rom2)

    //tmp

    dh := mpk.pairing.NewG1()
    dh.Add(h,d1)

    gdn := mpk.pairing.NewG1()
    gdn.MulBig(g,new(big.Int).Neg(digest))

    comd := mpk.pairing.NewG1()
    comd.Add(commitment.committment,gdn)
    b := pbc.DDH(rg,dh,comd,g,mpk.pairing)
    if b == false {
	fmt.Println("==========checkcommitment==== ")
}
    return b
}
*/

func getBasePoint(pairing *pbc.Pairing) *pbc.Element {
    var p *pbc.Element
    cof := pairing.NewZr()
    num,_ := new(big.Int).SetString("10007920040268628970387373215664582404186858178692152430205359413268619141100079249246263148037326528074908",10)
    cof.SetBig(num)

    order,_ := new(big.Int).SetString("730750818665451459101842416358141509827966402561",10)
    q := pairing.NewZr()
    q.SetBig(order)

    for {
	    p = pairing.NewG1()
	    p.Rand()
	    ge := pairing.NewG1()
	    ge.MulZn(p,cof)

	    pq := pairing.NewG1()
	    pq.MulZn(ge,q)

	    if ge.Is0() || pq.Is0() {
		return ge
	    }
    }

    return nil 
}

func generateMasterPK() *CmtMasterPublicKey {
	pairing, err := pbc.NewPairingFromString("type a\nq 7313295762564678553220399414112155363840682896273128302543102778210584118101444624864132462285921835023839111762785054210425140241018649354445745491039387\nh 10007920040268628970387373215664582404186858178692152430205359413268619141100079249246263148037326528074908\nr 730750818665451459101842416358141509827966402561\nexp2 159\nexp1 17\nsign1 1\nsign0 1\n")
	if err != nil {
		fmt.Println("preload pairing fail.\n")
	}

	g := getBasePoint(pairing)
	q,_ := new(big.Int).SetString("730750818665451459101842416358141509827966402561",10)
	h := pbc.RandomPointInG1(pairing)
	cmpk := new(CmtMasterPublicKey)
	cmpk.New(g,q,h,pairing)
	return cmpk
}

func (this *MTDCommitment) cmtOpen() *Open {
    return this.open
}

func (this *MTDCommitment) cmtCommitment() *Commitment {
    return this.commitment
}
