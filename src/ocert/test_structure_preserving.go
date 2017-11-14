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
	sharedParams := GenerateSharedParams()
	VK, SK := SKeyGen(sharedParams)

	P := new(Pseudonym)
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	P.C = pairing.NewG1().Rand().Bytes()
	P.D = pairing.NewG1().Rand().Bytes()
	PKc := new(ClientPublicKey)
	PKc.PK = pairing.NewG2().Rand().Bytes()
	ecert := SSign(sharedParams, SK, P, PKc)

	if SVerify(sharedParams, VK, P, PKc, ecert) {
		fmt.Println("Stest passes")
		return true
	} else {
		return false
	}
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