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
func Ptest(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if verbose {fmt.Println("Creating Equations")}

  // Shared generators
  G  := pairing.NewG1().SetBytes(sharedParams.G1)
  H  := pairing.NewG2().SetBytes(sharedParams.G2)

  // Verifcation key
  VK, SK := SKeyGen(sharedParams)

  // Auditor Key pair
  PKa, SKa := EKeyGen(sharedParams)
  _ = SKa

  // Pseudonym Generation
  clientID := new(ClientID)
  clientID.ID = pairing.NewG1().Rand().Bytes()
  P := EEnc(sharedParams, PKa, clientID)

  // Rerandomization
  Pprime, rprime := ERerand(sharedParams, PKa, P)

  // Client Keypair
  PKc := new(ClientPublicKey)
  Xc := pairing.NewZr().Rand().Bytes()
  PKc.PK = pairing.NewG2().MulZn(H, pairing.NewZr().SetBytes(Xc)).Bytes()

  // Ecert Generation
  ecert := SSign(sharedParams, SK, P, PKc)

  // Construct Var
  eqs := new(SystemOfEquations)
  vars := new(ProofVariables)
  vars.PKa = PKa
  vars.P = P
  vars.VK = VK
  vars.RPrime = rprime
  vars.PKc = PKc
  vars.Xc = Xc
  vars.E = ecert

  pi := PSetup(sharedParams, eqs, vars)

  if verbose {fmt.Println("Testing Structure Integrity")}
  retValEq1 := len(pi.Eq1.Theta) == 2 && len(pi.Eq1.Pi) == 1 &&
    len(pi.Eq1.d) == 1 && len(pi.Eq1.cprime) == 1 &&
    len(pi.Eq1.c) == 0 && len(pi.Eq1.dprime) == 0
  if verbose {fmt.Println("Eq1 Structure:\t", retValEq1)}

  retValEq2 := len(pi.Eq2.Theta) == 1 && len(pi.Eq2.Pi) == 2 &&
    len(pi.Eq2.d) == 0 && len(pi.Eq2.cprime) == 0 &&
    len(pi.Eq2.c) == 1 && len(pi.Eq2.dprime) == 1
  if verbose {fmt.Println("Eq2 Structure:\t", retValEq2)}

  retValEq3 := len(pi.Eq3.Theta) == 1 && len(pi.Eq3.Pi) == 2 &&
    len(pi.Eq3.d) == 0 && len(pi.Eq3.cprime) == 0 &&
    len(pi.Eq3.c) == 1 && len(pi.Eq3.dprime) == 1
  if verbose {fmt.Println("Eq3 Structure:\t", retValEq3)}

  retValEq4 := len(pi.Eq4.Theta) == 2 && len(pi.Eq4.Pi) == 2 &&
    len(pi.Eq4.d) == 0 && len(pi.Eq4.cprime) == 0 &&
    len(pi.Eq4.c) == 4 && len(pi.Eq4.dprime) == 0
  if verbose {fmt.Println("Eq4 Structure:\t", retValEq4)}

  retValEq5 := len(pi.Eq5.Theta) == 2 && len(pi.Eq5.Pi) == 2 &&
    len(pi.Eq5.d) == 2 && len(pi.Eq5.cprime) == 0 &&
    len(pi.Eq5.c) == 1 && len(pi.Eq5.dprime) == 0
  if verbose {fmt.Println("Eq5 Structure:\t", retValEq5)}

  // Create constants for verify
  consts := new(ProofConstants)
  consts.VK = VK
  consts.PKa = PKa
  consts.Egh = pairing.NewGT().Pair(G, H).Bytes()
  consts.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VK.Z)).Bytes()
  consts.PPrime = Pprime

  Cprime := pairing.NewG1().MulZn(G, pairing.NewZr().SetBytes(rprime))
  Cprime = pairing.NewG1().Add(Cprime, pairing.NewG1().SetBytes(P.C))
  if verbose {
    fmt.Println("Cprime Equality:", reflect.DeepEqual(Cprime.Bytes(), Pprime.C))
    fmt.Println("ERerand C':        ", Pprime.C)
    fmt.Println("Reconstruced  C':  ", Cprime.Bytes())
  }

  Dprime := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(PKa.PK),
    pairing.NewZr().SetBytes(rprime))
  Dprime = pairing.NewG1().Add(Dprime, pairing.NewG1().SetBytes(P.D))
  if verbose {
    fmt.Println("Dprime Equality:", reflect.DeepEqual(Cprime.Bytes(), Pprime.C))
    fmt.Println("ERerand D':        ", Pprime.D)
    fmt.Println("Reconstruced  D':  ", Dprime.Bytes())
  }

  return retValEq1 && retValEq2 && retValEq3 &&
         retValEq4 && retValEq5 && PProve(sharedParams, pi, consts)
}

