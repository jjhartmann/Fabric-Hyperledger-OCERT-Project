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
	pi := PSetup(sharedParams, eqs, vars)

	consts := new(ProofConstants)
	return PProve(sharedParams, pi, consts)
}

// Test mapping between G and B
func IotaRhoTest(verbose bool) bool {
	fmt.Println("Testing Iota1 and Iota2 Conversion B")
	sharedParams := GenerateSharedParams()
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

	// Generate element to test on conversion
	Z1 := pairing.NewG1().Rand()
	b1pair := Iota1(pairing, Z1)

	Z2 := pairing.NewG2().Rand()
	b2pair := Iota2(pairing, Z2)

	tb1 := pairing.NewG1().SetBytes(b1pair.b1)
	tb2 := pairing.NewG1().SetBytes(b1pair.b2)
  tb3 := pairing.NewG2().SetBytes(b2pair.b1)
  tb4 := pairing.NewG2().SetBytes(b2pair.b2)


	if (verbose) {
	  fmt.Println("Iota1")
    fmt.Printf("Z1 = %s\n", Z1)
		fmt.Printf("B1.b1 = %s\n", tb1)
		fmt.Printf("B1.b1 = %s\n", tb2)

    fmt.Println("Iota2")
    fmt.Printf("Z2 = %s\n", Z2)
    fmt.Printf("B2.b1 = %s\n", tb3)
    fmt.Printf("B2.b1 = %s\n", tb4)
	}

	// Alpha is a random integer of order prime (Secret Key)
	fmt.Println("Testing Rho1 and Rho2 Conversion to G")
	alpha := pairing.NewZr().Rand()
	Zprime1 := Rho1(pairing, b1pair, alpha)
	ret1 := Zprime1.Equals(Z1)
  Zprime2 := Rho2(pairing, b2pair, alpha)
  ret2 := Zprime2.Equals(Z2)

	if (verbose) {
    fmt.Printf("Zprime1 = %s\n", Zprime1)
    fmt.Printf("Zprime2 = %s\n", Zprime2)
    fmt.Println("Test1 ==", ret1)
    fmt.Println("Test2 ==", ret2)
	}

	return ret1 && ret2
}


// Test mapping between Zp and B
func TestIotaRhoPrime(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  fmt.Println("Creating CRS Sigma")
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  // Test IotaPrim: Zp -> B1
  fmt.Println("Calling IotaPrime")
  z := pairing.NewZr().Rand()
  B1 := IotaPrime1(pairing, z, sigma)
  B2 := IotaPrime2(pairing, z, sigma)

  if (verbose){
    b1 := pairing.NewG1().SetBytes(B1.b1)
    b2 := pairing.NewG1().SetBytes(B1.b2)
    b3 := pairing.NewG2().SetBytes(B2.b1)
    b4 := pairing.NewG2().SetBytes(B2.b1)

    fmt.Printf("z = %s\n", z)
    fmt.Printf("b1 = %s\n", b1)
    fmt.Printf("b2 = %s\n", b2)
    fmt.Printf("b3 = %s\n", b3)
    fmt.Printf("b4 = %s\n", b4)

  }

  fmt.Println("Calling RhoPrime")
  zP1 := RhoPrime1(pairing, B1, alpha)
  zP2 := RhoPrime2(pairing, B2, alpha)


  // To Check, we need to multiple the generator g1 by z to see
  // if the conversin back is successful.
  fmt.Println("Testing Equality:")
  P1 := pairing.NewG1().SetBytes(sigma.U[0].u1)
  P2 := pairing.NewG2().SetBytes(sigma.V[0].u1)

  retU := zP1.Equals(pairing.NewG1().MulZn(P1, z))
  retV := zP2.Equals(pairing.NewG2().MulZn(P2, z))

  if (verbose){
    fmt.Println("retU =", retU)
    fmt.Println("retV =", retV)
  }

  return retU && retV
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