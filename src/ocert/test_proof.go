package ocert

import (
	"os"
	"fmt"
	"github.com/Nik-U/pbc"
  "reflect"
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
 Test FMap
 F: B1^2 * B2^2 -> BT^4
 */
func TestFMap(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  fmt.Println("Creating CRS Sigma")
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  fmt.Println("Creating Elements in B1 & B2")
  z := pairing.NewZr().Rand() // testing element to map
  Y := pairing.NewG2().Rand() // test element in G2
  B1 := IotaPrime1(pairing, z, sigma)
  B2 := Iota2(pairing, Y)

  fmt.Println("Mapping into BT")
  BT := FMap(pairing, B1, B2)

  // Manual create pairs
  X1 := pairing.NewG1().SetBytes(B1.b1)
  X2 := pairing.NewG1().SetBytes(B1.b2)
  Y1 := pairing.NewG2().SetBytes(B2.b1)
  Y2 := pairing.NewG2().SetBytes(B2.b2)

  ret1 := pairing.NewGT().Pair(X1, Y1).Equals(pairing.NewGT().SetBytes(BT.el11))
  ret2 := pairing.NewGT().Pair(X1, Y2).Equals(pairing.NewGT().SetBytes(BT.el12))
  ret3 := pairing.NewGT().Pair(X2, Y1).Equals(pairing.NewGT().SetBytes(BT.el21))
  ret4 := pairing.NewGT().Pair(X2, Y2).Equals(pairing.NewGT().SetBytes(BT.el22))

  if (verbose){
    fmt.Println("BT[1, 1] == e(X1, Y1): ", ret1)
    fmt.Println("BT[1, 2] == e(X1, Y2): ", ret2)
    fmt.Println("BT[2, 1] == e(X2, Y1): ", ret3)
    fmt.Println("BT[2, 2] == e(X2, Y2): ", ret4)
  }

  return ret1 && ret2 && ret3 && ret4
}


func TestIotaHat(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  fmt.Println("Creating CRS Sigma")
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  fmt.Println("Create IotaHat int BT")
  AT := pairing.NewG2().Rand() // Element to test
  BT := IotaHat(pairing, AT, sigma)

  fmt.Println("Testing: F(ι'1(1), ι2(AT)) = F(u, (O,AT))")
  c := pairing.NewZr().SetInt32(1)
  B1_1 := IotaPrime1(pairing, c, sigma)
  B2_1 := Iota2(pairing, AT)
  BT_1 := FMap(pairing, B1_1 , B2_1)

  B1_2 := new(BPair)
  B1_2.b1 = sigma.u.u1
  B1_2.b2 = sigma.u.u2
  B2_2 := new(BPair)
  B2_2.b1 = pairing.NewG2().Set0().Bytes()
  B2_2.b2 = AT.Bytes()
  BT_2 := FMap(pairing, B1_2, B2_2)


  ret1 := reflect.DeepEqual(BT, BT_1)
  ret2 := reflect.DeepEqual(BT_1, BT_2)


  if (verbose){
    fmt.Println("BT:   ", BT)
    fmt.Println("BT_1: ", BT_1)
    fmt.Println("BT_2: ", BT_2)
    fmt.Println("ι^T(AT) == F(ι'1(1), ι2(AT)): ", ret1)
    fmt.Println("F(ι'1(1), ι2(AT)) = F(u, (O,AT)): ", ret2)
  }

  return ret1 && ret2
}



func TestCompleteMatrixMapping(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  fmt.Println("Creating CRS Sigma")
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  fmt.Println("Testing Multi-Scalar Multiplication Mapping Matrix for G2")
  fmt.Println("    - F(ι'1(x), ι2(Y)) = F(ι'1(1), ι2(xY)) = ιT(xY)")
  c := pairing.NewZr().SetInt32(1)
  x := pairing.NewZr().Rand()
  Y := pairing.NewG2().Rand()
  xY := MultiScalar_f_G2_map(pairing, x, Y)

  B1x := IotaPrime1(pairing, x, sigma)
  B1c := IotaPrime1(pairing, c, sigma)
  B2Y := Iota2(pairing, Y)
  B2xY := Iota2(pairing, xY)


  BTxY := IotaHat(pairing, xY, sigma)
  BTF1 := FMap(pairing, B1x, B2Y)
  BTF2 := FMap(pairing, B1c, B2xY)

  ret1 := reflect.DeepEqual(BTxY, BTF1)
  ret2 := reflect.DeepEqual(BTF1, BTF2)

  if (verbose) {
    fmt.Println("BTxY:   ", BTxY)
    fmt.Println("BTF1:   ", BTF1)
    fmt.Println("BTF2:   ", BTF2)
    fmt.Println("F(ι'1(1), ι2(xY)) = ιT(xY):           ", ret1)
    fmt.Println("F(ι'1(x), ι2(Y)) = F(ι'1(1), ι2(xY)): ", ret2)
  }



  fmt.Println()
  fmt.Println("Testing Multi-Scalar Multiplication Mapping Matrix for G1")
  fmt.Println("    - F(ι'1(x), ι2(Y)) = F(ι'1(1), ι2(xY)) = ιT(xY)")
  MSc := pairing.NewZr().SetInt32(1)
  MSx := pairing.NewZr().Rand()
  MSY := pairing.NewG1().Rand()
  MSxY := MultiScalar_f_G1_map(pairing, MSx, MSY)

  MSB1x := IotaPrime2(pairing, MSx, sigma)
  MSB1c := IotaPrime2(pairing, MSc, sigma)
  MSB2Y := Iota1(pairing, MSY)
  MSB2xY := Iota1(pairing, MSxY)


  MSBTxY := IotaHat(pairing, MSxY, sigma)
  MSBTF1 := FMap(pairing, MSB1x, MSB2Y)
  MSBTF2 := FMap(pairing, MSB1c, MSB2xY)

  ret4 := reflect.DeepEqual(MSBTxY, MSBTF1)
  ret5 := reflect.DeepEqual(MSBTF1, MSBTF2)

  if (verbose) {
    fmt.Println("BTxY:   ", MSBTxY)
    fmt.Println("BTF1:   ", MSBTF1)
    fmt.Println("BTF2:   ", MSBTF2)
    fmt.Println("F(ι'1(1), ι2(xY)) = ιT(xY):           ", ret4)
    fmt.Println("F(ι'1(x), ι2(Y)) = F(ι'1(1), ι2(xY)): ", ret5)
  }



  fmt.Println()
  fmt.Println("Testing Pairing Product Mapping Matrix")
  PPX := pairing.NewG1().Rand()
  PPY := pairing.NewG2().Rand()
  PPgt := ProductPairing_e_GT_map(pairing, PPX, PPY)

  PPB1 := Iota1(pairing, PPX)
  PPB2 := Iota2(pairing, PPY)
  PPBTi := IotaT(pairing, PPgt)
  PPBTF := FMap(pairing, PPB1, PPB2)

  ret3 := reflect.DeepEqual(PPBTi, PPBTF)

  if (verbose){
    fmt.Println("PPBTi:   ", PPBTi)
    fmt.Println("PPBTF:   ", PPBTF)
    fmt.Println("F(ι1(X), ι2(Y)) = ιT(e(X,Y)): ", ret3)
  }

  return ret2 && ret1 && ret3
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