// Test mapping between G and B
func TestIotaRho(verbose bool) bool {

  if (verbose) {fmt.Println("Testing Iota1 and Iota2 Conversion B")}
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
  if (verbose) {fmt.Println("Testing Rho1 and Rho2 Conversion to G")}
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

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  // Test IotaPrim: Zp -> B1
  if (verbose) {fmt.Println("Calling IotaPrime")}
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

  if (verbose) {fmt.Println("Calling RhoPrime")}
  zP1 := RhoPrime1(pairing, B1, alpha)
  zP2 := RhoPrime2(pairing, B2, alpha)


  // To Check, we need to multiple the generator g1 by z to see
  // if the conversin back is successful.
  if (verbose) {fmt.Println("Testing Equality:")}
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

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  if (verbose) {fmt.Println("Creating Elements in B1 & B2")}
  z := pairing.NewZr().Rand() // testing element to map
  Y := pairing.NewG2().Rand() // test element in G2
  B1 := IotaPrime1(pairing, z, sigma)
  B2 := Iota2(pairing, Y)

  if (verbose) {fmt.Println("Mapping into BT")}
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

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {fmt.Println("Create IotaHat int BT")}
  AT := pairing.NewG2().Rand() // Element to test
  BT := IotaHat(pairing, AT, sigma)

  if (verbose) {fmt.Println("Testing: F(ι'1(1), ι2(AT)) = F(u, (O,AT))")}
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

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {
    fmt.Println("Testing Multi-Scalar Multiplication Mapping Matrix for G2")
    fmt.Println("    - F(ι'1(x), ι2(Y)) = F(ι'1(1), ι2(xY)) = ιT(xY)")
  }
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



  if (verbose) {
    fmt.Println()
    fmt.Println("Testing Multi-Scalar Multiplication Mapping Matrix for G1")
    fmt.Println("    - F(ι'1(x), ι2(Y)) = F(ι'1(1), ι2(xY)) = ιT(xY)")
  }
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



  if (verbose) {
    fmt.Println()
    fmt.Println("Testing Pairing Product Mapping Matrix")
  }
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

  return ret2 && ret1 && ret3 && ret4 && ret5
}

func TestSimpleCommitment(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  // Create random value in G1
  X := pairing.NewG1().Rand()
  Xi := Iota1(pairing, X)

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  // Test with Functions
  chi := []*pbc.Element{X}
  C, Ru, rmat := CreateCommitmentOnG1(pairing, chi, sigma)
  _ = C
  _ = Ru

  if (verbose){
    fmt.Println("Function Commitsments C")
    fmt.Println(pairing.NewG1().SetBytes(C[0].b1))
    fmt.Println(pairing.NewG1().SetBytes(C[0].b2))
  }


  if (verbose){
    fmt.Println("Function Commitsments Ru")
    fmt.Println(pairing.NewG1().SetBytes(Ru[0].b1))
    fmt.Println(pairing.NewG1().SetBytes(Ru[0].b2))
  }

  // Simple Toy Test
  r1 := rmat.mat[0][0]
  r2 := rmat.mat[0][1]

  b1 := pairing.NewG1().SetBytes(Xi.b1)
  b2 := pairing.NewG1().SetBytes(Xi.b2)

  if (verbose){
    fmt.Println("B1 and B2")
    fmt.Println(b1)
    fmt.Println(b2)
  }

  ruP1 := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(sigma.U[0].u1), r1)
  ruQ1 := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(sigma.U[0].u2), r1)
  ruP2 := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(sigma.U[1].u1), r2)
  ruQ2 := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(sigma.U[1].u2), r2)

  // Add components together
  P := pairing.NewG1().Add(ruP1, ruP2)
  Q := pairing.NewG1().Add(ruQ1, ruQ2)

  if (verbose){
    fmt.Println("P and Q")
    fmt.Println(P)
    fmt.Println(Q)
  }

  c1 := pairing.NewG1().Add(b1, P)
  c2 := pairing.NewG1().Add(b2, Q)

  if (verbose){
    fmt.Println("C1 and C2")
    fmt.Println(c1)
    fmt.Println(c2)
  }

  // Subtract
  cp1 := pairing.NewG1().Sub(c1, P)
  cp2 := pairing.NewG1().Sub(c2, Q)

  if (verbose){
    fmt.Println("C'1 and C'2")
    fmt.Println(cp1)
    fmt.Println(cp2)
  }

  if (verbose) {
    fmt.Println("Testing Equality")
  }

  ret1 := cp1.Equals(b1)
  ret2 := cp2.Equals(b2)

  ret3 := c1.Equals(pairing.NewG1().SetBytes(C[0].b1))
  ret4 := c2.Equals(pairing.NewG1().SetBytes(C[0].b2))

  return ret1 && ret2 && ret3 && ret4
}

