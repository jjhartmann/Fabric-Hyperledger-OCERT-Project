package ocert

import (
	"os"
	"fmt"
	"github.com/Nik-U/pbc"
)

/*
 * Run a single test
 */
func Ptest() bool {
	sharedParams := GenerateSharedParams()
	eqs := new(SystemOfEquations)
	vars := new(ProofVariables)
	pi := Setup(sharedParams, eqs, vars)

	consts := new(ProofConstants)
	return Prove(sharedParams, pi, consts)
}

func IotaTest(verbose bool) bool {
	fmt.Println("Testing Iota1 and Iota2 Conversion B")
	sharedParams := GenerateSharedParams()
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

	// Generate element to test on conversion
	Z := pairing.NewG1().Rand()
	b1pair := Iota1(pairing, Z)

	tb1 := pairing.NewG1().SetBytes(b1pair.b1)
	tb2 := pairing.NewG1().SetBytes(b1pair.b2)

	if (verbose) {
		fmt.Printf("Z = %s\n", Z)
		fmt.Printf("B1.b1 = %s\n", tb1)
		fmt.Printf("B1.b1 = %s\n", tb2)
	}

	// Alpha is a random integer of order prime (Secret Key)
	fmt.Println("Testing Rho1 and Rho2 Conversion to G")
	alpha := pairing.NewZr().Rand()
	Zprime := Rho1(pairing, b1pair, alpha)
	ret := Zprime.Equals(Z)

	if (verbose) {
		fmt.Printf("Zprime = %s\n", Zprime)
		fmt.Println("Test ==", ret)
	}

	return ret
}


/*
 * Run test b times
 */
func RunPTest(b int) {
	for i := 0; i < b; i++ {
		if !Ptest() {
			os.Exit(1)
		}
	} 
}