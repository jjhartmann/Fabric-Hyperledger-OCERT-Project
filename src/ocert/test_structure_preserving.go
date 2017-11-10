package ocert

import (
	"os"
	"fmt"
	"github.com/Nik-U/pbc"
)

/*
 * Run a single test
 */
func Stest() bool {
	fmt.Println("Stest")
	sharedParams := new(SharedParams)
	params := pbc.GenerateF(160)
	pairing := params.NewPairing()
	g1 := pairing.NewG1().Rand()
	g2 := pairing.NewG2().Rand()
	sharedParams.Params = params.String()
	sharedParams.G1 = g1.Bytes()
	sharedParams.G2 = g2.Bytes()

	VK, SK := SKeyGen(sharedParams)

	P := new(Pseudonym)
	P.C = pairing.NewG1().Rand().Bytes()
	P.D = pairing.NewG1().Rand().Bytes()
	PKc := new(ClientPublicKey)
	PKc.PK = pairing.NewG2().Rand().Bytes()
	ecert := SSign(sharedParams, SK, P, PKc)

	return SVerify(sharedParams, VK, P, PKc, ecert)
}

/*
 * Run test b times
 */
func RunSTest(b int) {
	for i := 0; i < b; i++ {
		if !Stest() {
			os.Exit(1)
		}
	} 
}