func TestCreateCommitmentsG1(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {fmt.Println("Create Commitments On G1")}
  chi := []*pbc.Element{
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
    pairing.NewG1().Rand(),
  }
  C, Ru, _ := CreateCommitmentOnG1(pairing, chi, sigma)
  ret1 := (len(chi) == len(C) && len(C) == len(Ru))

  if (verbose){
    fmt.Println("Length Consistency Test: ", ret1)
    for i:=0; i<len(C); i++ {
      fmt.Printf("%s\t", pairing.NewG1().SetBytes(C[i].b1))
      fmt.Printf("%s\n", pairing.NewG1().SetBytes(C[i].b2))
    }
    fmt.Println()
    fmt.Println()
  }


  if (verbose){fmt.Println("Testing Equality: ι1(X) = C - Ru")}
  ret2 := true
  for i:=0; i<len(C); i++ {
    //Bc := C[i] // Pair b1 and b2 in B1
    Bp1 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(C[i].b1),
      pairing.NewG1().SetBytes(Ru[i].b1))
    Bp2 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(C[i].b2),
      pairing.NewG1().SetBytes(Ru[i].b2))

    //Cp := C[i].SubinG1(pairing, C[i], Ru[i])
    //Bp1 := pairing.NewG1().SetBytes(Cp.b1)
    //Bp2 := pairing.NewG1().SetBytes(Cp.b2)

    Bi := Iota1(pairing, chi[i])
    Bi1 := pairing.NewG1().SetBytes(Bi.b1)
    Bi2 := pairing.NewG1().SetBytes(Bi.b2)

    tmp1 := Bp1.Equals(Bi1)
    tmp2 := Bp2.Equals(Bi2)

    if (verbose){
      fmt.Println("Testing Equality: ", i, tmp1 && tmp2)
      fmt.Printf("%s\t",Bp1)
      fmt.Printf("%s\n",Bp2)
      fmt.Printf("%s\t",Bi1)
      fmt.Printf("%s\n",Bi2)
    }
    ret2 = (ret2 && tmp1 && tmp2)
  }

  return ret1 && ret2
}

func TestCreateCommitmentPrimeOnG1(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {fmt.Println("Create Commitments On G1")}
  x := []*pbc.Element{
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
  }
  C, Ru, _ := CreateCommitmentPrimeOnG1(pairing, x, sigma)
  ret1 := (len(x) == len(C) && len(C) == len(Ru))

  if (verbose){
    fmt.Println("Length Consistency Test: ", ret1)
    for i:=0; i<len(C); i++ {
      fmt.Printf("%s\t", pairing.NewG1().SetBytes(C[i].b1))
      fmt.Printf("%s\n", pairing.NewG1().SetBytes(C[i].b2))
    }
    fmt.Println()
    fmt.Println()
  }


  if (verbose){fmt.Println("Testing Equality: ι'1(X) = C - Ru")}
  ret2 := true
  for i:=0; i<len(C); i++ {
    //Bc := C[i] // Pair b1 and b2 in B1
    //Br := Ru[i]
    Bp1 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(C[i].b1),
      pairing.NewG1().SetBytes(Ru[i].b1))
    Bp2 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(C[i].b2),
      pairing.NewG1().SetBytes(Ru[i].b2))

    Bi := IotaPrime1(pairing, x[i], sigma)
    Bi1 := pairing.NewG1().SetBytes(Bi.b1)
    Bi2 := pairing.NewG1().SetBytes(Bi.b2)

    tmp1 := Bp1.Equals(Bi1)
    tmp2 := Bp2.Equals(Bi2)

    if (verbose){
      fmt.Println("Testing Equality: ", i, tmp1 && tmp2)
      fmt.Printf("%s\t",Bp1)
      fmt.Printf("%s\n",Bp2)
      fmt.Printf("%s\t",Bi1)
      fmt.Printf("%s\n",Bi2)
    }
    ret2 = (ret2 && tmp1 && tmp2)
  }

  return ret1 && ret2
}

func TestCreateCommitmentsG2(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {fmt.Println("Create Commitments On G2")}
  Y := []*pbc.Element{
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
    pairing.NewG2().Rand(),
  }
  C, Su, _ := CreateCommitmentOnG2(pairing, Y, sigma)
  ret1 := (len(Y) == len(C) && len(C) == len(Su))

  if (verbose){
    fmt.Println("Length Consistency Test: ", ret1)
    for i:=0; i<len(C); i++ {
      fmt.Printf("%s\t", pairing.NewG2().SetBytes(C[i].b1))
      fmt.Printf("%s\n", pairing.NewG2().SetBytes(C[i].b2))
    }
    fmt.Println()
    fmt.Println()
  }


  if (verbose){fmt.Println("Testing Equality: ι2(X) = C - Ru")}
  ret2 := true
  for i:=0; i<len(C); i++ {
    //Bc := C[i] // Pair b1 and b2 in B1
    //Br := Ru[i]
    Bp1 := pairing.NewG2().Sub(pairing.NewG2().SetBytes(C[i].b1),
      pairing.NewG2().SetBytes(Su[i].b1))
    Bp2 := pairing.NewG2().Sub(pairing.NewG2().SetBytes(C[i].b2),
      pairing.NewG2().SetBytes(Su[i].b2))

    Bi := Iota2(pairing, Y[i])
    Bi1 := pairing.NewG2().SetBytes(Bi.b1)
    Bi2 := pairing.NewG2().SetBytes(Bi.b2)

    tmp1 := Bp1.Equals(Bi1)
    tmp2 := Bp2.Equals(Bi2)

    if (verbose){
      fmt.Println("Testing Equality: ", i, tmp1 && tmp2)
      fmt.Printf("%s\t",Bp1)
      fmt.Printf("%s\n",Bp2)
      fmt.Printf("%s\t",Bi1)
      fmt.Printf("%s\n",Bi2)
    }
    ret2 = (ret2 && tmp1 && tmp2)
  }

  return ret1 && ret2
}


func TestCreateCommitmentPrimeOnG2(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if (verbose) {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {fmt.Println("Create Commitment Primes On G2")}
  y := []*pbc.Element{
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
    pairing.NewZr().Rand(),
  }
  C, Su, _ := CreateCommitmentPrimeOnG2(pairing, y, sigma)
  ret1 := (len(y) == len(C) && len(C) == len(Su))

  if (verbose){
    fmt.Println("Length Consistency Test: ", ret1)
    for i:=0; i<len(C); i++ {
      fmt.Printf("%s\t", pairing.NewG2().SetBytes(C[i].b1))
      fmt.Printf("%s\n", pairing.NewG2().SetBytes(C[i].b2))
    }
    fmt.Println()
    fmt.Println()
  }


  if (verbose){fmt.Println("Testing Equality: ι'2(X) = C - Ru")}
  ret2 := true
  for i:=0; i<len(C); i++ {
    //Bc := C[i] // Pair b1 and b2 in B1
    //Br := Ru[i]
    Bp1 := pairing.NewG2().Sub(pairing.NewG2().SetBytes(C[i].b1),
      pairing.NewG2().SetBytes(Su[i].b1))
    Bp2 := pairing.NewG2().Sub(pairing.NewG2().SetBytes(C[i].b2),
      pairing.NewG2().SetBytes(Su[i].b2))

    Bi := IotaPrime2(pairing, y[i], sigma)
    Bi1 := pairing.NewG2().SetBytes(Bi.b1)
    Bi2 := pairing.NewG2().SetBytes(Bi.b2)

    tmp1 := Bp1.Equals(Bi1)
    tmp2 := Bp2.Equals(Bi2)

    if (verbose){
      fmt.Println("Testing Equality: ", i, tmp1 && tmp2)
      fmt.Printf("%s\t",Bp1)
      fmt.Printf("%s\n",Bp2)
      fmt.Printf("%s\t",Bi1)
      fmt.Printf("%s\n",Bi2)
    }
    ret2 = (ret2 && tmp1 && tmp2)
  }

  return ret1 && ret2
}


func TestEquation1ProofGen(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if verbose {fmt.Println("Test Proof Generation for Eq1")}

  Xc := pairing.NewZr().Rand() // Client Secret Key (variable)
  H  := pairing.NewG2().SetBytes(sharedParams.G2) // Shared Generator ??
  PKc := pairing.NewG2().Rand() // Public Key (variable

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  proof := ProveEquation1(pairing, Xc, H, PKc, sigma)

  if verbose {
    fmt.Println("P.Theta -------------------- ", proof.Theta)
    for i := 0; i < len(proof.Theta); i++ {
      fmt.Println("\tEl:", i, proof.Theta[i])
    }
    fmt.Println("P.Pi ----------------------- ", proof.Pi)
    for i := 0; i < len(proof.Pi); i++ {
      fmt.Println("\tEl:", i, proof.Pi[i])
    }
    fmt.Println("P.c ------------------------ ", proof.c)
    for i := 0; i < len(proof.c); i++ {
      fmt.Println("\tEl:", i, proof.c[i])
    }
    fmt.Println("P.d ------------------------ ", proof.d)
    for i := 0; i < len(proof.d); i++ {
      fmt.Println("\tEl:", i, proof.d[i])
    }
    fmt.Println("P.cprime ------------------- ", proof.cprime)
    for i := 0; i < len(proof.cprime); i++ {
      fmt.Println("\tEl:", i, proof.cprime[i])
    }
    fmt.Println("P.dprime ------------------- ", proof.dprime)
    for i := 0; i < len(proof.dprime); i++ {
      fmt.Println("\tEl:", i, proof.dprime[i])
    }
  }

  return len(proof.Theta) == 2 && len(proof.Pi) == 1 &&
      len(proof.d) == 1 && len(proof.cprime) == 1 &&
      len(proof.c) == 0 && len(proof.dprime) == 0
}

func TestEquation2ProofGen(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if verbose {fmt.Println("Test Proof Generation for Eq2")}

  rprime := pairing.NewZr().Rand() // Client Secret Key (variable)
  C  := pairing.NewG1().SetBytes(sharedParams.G1)
  G  := pairing.NewG1().SetBytes(sharedParams.G1) // Shared Generator ??
  //rprime := pairing.NewG2().Rand() // Public Key (variable

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  proof := ProveEquation2(pairing, rprime, G, C, sigma)

  if verbose {
    fmt.Println("P.Theta -------------------- ", proof.Theta)
    for i := 0; i < len(proof.Theta); i++ {
      fmt.Println("\tEl:", i, proof.Theta[i])
    }
    fmt.Println("P.Pi ----------------------- ", proof.Pi)
    for i := 0; i < len(proof.Pi); i++ {
      fmt.Println("\tEl:", i, proof.Pi[i])
    }
    fmt.Println("P.c ------------------------ ", proof.c)
    for i := 0; i < len(proof.c); i++ {
      fmt.Println("\tEl:", i, proof.c[i])
    }
    fmt.Println("P.d ------------------------ ", proof.d)
    for i := 0; i < len(proof.d); i++ {
      fmt.Println("\tEl:", i, proof.d[i])
    }
    fmt.Println("P.cprime ------------------- ", proof.cprime)
    for i := 0; i < len(proof.cprime); i++ {
      fmt.Println("\tEl:", i, proof.cprime[i])
    }
    fmt.Println("P.dprime ------------------- ", proof.dprime)
    for i := 0; i < len(proof.dprime); i++ {
      fmt.Println("\tEl:", i, proof.dprime[i])
    }
  }

  return len(proof.Theta) == 1 && len(proof.Pi) == 2 &&
      len(proof.d) == 0 && len(proof.cprime) == 0 &&
      len(proof.c) == 1 && len(proof.dprime) == 1
}

func TestEquation4ProofGen(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  VK, SK := SKeyGen(sharedParams)
  _ = VK

  P := new(Pseudonym)
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  P.C = pairing.NewG1().Rand().Bytes()
  P.D = pairing.NewG1().Rand().Bytes()
  PKc := new(ClientPublicKey)
  PKc.PK = pairing.NewG2().Rand().Bytes()
  ecert := SSign(sharedParams, SK, P, PKc)

  if verbose {fmt.Println("Test Proof Generation for Eq1")}

  //Vars
  R := pairing.NewG1().SetBytes(ecert.R)
  S := pairing.NewG1().SetBytes(ecert.S)
  C := pairing.NewG1().SetBytes(P.C)
  D := pairing.NewG1().SetBytes(P.D)

  //Constants
  V := pairing.NewG2().SetBytes(VK.V)
  H := pairing.NewG2().SetBytes(sharedParams.G2)
  W1 := pairing.NewG2().SetBytes(VK.W1)
  W2 := pairing.NewG2().SetBytes(VK.W2)


  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  proof := ProveEquation4(pairing, R, S, C, D, V, H, W1, W2, sigma)

  if verbose {
    fmt.Println("P.Theta -------------------- ", proof.Theta)
    for i := 0; i < len(proof.Theta); i++ {
      fmt.Println("\tEl:", i, proof.Theta[i])
    }
    fmt.Println("P.Pi ----------------------- ", proof.Pi)
    for i := 0; i < len(proof.Pi); i++ {
      fmt.Println("\tEl:", i, proof.Pi[i])
    }
    fmt.Println("P.c ------------------------ ", proof.c)
    for i := 0; i < len(proof.c); i++ {
      fmt.Println("\tEl:", i, proof.c[i])
    }
    fmt.Println("P.d ------------------------ ", proof.d)
    for i := 0; i < len(proof.d); i++ {
      fmt.Println("\tEl:", i, proof.d[i])
    }
    fmt.Println("P.cprime ------------------- ", proof.cprime)
    for i := 0; i < len(proof.cprime); i++ {
      fmt.Println("\tEl:", i, proof.cprime[i])
    }
    fmt.Println("P.dprime ------------------- ", proof.dprime)
    for i := 0; i < len(proof.dprime); i++ {
      fmt.Println("\tEl:", i, proof.dprime[i])
    }
  }

  return len(proof.Theta) == 2 && len(proof.Pi) == 2 &&
    len(proof.d) == 0 && len(proof.cprime) == 0 &&
    len(proof.c) == 4 && len(proof.dprime) == 0
}


func TestEquation5ProofGen(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  VK, SK := SKeyGen(sharedParams)
  _ = VK

  P := new(Pseudonym)
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  P.C = pairing.NewG1().Rand().Bytes()
  P.D = pairing.NewG1().Rand().Bytes()
  PKc := new(ClientPublicKey)
  PKc.PK = pairing.NewG2().Rand().Bytes()
  ecert := SSign(sharedParams, SK, P, PKc)

  if verbose {fmt.Println("Test Proof Generation for Eq1")}

  R := pairing.NewG1().SetBytes(ecert.R)
  T := pairing.NewG2().SetBytes(ecert.T)
  PK := pairing.NewG2().SetBytes(PKc.PK)
  U := pairing.NewG1().SetBytes(VK.U)

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  proof := ProveEquation5(pairing, R, T, PK, U, sigma)

  if verbose {
    fmt.Println("P.Theta -------------------- ", proof.Theta)
    for i := 0; i < len(proof.Theta); i++ {
      fmt.Println("\tEl:", i, proof.Theta[i])
    }
    fmt.Println("P.Pi ----------------------- ", proof.Pi)
    for i := 0; i < len(proof.Pi); i++ {
      fmt.Println("\tEl:", i, proof.Pi[i])
    }
    fmt.Println("P.c ------------------------ ", proof.c)
    for i := 0; i < len(proof.c); i++ {
      fmt.Println("\tEl:", i, proof.c[i])
    }
    fmt.Println("P.d ------------------------ ", proof.d)
    for i := 0; i < len(proof.d); i++ {
      fmt.Println("\tEl:", i, proof.d[i])
    }
    fmt.Println("P.cprime ------------------- ", proof.cprime)
    for i := 0; i < len(proof.cprime); i++ {
      fmt.Println("\tEl:", i, proof.cprime[i])
    }
    fmt.Println("P.dprime ------------------- ", proof.dprime)
    for i := 0; i < len(proof.dprime); i++ {
      fmt.Println("\tEl:", i, proof.dprime[i])
    }
  }

  return len(proof.Theta) == 2 && len(proof.Pi) == 2 &&
    len(proof.d) == 2 && len(proof.cprime) == 0 &&
    len(proof.c) == 1 && len(proof.dprime) == 0
}


func TestEquation1Verify(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if verbose {fmt.Println("Test Proof Generation for Eq1")}

  Xc := pairing.NewZr().Rand() // Client Secret Key (variable)
  H  := pairing.NewG2().SetBytes(sharedParams.G2) // Shared Generator ??
  PKc := pairing.NewG2().MulZn(H, Xc) // Public Key (variable
  negPKc := pairing.NewG2().Neg(PKc)

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  if verbose {fmt.Println("Generate Proof")}
  proof := ProveEquation1(pairing, Xc, H, negPKc, sigma)

  if verbose {fmt.Println("Tetsting Initital Euqation: XcH + (-1)PKc = 0")}
  tau := pairing.NewG2().Set0()
  if verbose {fmt.Println(tau)}


  if verbose {fmt.Println("Verify Proof")}
  ret := VerifyEquation1(pairing, proof, H, tau, sigma)

  if verbose {
    fmt.Println("Verify Restul: ", ret)
  }



  return ret
}

func TestEquation2Verify(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if verbose {fmt.Println("Creating Equations")}

  // Shared generators
  G  := pairing.NewG1().SetBytes(sharedParams.G1)
  H  := pairing.NewG2().SetBytes(sharedParams.G2)
  _ = H

  // Verifcation key
  VK, SK := SKeyGen(sharedParams)
  _ = VK
  _ = SK

  // Auditor Key pair
  PKa, SKa := EKeyGen(sharedParams)
  _ = SKa

  // Pseudonym Generation
  clientID := new(ClientID)
  clientID.ID = pairing.NewG1().Rand().Bytes()
  P := EEnc(sharedParams, PKa, clientID)
  C := pairing.NewG1().SetBytes(P.C)

  // Rerandomization
  Pprime, rprime := ERerand(sharedParams, PKa, P)
  r := pairing.NewZr().SetBytes(rprime)
  Cprime := pairing.NewG1().SetBytes(Pprime.C)

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  if verbose {fmt.Println("Generating proof:")}
  proof := ProveEquation2(
    pairing,
    r,
    G,
    C,
    sigma)

  if verbose {fmt.Println("Testing second Equation: C + rprime * G = Cprime")}
  tau := pairing.NewG1().MulZn(G, r)//pairing.NewG1().SetBytes(Pprime.C)
  tau = pairing.NewG1().Add(C, tau)
  if verbose {fmt.Println(tau)}

  if verbose {fmt.Println("Verify Proof")}
  ret := VerifyEquation2(pairing, proof, G, Cprime, sigma)

  if verbose {
    fmt.Println("Verify Result: ", ret)
  }

  return ret
}


func TestEquation3Verify(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if verbose {fmt.Println("Creating Equations")}

  // Shared generators
  G  := pairing.NewG1().SetBytes(sharedParams.G1)
  H  := pairing.NewG2().SetBytes(sharedParams.G2)
  _ = H
  _ = G

  // Verifcation key
  VK, SK := SKeyGen(sharedParams)

  _ = VK
  _ = SK

  // Auditor Key pair
  PKa, SKa := EKeyGen(sharedParams)
  PK := pairing.NewG1().SetBytes(PKa.PK)
  _ = SKa

  // Pseudonym Generation
  clientID := new(ClientID)
  clientID.ID = pairing.NewG1().Rand().Bytes()
  P := EEnc(sharedParams, PKa, clientID)
  D := pairing.NewG1().SetBytes(P.D)

  // Rerandomization
  Pprime, rprime := ERerand(sharedParams, PKa, P)
  r := pairing.NewZr().SetBytes(rprime)
  Dprime := pairing.NewG1().SetBytes(Pprime.D)

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  if verbose {fmt.Println("Generating proof:")}
  proof := ProveEquation2(
    pairing,
    r,
    PK,
    D,
    sigma)

  if verbose {fmt.Println("Testing second Equation: D + rprime * G = Dprime")}
  tau := pairing.NewG1().MulZn(PK, r)
  tau = pairing.NewG1().Add(D, tau)
  if verbose {fmt.Println(tau)}

  if verbose {fmt.Println("Verify Proof")}
  ret := VerifyEquation2(pairing, proof, PK, Dprime, sigma)

  if verbose {
    fmt.Println("Verify Result: ", ret)
  }

  return ret
}

func TestEquation4Verify(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  VK, SK := SKeyGen(sharedParams)
  _ = VK

  P := new(Pseudonym)
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  P.C = pairing.NewG1().Rand().Bytes()
  P.D = pairing.NewG1().Rand().Bytes()
  PKc := new(ClientPublicKey)
  PKc.PK = pairing.NewG2().Rand().Bytes()
  ecert := SSign(sharedParams, SK, P, PKc)

  if verbose {fmt.Println("Test Proof Generation for Eq1")}

  //Vars
  R := pairing.NewG1().SetBytes(ecert.R)
  S := pairing.NewG1().SetBytes(ecert.S)
  C := pairing.NewG1().SetBytes(P.C)
  D := pairing.NewG1().SetBytes(P.D)

  //Constants
  V := pairing.NewG2().SetBytes(VK.V)
  H := pairing.NewG2().SetBytes(sharedParams.G2)
  W1 := pairing.NewG2().SetBytes(VK.W1)
  W2 := pairing.NewG2().SetBytes(VK.W2)


  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  proof := ProveEquation4(pairing, R, S, C, D, V, H, W1, W2, sigma)

  if verbose {fmt.Println("Tetsting Initital Euqation: e(R, T)e(U, PKc) = e(G, G)")}
  tau := pairing.NewGT().Pair(pairing.NewG1().SetBytes(sharedParams.G1),
    pairing.NewG2().SetBytes(VK.Z))
  if verbose {fmt.Println(tau)}


  if verbose {fmt.Println("Verify Proof")}
  ret := VerifyEquation4(pairing, proof, V, H, W1, W2, tau, sigma)

  if verbose {
    fmt.Println("Verify Restul: ", ret)
  }

  return ret
}



func TestEquation5Verify(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  VK, SK := SKeyGen(sharedParams)
  _ = VK

  P := new(Pseudonym)
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  P.C = pairing.NewG1().Rand().Bytes()
  P.D = pairing.NewG1().Rand().Bytes()
  PKc := new(ClientPublicKey)
  PKc.PK = pairing.NewG2().Rand().Bytes()
  ecert := SSign(sharedParams, SK, P, PKc)

  eRT := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().SetBytes(ecert.T))
  if verbose {fmt.Println("eRT:", eRT)}

  eUPKc := pairing.NewGT().Pair(pairing.NewG1().SetBytes(VK.U),
    pairing.NewG2().SetBytes(PKc.PK))
  if verbose {fmt.Println("eUPKc:", eUPKc)}

  eGH := pairing.NewGT().Pair(pairing.NewG1().SetBytes(sharedParams.G1),
    pairing.NewG2().SetBytes(sharedParams.G2))
  if verbose {fmt.Println("eGH:", eGH)}

  // Test Revised
  eR1 := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().Set1())
  if verbose {fmt.Println("eR1:", eR1)}

  e1T := pairing.NewGT().Pair(pairing.NewG1().Set1(),
    pairing.NewG2().SetBytes(ecert.T))
  if verbose {fmt.Println("e1T:", e1T)}

  res := pairing.NewGT().Mul(eRT, eUPKc)
  res = pairing.NewGT().Mul(res, eR1)
  res = pairing.NewGT().Mul(res, e1T)
  if verbose {fmt.Println("Res:", res)}

  //////////////////////////////////////////////////////////////

  if verbose {fmt.Println("Test Proof Generation for Eq1")}
  R := pairing.NewG1().SetBytes(ecert.R)
  T := pairing.NewG2().SetBytes(ecert.T)
  PK := pairing.NewG2().SetBytes(PKc.PK)
  U := pairing.NewG1().SetBytes(VK.U)

  if verbose {fmt.Println("Creating CRS CommonReferenceString")}
  alpha := pairing.NewZr().Rand() // Another Secret Key..
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  if verbose {fmt.Println("Generate Proof")}
  proof := ProveEquation5(pairing, R, T, PK, U, sigma)

  if verbose {fmt.Println("Tetsting Initital Euqation: e(R, T)e(U, PKc) = e(G, G)")}
  tau := eGH
  if verbose {fmt.Println(tau)}


  if verbose {fmt.Println("Verify Proof")}
  ret := VerifyEquation5(pairing, proof, U, tau, sigma)

  if verbose {
    fmt.Println("Verify Restul: ", ret)
  }

  return ret
}





/*
 * Run test b times
 */
func RunPTest(b int) {
  for i := 0; i < b; i++ {
    if !Ptest(false) {
      os.Exit(1)
    }
  }
}

func TestEquation4Equality(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  VK, SK := SKeyGen(sharedParams)
  _ = VK

  P := new(Pseudonym)
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  P.C = pairing.NewG1().Rand().Bytes()
  P.D = pairing.NewG1().Rand().Bytes()
  PKc := new(ClientPublicKey)
  PKc.PK = pairing.NewG2().Rand().Bytes()
  ecert := SSign(sharedParams, SK, P, PKc)

  eRV := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().SetBytes(VK.V))
  if verbose {fmt.Println("eRV:", eRV)}

  eSH := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.S),
    pairing.NewG2().SetBytes(sharedParams.G2))
  if verbose {fmt.Println("eSH:", eSH)}

  eCW1 := pairing.NewGT().Pair(pairing.NewG1().SetBytes(P.C),
    pairing.NewG2().SetBytes(VK.W1))
  if verbose {fmt.Println("eCW1:", eCW1)}

  eDW2 := pairing.NewGT().Pair(pairing.NewG1().SetBytes(P.D),
    pairing.NewG2().SetBytes(VK.W2))
  if verbose {fmt.Println("eDW2:", eDW2)}

  eGZ := pairing.NewGT().Pair(pairing.NewG1().SetBytes(sharedParams.G1),
    pairing.NewG2().SetBytes(VK.Z))
  if verbose {fmt.Println("eGZ:", eGZ)}

  res := pairing.NewGT().Mul(eRV, eSH)
  res = pairing.NewGT().Mul(res, eCW1)
  res = pairing.NewGT().Mul(res, eDW2)
  if verbose {fmt.Println("Res:", res)}

  ret1 := res.Equals(eGZ)


  return ret1
}


func TestEquation5Equality(verbose bool) bool {
  sharedParams := GenerateSharedParams()
  VK, SK := SKeyGen(sharedParams)
  _ = VK

  P := new(Pseudonym)
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  P.C = pairing.NewG1().Rand().Bytes()
  P.D = pairing.NewG1().Rand().Bytes()
  PKc := new(ClientPublicKey)
  PKc.PK = pairing.NewG2().Rand().Bytes()
  ecert := SSign(sharedParams, SK, P, PKc)

  eRT := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().SetBytes(ecert.T))
  if verbose {fmt.Println("eRT:", eRT)}

  eUPKc := pairing.NewGT().Pair(pairing.NewG1().SetBytes(VK.U),
    pairing.NewG2().SetBytes(PKc.PK))
  if verbose {fmt.Println("eUPKc:", eUPKc)}

  eGH := pairing.NewGT().Pair(pairing.NewG1().SetBytes(sharedParams.G1),
    pairing.NewG2().SetBytes(sharedParams.G2))
  if verbose {fmt.Println("eGH:", eGH)}

  // Test Revised
  eR1 := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().Set1())
  if verbose {fmt.Println("eR1:", eR1)}
  e1T := pairing.NewGT().Pair(pairing.NewG1().Set1(),
    pairing.NewG2().SetBytes(ecert.T))
  if verbose {fmt.Println("e1T:", e1T)}

  eUPKcT := pairing.NewGT().Add(eUPKc, e1T)
  res := pairing.NewGT().Mul(eUPKcT, eR1)
  res = pairing.NewGT().Mul(res, eRT)
  if verbose {fmt.Println("Res:", res)}

  ret1 := res.Equals(eGH)

  /////////////////////////////////////////////

  eRPkc0 := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().MulZn(pairing.NewG2().SetBytes(PKc.PK),
      pairing.NewZr().SetInt32(0)))
  eRT1 := pairing.NewGT().Pair(pairing.NewG1().SetBytes(ecert.R),
    pairing.NewG2().MulZn(pairing.NewG2().SetBytes(ecert.T),
      pairing.NewZr().SetInt32(1)))

  eRTPKc01 := pairing.NewGT().Add(eRPkc0, eRT1)
  fmt.Println("eRTPKc01:",eRTPKc01)
  ret2 := eRTPKc01.Equals(eRT)


  /////////////////////////////////////////////
  eRand1 := pairing.NewGT().Pair(pairing.NewG1().Rand(), pairing.NewG2().Set1())
  fmt.Println("eRand1:", eRand1)
  eRand2 := pairing.NewGT().Pair(pairing.NewG1().Set1(), pairing.NewG2().Rand())
  fmt.Println("eRand2:", eRand2)
  eRandRet := pairing.NewGT().Add(res, eRand1)
  eRandRet = pairing.NewGT().Add(eRandRet, eRand2)
  fmt.Println("eRandRet:", eRandRet)
  ret3 := eRandRet.Equals(res)

  return ret1 && ret2 && ret3
}

/*
 * Run all Tests
 */
func RunAllPTests(verbose bool) {
  fmt.Println("Iota and Rho:         ", TestIotaRho(verbose))
  fmt.Println("Iota and Rho Prime:   ", TestIotaRhoPrime(verbose))
  fmt.Println("Iota Hat:             ", TestIotaHat(verbose))
  fmt.Println("F function Map:       ", TestFMap(verbose))
  fmt.Println("Matrix Map:           ", TestCompleteMatrixMapping(verbose))
  fmt.Println("Simple Commitment     ", TestSimpleCommitment(verbose))
  fmt.Println("Commitment: G1->B1    ", TestCreateCommitmentsG1(verbose))
  fmt.Println("Commitment: Zp->B1    ", TestCreateCommitmentPrimeOnG1(verbose))
  fmt.Println("Commitment: G2->B2    ", TestCreateCommitmentsG2(verbose))
  fmt.Println("Commitment: Zp->B2    ", TestCreateCommitmentPrimeOnG2(verbose))
  fmt.Println("Proof Generation EQ1  ", TestEquation1ProofGen(verbose))
  fmt.Println("Proof Generation EQ2-3", TestEquation2ProofGen(verbose))
  fmt.Println("Proof Generation EQ4  ", TestEquation4ProofGen(verbose))
  fmt.Println("Proof Generation EQ5  ", TestEquation5ProofGen(verbose))
  fmt.Println("Proof Verify EQ1      ", TestEquation1Verify(verbose))
  fmt.Println("Proof Verify EQ2      ", TestEquation2Verify(verbose))
  fmt.Println("Proof Verify EQ3      ", TestEquation3Verify(verbose))
  fmt.Println("Proof Verify EQ4      ", TestEquation4Verify(verbose))
  fmt.Println("Proof Verify EQ5      ", TestEquation5Verify(verbose))
  fmt.Println("Test Entire Pipeline  ", Ptest(verbose))

}
