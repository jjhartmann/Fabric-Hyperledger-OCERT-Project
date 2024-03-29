/*
 *
 * Copyright 2017 Kewei Shi, Jeremy Hartmann, Tuhin Tiwari and Dharvi Verma
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/*
 * Based on Non-interactive Proof Systems for Bilinear Groups
 */

package ocert

import (
    "github.com/Nik-U/pbc"
    "reflect"
)

/*
 * Set up the proof of knowledge, called by the client. It takes a system
 * of equations(e.g. pairing product equations and multi-scalar multiplication
 * equations) and outputs proof (e.g pi and theta ...)
 */
func PSetup(sharedParams *SharedParams, vars *ProofVariables) *ProofOfKnowledge {
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().Rand()
    g2 := pairing.NewG2().Rand()
    gt := pairing.NewGT().Pair(g1, g2)
    _ = gt

    pi := new(ProofOfKnowledge)
    alpha := pairing.NewZr().Rand() // Seed for CRS
    sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

    // Setup proof of eq1
    Xc := pairing.NewZr().SetBytes(vars.Xc)
    H := pairing.NewG2().SetBytes(sharedParams.G2)
    negPKc := pairing.NewG2().Neg(pairing.NewG2().SetBytes(vars.PKc.PK))
    pi.Eq1 = ProveEquation1(pairing, Xc, H, negPKc, sigma)

    // Setup proof of eq2
    C := pairing.NewG1().SetBytes(vars.P.C)
    G := pairing.NewG1().SetBytes(sharedParams.G1)
    rprime := pairing.NewZr().SetBytes(vars.RPrime)
    pi.Eq2 = ProveEquation2(pairing, rprime, G, C, sigma)

    // Setup proof of eq3
    D := pairing.NewG1().SetBytes(vars.P.D)
    PKa := pairing.NewG1().SetBytes(vars.PKa.PK)
    pi.Eq3 = ProveEquation2(pairing, rprime, PKa, D, sigma)

    // Setup proof of eq4
    // Vars
    R := pairing.NewG1().SetBytes(vars.E.R)
    S := pairing.NewG1().SetBytes(vars.E.S)
    _ = C
    _ = D

    // Constants
    V := pairing.NewG2().SetBytes(vars.VK.V)
    W1 := pairing.NewG2().SetBytes(vars.VK.W1)
    W2 := pairing.NewG2().SetBytes(vars.VK.W2)
    _ = H

    pi.Eq4 = ProveEquation4(pairing, R, S, C, D, V, H, W1, W2, sigma)

    // Setup proof of eq5
    _ = R
    T := pairing.NewG2().SetBytes(vars.E.T)
    PKc := pairing.NewG2().SetBytes(vars.PKc.PK)
    U := pairing.NewG1().SetBytes(vars.VK.U)

    pi.Eq5 = ProveEquation5(pairing, R, T, PKc, U, sigma)

    // Set CRSf
    pi.sigma = sigma
    return pi
}

/*
 * Validate the proof of knowledage, return true if all the equations
 * in the system hold.
 */
func PProve(sharedParams *SharedParams, pi *ProofOfKnowledge, consts *ProofConstants) bool {
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().Rand()
    g2 := pairing.NewG2().Rand()
    gt := pairing.NewGT().Pair(g1, g2)
    _ = gt

    // Validate eq1
    H := pairing.NewG2().SetBytes(sharedParams.G2)
    Zero := pairing.NewG2().Set0()
    retVal := VerifyEquation1(pairing, pi.Eq1, H, Zero, pi.sigma)

    // fmt.Println("EQ1:", retVal)

    // Validate eq2
    G := pairing.NewG1().SetBytes(sharedParams.G1)
    Cprime := pairing.NewG1().SetBytes(consts.PPrime.C)
    retVal2 := VerifyEquation2(pairing, pi.Eq2, G, Cprime, pi.sigma)
    retVal = retVal && retVal2

    // fmt.Println("EQ2:", retVal2)

    // Validate eq3
    PKa := pairing.NewG1().SetBytes(consts.PKa.PK)
    Dprime := pairing.NewG1().SetBytes(consts.PPrime.D)
    retVal3 := VerifyEquation2(pairing, pi.Eq3, PKa, Dprime, pi.sigma)
    retVal = retVal && retVal3

    // fmt.Println("EQ3:", retVal3)

    // Validate eq4
    V := pairing.NewG2().SetBytes(consts.VK.V)
    _ = H
    W1 := pairing.NewG2().SetBytes(consts.VK.W1)
    W2 := pairing.NewG2().SetBytes(consts.VK.W2)
    eGZ := pairing.NewGT().SetBytes(consts.Egz)
    retVal4 := VerifyEquation4(pairing, pi.Eq4, V, H, W1, W2, eGZ, pi.sigma)
    retVal = retVal && retVal4

    // fmt.Println("EQ4:", retVal4)

    // Validate eq5
    U := pairing.NewG1().SetBytes(consts.VK.U)
    eGH := pairing.NewGT().SetBytes(consts.Egh)
    _ = U
    _ = eGH
    retVal5 := VerifyEquation5(pairing, pi.Eq5, U, eGH, pi.sigma)
    retVal = retVal && retVal5

    // fmt.Println("EQ5:", retVal)
    
    return retVal
}


/*
 * Create proof for equation: xc * H + (-1)PKc = 0
 *   Multi-Scalar Multiplication in G2
 *   xc, (-1) from group Zp: Zp -> B1
 *   H, PKc from group G2:
 *
 * Proof:
 *    Pi     := r*ι_2(H) + r*lambda*ι_2(PKc)  + (r*lambda*S - T')*V
 *    Theta  := S'*ι'_1(-1) + S*lambda*ι'_1(xc) + Tu_1
 */
func ProveEquation1(pairing *pbc.Pairing, xc *pbc.Element, H *pbc.Element, PKc *pbc.Element, sigma *CommonReferenceString) *ProofOfEquation{
    proof := new(ProofOfEquation)

    // Create commitment in B1 for Xc
    cprime, _, R := CreateCommitmentPrimeOnG1(pairing, []*pbc.Element{xc}, sigma)
    if R.cols != 1 && R.rows != 1 {
        panic("Issues in conversion and creation of samples in Zp for R")
    }
    r := R.mat[0][0]

    // Create commitment in B2 for PKc
    d, _, S := CreateCommitmentOnG2(pairing, []*pbc.Element{PKc}, sigma)
    if S.rows != 1 && S.cols != 2 {
        panic("Issues in conversion and creation of samples in Zp for S")
    }

    /////////////////////////////////////
    // Pi: In G2
    /////////////////////////////////////

    // Multiply Scalar from Zn on B elements
    Hi := Iota2(pairing, H)
    Hir := Hi.MulScalarInG2(pairing, r)           // r*ι_2(H)
    // +
    // Create Phi := (r*lambada*S - T)
    T := NewRMatrix(pairing, 2, 1)
    Ti := T.InvertMatrix()                 // T' invert
    Vphi := Ti.MulCommitmentKeysG2(pairing, sigma.V) // (r*gamma*S - T')V (commitment key in G2)
    if len(Vphi) > 1{
        panic("VPhi Should have len == 1")
    }
    //=
    // Construct Pi (Hir - Vphi)
    Tv := Vphi[0]
    pi := new(BPair)
    pi.b1 = pairing.NewG2().Sub(pairing.NewG2().SetBytes(Hir.b1),
        pairing.NewG2().SetBytes(Tv.b1)).Bytes()
    pi.b2 = pairing.NewG2().Sub(pairing.NewG2().SetBytes(Hir.b2),
        pairing.NewG2().SetBytes(Tv.b2)).Bytes()

    ////////////////////////////////////////////
    // Theta: In G1
    ////////////////////////////////////////////
    Si := S.InvertMatrix()                        // S Invert = S'
    pos1 := IotaPrime1(pairing, pairing.NewZr().Set1(), sigma)
    Spos := Si.MulBScalarinB1(pairing, *pos1)     // S'*ι'_1(1)
    // +
    Tu := T.MulCommitmentKeysG1(pairing, []CommitmentKey{sigma.U[0]}) // Tu_1

    if len(Spos) != len(Tu) {
        panic("All section lengths need to be equivalent")
    }

    // Construct theta
    theta := []*BPair{}

    for i := 0; i < len(Spos); i++ {
        Snegxc := Spos[i][0].AddinG1(pairing, Tu[i])
        theta = append(theta, Snegxc)
    }

    // Collect elements
    proof.Theta = theta
    proof.Pi = []*BPair{pi}
    proof.cprime = cprime
    proof.d = d

    return proof
}

/*
 * Create proof for equation: C + r' * G = C'
 *   Multi-Scalar Multiplication in G1
 *   r' from group Zp: Zp -> B1
 *   C, G from group G1:
 *   C' is a Constant
 *
 * Proof:
 *    Pi     := R'*ι'_2(1) + R'*lambda*ι'_2(rprime)  + (R'*lambda*s - T')*v_1
 *    Theta  := s*ι_1(C) + s*lambda*ι_1(G) + TU
 */
func ProveEquation2(pairing *pbc.Pairing, rprime *pbc.Element, G *pbc.Element, C *pbc.Element, sigma *CommonReferenceString) *ProofOfEquation{
    proof := new(ProofOfEquation)

    // Create commitment in B1 for C
    c, _, R := CreateCommitmentOnG1(pairing, []*pbc.Element{C}, sigma)
    if R.rows != 1 && R.cols != 2 {
        panic("Issues in conversion and creation of samples in Zp for R")
    }

    // Create commitment in B2 for r' TODO: This should be in B2 - Done : Check?
    dprime, _, S := CreateCommitmentPrimeOnG2(pairing, []*pbc.Element{rprime}, sigma)
    if S.cols != 1 && S.rows != 1 {
        panic("Issues in conversion and creation of samples in Zp for S")
    }
    s := S.mat[0][0]

    // Convert parameters to B groups

    ////////////////////////////////////////////
    // Pi: In G2
    ////////////////////////////////////////////

    Ri := R.InvertMatrix()                        // R Invert = R'
    pos1 := IotaPrime2(pairing, pairing.NewZr().Set1(), sigma)
    Rpos := Ri.MulBScalarinB2(pairing, *pos1)     // R'*ι'_2(1)
    // +
    T := NewRMatrix(pairing, 1, 2)
    Ti := T.InvertMatrix()                 // T' invert
    Tv := Ti.MulCommitmentKeysG2(pairing, []CommitmentKey{sigma.V[0]})

    if len(Rpos) != len(Tv) {
        panic("All section lengths need to be equivalent")
    }

    // Construct pi
    pi := []*BPair{}

    for i := 0; i < len(Rpos); i++ {
        r_tmp := Rpos[i][0]
        tv_tmp := Tv[i]

        Bp1 := pairing.NewG2().Sub(pairing.NewG2().SetBytes(r_tmp.b1),
            pairing.NewG2().SetBytes(tv_tmp.b1))
        Bp2 := pairing.NewG2().Sub(pairing.NewG2().SetBytes(r_tmp.b2),
            pairing.NewG2().SetBytes(tv_tmp.b2))

        Bpair_tmp := new(BPair)
        Bpair_tmp.b1 = Bp1.Bytes()
        Bpair_tmp.b2 = Bp2.Bytes()

        pi = append(pi, Bpair_tmp)
    }

    /////////////////////////////////////
    // Theta: In G1
    /////////////////////////////////////

    // Multiply Scalar from Zn on B elements
    Gi := Iota1(pairing, G)
    Gir := Gi.MulScalarInG1(pairing, s)           // s*ι_1(C)

    // Multiple Phi by commitment keys
    Thetai := T.MulCommitmentKeysG1(pairing, sigma.U) // TU (commitment key in G1)
    if len(Thetai) > 1{
        panic("Thetai Should have len == 1")
    }

    // Construct theta (Gir + Cir + Thetai)
    theta := Gir.AddinG1(pairing, Thetai[0])

    // Collect elements
    proof.Theta = []*BPair {theta}
    proof.Pi = pi
    proof.c = c
    proof.dprime = dprime

    return proof
}



/*
 * Proof Equation 4
 */
func ProveEquation4(pairing *pbc.Pairing,
    R *pbc.Element,
    S *pbc.Element,
    C *pbc.Element,
    D *pbc.Element,
    V *pbc.Element,
    H *pbc.Element,
    W1 *pbc.Element,
    W2 *pbc.Element,
    sigma *CommonReferenceString) *ProofOfEquation {
    proof := new(ProofOfEquation)

    //// Create commitment in B2
    Xvec := []*pbc.Element{R, S, C, D}
    c_commit, _, Rmat := CreateCommitmentOnG1(pairing, Xvec, sigma) // T, PKc
    if Rmat.rows != 4 && Rmat.cols != 2 && len(c_commit) == 4 {
        panic("Issues in conversion and creation of samples in Zp for S")
    }

    //////////////////////////////////
    // PI: for G2
    /////////////////////////////////////

    Ri := Rmat.InvertMatrix()
    if len(Ri.mat) !=2 && len(Ri.mat[0]) != 4 {
        panic("Issue when inverting the R matrix")
    }

    B := new(BMatrix)
    B.mat = [][]*BPair{
        []*BPair{Iota2(pairing, V)},
        []*BPair{Iota2(pairing, H)},
        []*BPair{Iota2(pairing, W1)},
        []*BPair{Iota2(pairing, W2)},
    }
    B.rows = 4
    B.cols = 1

    RB := Ri.MultBPairMatrixG2(pairing, B)
    if len(RB.mat) != 2 && len(RB.mat[0]) != 1 {
        panic("Issue in dimensionality when mult in B matrix")
    }

    // -

    Tmat := NewRMatrix(pairing, 2, 2)
    Ti := Tmat.InvertMatrix()
    Tv := Ti.MulCommitmentKeysG2(pairing, sigma.V)
    if len(Tv) != 2 {
        panic("Issuen in commitment key multiplication.")
    }

    // ) =  Construct pi
    if len(Tv) != len(RB.mat) {
        panic("All preliminary output for PI needs to have equivalent dimensionality")
    }
    pi := []*BPair{}
    _ = RB      // *BMatrix 2x1
    _ = Tv    // []*BPair 2x1

    for i := 0; i < len(Tv); i++ {
        tmpRB := RB.mat[i][0]
        tmpTv := Tv[i]

        tmpBPair := new(BPair)
        tmpBPair.b1 = pairing.NewG2().Sub(pairing.NewG2().SetBytes(tmpRB.b1),
            pairing.NewG2().SetBytes(tmpTv.b1)).Bytes()
        tmpBPair.b2 = pairing.NewG2().Sub(pairing.NewG2().SetBytes(tmpRB.b2),
            pairing.NewG2().SetBytes(tmpTv.b2)).Bytes()

        pi = append(pi, tmpBPair)
    }


    //////////////////////////////////
    // Theta: for G1
    /////////////////////////////////////

    theta := Tmat.MulCommitmentKeysG1(pairing, sigma.U)


    proof.c = c_commit
    proof.Pi = pi
    proof.Theta = theta
    return proof
}


/*
 * Proof Equation 5
 */
func ProveEquation5(pairing *pbc.Pairing,
    R *pbc.Element,
    T *pbc.Element,
    PKc *pbc.Element,
    U *pbc.Element,
    sigma *CommonReferenceString) *ProofOfEquation {
    proof := new(ProofOfEquation)

    // Create commitment in B1 for Xc
    c, _, Rmat := CreateCommitmentOnG1(pairing, []*pbc.Element{R}, sigma)
    if Rmat.cols != 2 && Rmat.rows != 1 {
        panic("Issues in conversion and creation of samples in Zp for R")
    }

    // Create commitment in B2 for PKc
    Ymat := [][]*BPair{
        {Iota2(pairing, T)},
        {Iota2(pairing, PKc)},
    }
    Yvec := []*pbc.Element{T, PKc}
    d, _, Smat := CreateCommitmentOnG2(pairing, Yvec, sigma) // T, PKc
    if Smat.rows != 2 && Smat.cols != 2 {
        panic("Issues in conversion and creation of samples in Zp for S")
    }

    //////////////////////////////////
    // PI: for G2
    /////////////////////////////////////
    //Gamma := NewIdentiyMatrix(pairing, 1, 2)
    Gamma := new(RMatrix)
    Gamma.mat = [][]*pbc.Element{
        {pairing.NewZr().Set1(), pairing.NewZr().Set0()},
    }
    Gamma.rows = 1
    Gamma.cols = 2
    //fmt.Println("Gamma:", Gamma.mat)

    Ri := Rmat.InvertMatrix()
    Tmat := NewRMatrix(pairing, 2, 2)
    Ti := Tmat.InvertMatrix()

    Biota := Iota2(pairing, pairing.NewG2().Set0())
    Rb := Ri.MulBScalarinB2(pairing, *Biota)
    if len(Rb) != 2 && len(Rb[0]) != 1 {
        panic("Issue in dimensionality of RB")
    }

    // +

    YBMat := new(BMatrix) // TODO: Can clean this up
    YBMat.mat = Ymat
    YBMat.rows = 2
    YBMat.cols = 1
    Ygamma := Gamma.MultBPairMatrixG2(pairing, YBMat)
    RYg := Ri.MultBPairMatrixG2(pairing, Ygamma)
    if RYg.rows != 2 && RYg.cols != 1 {
        panic("Issue when multipling equations")
    }

    // + (

    Sg := Gamma.MultElementArrayZr(pairing, Smat.mat)
    if Sg.rows != 1 && Sg.cols != 2 {
        panic("Sg has wrong dimensionality.")
    }
    RSg := Ri.MultElementArrayZr(pairing, Sg.mat)
    if RSg.rows != 2 && RSg.cols != 2 {
        panic("RSg does not have correct dimensionality.")
    }
    RSgT := RSg.ElementWiseSub(pairing, Ti);
    RSgTv := RSgT.MulCommitmentKeysG2(pairing, sigma.V)
    if len(RSgTv) != 2 {
        panic("RSgTv dimensionality is incorrect should be 2")
    }

    // ) =  Construct pi

    if len(RSgTv) != RYg.rows && RYg.rows  == len(Rb) {
        panic("All preliminary output for PI needs to have equivalent dimensionality")
    }
    pi := []*BPair{}
    _ = Rb       // [][]*BPair
    _ = RYg      // *BMatrix 2x1
    _ = RSgTv    // []*BPair 2x1

    for i := 0; i < len(RSgTv); i++ {
        tmpRb := Rb[i][0]
        tmpRy := RYg.mat[i][0]
        tmpRs := RSgTv[i]

        tmpBPair := tmpRy.AddinG2(pairing, tmpRs)
        tmpBPair  = tmpBPair.AddinG2(pairing, tmpRb)
        pi = append(pi, tmpBPair)
    }

    //////////////////////////////////
    // Theta: for G1
    /////////////////////////////////////
    Si := Smat.InvertMatrix()
    Gi := Gamma.InvertMatrix()
    //fmt.Println("Gammai:", Gi.mat)

    Amat := new(BMatrix)
    Amat.mat = [][]*BPair{
        {Iota1(pairing, pairing.NewG1().Set0())},
        {Iota1(pairing, U)},
    }
    Amat.rows = 2
    Amat.cols = 1

    SiA := Si.MultBPairMatrixG1(pairing, Amat)
    if len(SiA.mat) != 2 && len(SiA.mat[0]) != 1 {
        panic("SiA dimensionality is wrong. Needs to be 2x1")
    }

    // +

    Riota := Iota1(pairing, R)
    RG := Gi.MulBScalarinB1(pairing, *Riota)

    RGmat := new(BMatrix)
    RGmat.mat = RG
    RGmat.rows = 2
    RGmat.cols = 1
    if len(RG) != 2 && len(RG[0]) != 1 {
        panic("RG dimensionality is wrong. Needs to be 2x1")
    }
    SRG := Si.MultBPairMatrixG1(pairing, RGmat)
    if len(SRG.mat) != 2 && len(SRG.mat[0]) != 1 {
        panic("SRG dimensionality is wrong. Needs to be 2x1")
    }
    // +
    Tu := Tmat.MulCommitmentKeysG1(pairing, sigma.U)
    if len(Tu) != 2 {
        panic("Tu dimensionality is wrong. Needs to be len 2")
    }
    // =
    if len(SiA.mat) != len(SRG.mat) && len(SRG.mat) != len(Tu){
        panic("Equation dimensionality is wrong. Needs to be len 2x1")
    }
    theta := []*BPair{}
    _ = SiA    // *BMatrix 2x1 // TODO: this can be streamlined
    _ = SRG    // *BMatrix 2x1
    _ = Tu     // []*BPair 2x1

    for i := 0; i < len(Tu); i++ {
        tmpRp := SiA.mat[i][0]
        tmpRy := SRG.mat[i][0]
        tmpRs := Tu[i]

        tmpBPair := tmpRp.AddinG1(pairing, tmpRy)
        tmpBPair = tmpBPair.AddinG1(pairing, tmpRs)

        theta = append(theta, tmpBPair)
    }

    proof.c = c
    proof.d = d
    proof.Pi = pi
    proof.Theta = theta
    return proof
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
 * Verifiy Equation 1
 *
 */
func VerifyEquation1(pairing *pbc.Pairing, proof *ProofOfEquation, H *pbc.Element, tau *pbc.Element, sigma *CommonReferenceString) bool {

    // Construct LHS
    pos1 := IotaPrime1(pairing, pairing.NewZr().Set1(), sigma)
    Fid := FMap(pairing, pos1, proof.d[0])
    // +
    Hi := Iota2(pairing, H)
    FcH := FMap(pairing, proof.cprime[0], Hi)
    // =
    LHS := Fid.AddinGT(pairing, FcH)

    // Construct RHS
    taui := IotaHat(pairing, tau, sigma)
    // +
    u := sigma.U[0].ConvertToBPair()
    Fup := FMap(pairing, u, proof.Pi[0])
    // +
    v1 := sigma.V[0].ConvertToBPair()
    Fsv1 := FMap(pairing, proof.Theta[0], v1)
    //+
    v2 := sigma.V[1].ConvertToBPair()
    Fsv2 := FMap(pairing, proof.Theta[1], v2)
    //=
    tmp2 := taui.AddinGT(pairing, Fup)
    tmp2 = tmp2.AddinGT(pairing, Fsv1)
    RHS := tmp2.AddinGT(pairing, Fsv2)

    // Perform Equality //TODO: Test for nil == nill
    ret := reflect.DeepEqual(LHS, RHS)
    return ret
}

/*
 * Verify Equation 2
 *
 */
func VerifyEquation2(pairing *pbc.Pairing, proof *ProofOfEquation, G *pbc.Element, tau *pbc.Element, sigma *CommonReferenceString) bool {

    // Construct LHS
    Gi := Iota1(pairing, G)
    FiGdp := FMap(pairing, Gi, proof.dprime[0])
    //+
    pos1 := IotaPrime2(pairing, pairing.NewZr().Set1(), sigma)
    Fcpos := FMap(pairing, proof.c[0], pos1)
    // =
    LHS := FiGdp.AddinGT(pairing, Fcpos)

    // Construct RHS
    taui := IotaHat2(pairing, tau, sigma)
    // +
    v := sigma.V[0].ConvertToBPair()
    Fvp := FMap(pairing, proof.Theta[0], v)
    // +
    u1 := sigma.U[0].ConvertToBPair()
    Fup1 := FMap(pairing, u1, proof.Pi[0])
    //+
    u2 := sigma.U[1].ConvertToBPair()
    Fup2 := FMap(pairing, u2, proof.Pi[1])
    //=
    tmp2 := taui.AddinGT(pairing, Fup1)
    tmp2 = tmp2.AddinGT(pairing, Fup2)
    RHS := tmp2.AddinGT(pairing, Fvp)

    // Perform Equality //TODO: Test for nil == nill
    ret := reflect.DeepEqual(LHS, RHS)
    return ret
}

/*
 * Verify Equation 4
 */
func VerifyEquation4(
    pairing *pbc.Pairing,
    proof *ProofOfEquation,
    V *pbc.Element,
    H *pbc.Element,
    W1 *pbc.Element,
    W2 *pbc.Element,
    tau *pbc.Element,
    sigma *CommonReferenceString) bool {

    // Construct LHS
    Vi := Iota2(pairing, V)
    Fv := FMap(pairing, proof.c[0], Vi)
    // +
    Hi := Iota2(pairing, H)
    Fh := FMap(pairing, proof.c[1], Hi)
    // +
    W1i := Iota2(pairing, W1)
    Fw1 := FMap(pairing, proof.c[2], W1i)
    // +
    W2i := Iota2(pairing, W2)
    Fw2 := FMap(pairing, proof.c[3], W2i)
    // =
    LHS := Fv.AddinGT(pairing, Fh)
    LHS = LHS.AddinGT(pairing, Fw1)
    LHS = LHS.AddinGT(pairing, Fw2)

    // Construct RHS
    Taui := IotaT(pairing, tau)
    // +
    u1 := sigma.U[0].ConvertToBPair()
    Fu1 := FMap(pairing, u1, proof.Pi[0])
    //+
    u2 := sigma.U[1].ConvertToBPair()
    Fu2 := FMap(pairing, u2, proof.Pi[1])
    // +
    v1 := sigma.V[0].ConvertToBPair()
    Fv1 := FMap(pairing, proof.Theta[0], v1)
    // +
    v2 := sigma.V[1].ConvertToBPair()
    Fv2 := FMap(pairing, proof.Theta[1], v2)
    // +
    RHS := Taui.AddinGT(pairing, Fu1)
    RHS = RHS.AddinGT(pairing, Fu2)
    RHS = RHS.AddinGT(pairing, Fv1)
    RHS = RHS.AddinGT(pairing, Fv2)

    return reflect.DeepEqual(LHS, RHS)
}


/*
 * Verify Equation 5
 */
func VerifyEquation5(pairing *pbc.Pairing, proof *ProofOfEquation, U *pbc.Element, tau *pbc.Element, sigma *CommonReferenceString) bool {

    // Construct LHS
    Uiota := Iota1(pairing, U)
    Fu := FMap(pairing, Uiota, proof.d[1])
    // +
    Fcd := FMap(pairing, proof.c[0], proof.d[0])
    // =
    LHS := Fu.AddinGT(pairing, Fcd)

    // Construct RHS
    Taui := IotaT(pairing, tau)
    // +
    u1 := sigma.U[0].ConvertToBPair()
    Fu1 := FMap(pairing, u1, proof.Pi[0])
    //+
    u2 := sigma.U[1].ConvertToBPair()
    Fu2 := FMap(pairing, u2, proof.Pi[1])
    // +
    v1 := sigma.V[0].ConvertToBPair()
    Fv1 := FMap(pairing, proof.Theta[0], v1)
    // +
    v2 := sigma.V[1].ConvertToBPair()
    Fv2 := FMap(pairing, proof.Theta[1], v2)
    // +
    RHS := Taui.AddinGT(pairing, Fu1)
    RHS = RHS.AddinGT(pairing, Fu2)
    RHS = RHS.AddinGT(pairing, Fv1)
    RHS = RHS.AddinGT(pairing, Fv2)

    return reflect.DeepEqual(RHS, LHS)
}


/*
 * Create Common refernce string sigma.
 * sigma = (u1, u2, v1, v2)
 * u1 and u1 are in B1
 * v1 and v2 are in B2
 *
 * u1 = (O, P)
 * u2 = t * u1
 */
func CreateCommonReferenceString(sharedParams *SharedParams, alpha *pbc.Element) *CommonReferenceString {
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

    // Proof should use different generators then what is stored in the
    // params. Florian: Using the same generators could case security issues
    // due to the discrete logrihtm problem
    // Since the groups are cyclic, this should not matter.
    g1 := pairing.NewG1().Rand()
    g2 := pairing.NewG2().Rand()
    sigma := new(CommonReferenceString)

    // Create commit keys for u1 and u2 on G1
    u11 := g1.Bytes()
    u12 := pairing.NewG1().MulZn(g1, alpha)

    t := pairing.NewZr().Rand()
    u21 := pairing.NewG1().MulZn(g1, t)
    u22 := pairing.NewG1().MulZn(u12, t)

    sigma.U = []CommitmentKey{
        CommitmentKey{u11, u12.Bytes()},
        CommitmentKey{u21.Bytes(), u22.Bytes()},
    }

    // Create commit keys v1 and v2 on G2
    v11 := g2.Bytes()
    v12 := pairing.NewG2().MulZn(g2, alpha)

    t2 := pairing.NewZr().Rand()
    v21 := pairing.NewG2().MulZn(g2, t2)
    v22 := pairing.NewG2().MulZn(v12, t2)

    sigma.V = []CommitmentKey{
        CommitmentKey{v11, v12.Bytes()},
        CommitmentKey{v21.Bytes(), v22.Bytes()},
    }

    // Create commitment keys u and v on Zn
    su1 := pairing.NewG1().Add(u21, pairing.NewG1().Set0())
    su2 := pairing.NewG1().Add(u22, g1)
    sigma.u = CommitmentKey{su1.Bytes(), su2.Bytes()}

    sv1 := pairing.NewG2().Add(v21, pairing.NewG2().Set0())
    sv2 := pairing.NewG2().Add(v22, g2)
    sigma.v = CommitmentKey{sv1.Bytes(), sv2.Bytes()}

    return sigma
}


/*
 * Create Commitment: G1 -> B1
 * - Creates a commitment of a variable from G1 to B1
 *   c := ι1(X) + Ru
 */
func CreateCommitmentOnG1(pairing *pbc.Pairing, chi []*pbc.Element, sigma *CommonReferenceString) ([]*BPair, []*BPair, *RMatrix){

    // Create RMatrix of random elements
    rows := len(chi)
    cols := len(sigma.U)
    rmat := NewRMatrix(pairing, rows, cols)

    // Create pairs in B1
    Ru := rmat.MulCommitmentKeysG1(pairing, sigma.U)

    // Create Commitment container
    C := []*BPair{}
    if (len(Ru) != len(chi)){
        panic("Error in CreateCommitmentOnG1: Ru and X needs to have the same length")
    }

    // Build commitments in B1
    for i:=0; i<len(chi); i++ {
        tmp := Iota1(pairing, chi[i])
        B := tmp.AddinG1(pairing, Ru[i])
        C = append(C, B)
    }

    return C, Ru, rmat
}

/*
 * Create Commitment: Zp -> B1
 * - Creates a commitment of a variable from Zp to B1
 *   c := ι'1(X) + Ru
 *
 *  x: Is in Zp
 */
func CreateCommitmentPrimeOnG1(pairing *pbc.Pairing, x []*pbc.Element, sigma *CommonReferenceString) ([]*BPair, []*BPair, *RMatrix){

    // Create RMatrix of random elements
    rows := len(x)
    cols := 1
    rmat := NewRMatrix(pairing, rows, cols)

    // Create pairs in B1
    Ru := rmat.MulCommitmentKeysG1(pairing, []CommitmentKey{sigma.U[0]})

    // Create Commitment container
    C := []*BPair{}
    if (len(Ru) != len(x)){
        panic("Error in CreateCommitmentOnG1: Ru and X needs to have the same length")
    }

    // Build commitments in B1
    for i:=0; i<len(x); i++ {
        tmp := IotaPrime1(pairing, x[i], sigma)
        B := tmp.AddinG1(pairing, Ru[i])
        C = append(C, B)
    }

    return C, Ru, rmat
}

/*
 * Create Commitment: G2 -> B2
 * - Creates a commitment of a variable from G1 to B1
 *   c := ι1(X) + Ru
 */
func CreateCommitmentOnG2(pairing *pbc.Pairing, Y []*pbc.Element, sigma *CommonReferenceString) ([]*BPair, []*BPair, *RMatrix){

    // Create RMatrix of random elements
    rows := len(Y)
    cols := len(sigma.V)
    rmat := NewRMatrix(pairing, rows, cols)

    // Create pairs in B1
    Su := rmat.MulCommitmentKeysG2(pairing, sigma.V)

    // Create Commitment container
    C := []*BPair{}
    if (len(Su) != len(Y)){
        panic("Error in CreateCommitmentOnG1: Ru and X needs to have the same length")
    }

    // Build commitments in B1
    for i:=0; i<len(Y); i++ {
        tmp := Iota2(pairing, Y[i])
        B := tmp.AddinG2(pairing, Su[i])
        C = append(C, B)
    }

    return C, Su, rmat
}

/*
 * Create Commitment: Zp -> B2
 * - Creates a commitment of a variable from Zp to B2
 *   c := ι'2(X) + Ru
 *
 *  x: Is in Zp
 */
func CreateCommitmentPrimeOnG2(pairing *pbc.Pairing, y []*pbc.Element, sigma *CommonReferenceString) ([]*BPair, []*BPair, *RMatrix){

    // Create RMatrix of random elements
    rows := len(y)
    cols := 1
    rmat := NewRMatrix(pairing, rows, cols)

    // Create pairs in B1
    Su := rmat.MulCommitmentKeysG2(pairing, []CommitmentKey{sigma.V[0]})

    // Create Commitment container
    C := []*BPair{}
    if (len(Su) != len(y)){
        panic("Error in CreateCommitmentOnG1: Ru and X needs to have the same length")
    }

    // Build commitments in B1
    for i:=0; i<len(y); i++ {
        tmp := IotaPrime2(pairing, y[i], sigma)
        B := tmp.AddinG2(pairing, Su[i])
        C = append(C, B)
    }

    return C, Su, rmat
}

/*
    Multi-Scalar Multiplication Mapping for G1
    f: (x, Y) -> xY
 */
func MultiScalar_f_G1_map(pairing *pbc.Pairing, y *pbc.Element, X *pbc.Element) *pbc.Element{
    return pairing.NewG1().MulZn(X, y)
}

/*
    Multi-Scalar Multiplication Mapping for G2
    f: (x, Y) -> xY
 */
func MultiScalar_f_G2_map(pairing *pbc.Pairing, x *pbc.Element, Y *pbc.Element) *pbc.Element{
    return pairing.NewG2().MulZn(Y, x)
}

func ProductPairing_e_GT_map(pairing *pbc.Pairing, X *pbc.Element, Y *pbc.Element) *pbc.Element{
    return pairing.NewGT().Pair(X, Y)
}

/*
 Mapping between B1 and B2 to BT (Groth & Sahai p. 25)
 B1 in A1^2
 B2 in A2^2
 BTMat in AT^4
 */
func FMap(pairing *pbc.Pairing, B1 *BPair, B2 *BPair) *BTMat {
    mat := new(BTMat)
    X1 := pairing.NewG1().SetBytes(B1.b1)
    X2 := pairing.NewG1().SetBytes(B1.b2)
    Y1 := pairing.NewG2().SetBytes(B2.b1)
    Y2 := pairing.NewG2().SetBytes(B2.b2)

    mat.el11 = pairing.NewGT().Pair(X1, Y1).Bytes()
    mat.el12 = pairing.NewGT().Pair(X1, Y2).Bytes()
    mat.el21 = pairing.NewGT().Pair(X2, Y1).Bytes()
    mat.el22 = pairing.NewGT().Pair(X2, Y2).Bytes()

    return mat
}

/*
 IotaHat: AT -> BT
 Here, the mapping is occuring from G2 -> B2^4
 */
func IotaHat(pairing *pbc.Pairing, Z *pbc.Element, sigma *CommonReferenceString) *BTMat {
    // Element from G2, first convert to B1 and B2
    // then map element into BT^4
    B1 := IotaPrime1(pairing, pairing.NewZr().SetInt32(1), sigma)
    B2 := Iota2(pairing, Z)

    // Map into BT
    mat := FMap(pairing, B1, B2)

    return mat
}

/*
 IotaHat2: AT -> BT
 Here, the mapping is occuring from G1 -> B1^4
 */
func IotaHat2(pairing *pbc.Pairing, Z *pbc.Element, sigma *CommonReferenceString) *BTMat {
    // Element from G1, first convert to B1 and B2
    // then map element into BT^4
    B1 := Iota1(pairing, Z)
    B2 := IotaPrime2(pairing, pairing.NewZr().SetInt32(1), sigma)

    // Map into BT
    mat := FMap(pairing, B1, B2)

    return mat
}

/*
 TODO: RhoHat - might not be possible since this function needs the inverse e^1
 */
func RhoHat() {

}

/*
 * Creates a mapping between elements in G1 and maps them
 * to elements in B1
 *
 * Iota1: G1 -> B1
 * Pairing: the pairing in the PBC lib described in CRS
 * Element: The element from G1 that is to be mapped to B1
 */
func Iota1(pairing *pbc.Pairing, el *pbc.Element) *BPair {
    pair := new(BPair)
    pair.b1 = pairing.NewG1().Set0().Bytes()
    pair.b2 = el.Bytes()
    return pair

}

/*
 * Takes an element in B1 which are a pair of elements in G1 (g_1, g_2)
 * and maps them back to G1
 *
 * Rho1: B1 -> G1
 * pairing: The pairing from the pbc library described in the CRS
 * BPair: the element in B1
 * Returns: element in G1
 */
func Rho1(pairing *pbc.Pairing, pair *BPair, alpha *pbc.Element) *pbc.Element {
    Z1 := pairing.NewG1().SetBytes(pair.b1)
    Z2 := pairing.NewG1().SetBytes(pair.b2)
    tmp := pairing.NewG1().MulZn(Z1, alpha)
    return pairing.NewG1().Sub(Z2, tmp)
}



/*
 * Creates a mapping between elements in G2 and maps them
 * to elements in B2
 *
 * Iota1: G2 -> B2
 * Pairing: the pairing in the PBC lib described in CRS
 * Element: The element from G2 that is to be mapped to B2
 */
func Iota2(pairing *pbc.Pairing, el *pbc.Element) *BPair {
    pair := new(BPair)
    pair.b1 = pairing.NewG2().Set0().Bytes()
    pair.b2 = el.Bytes()
    return pair
}

/*
 * Takes an element in B2 which are a pair of elements in G2 (g_1, g_2)
 * and maps them back to G2
 *
 * Rho1: B2 -> G2
 * pairing: The pairing from the pbc library described in the CRS
 * BPair: the element in B2
 * Returns: element in G2
 */
func Rho2(pairing *pbc.Pairing, pair *BPair, alpha *pbc.Element) *pbc.Element {
    Z1 := pairing.NewG2().SetBytes(pair.b1)
    Z2 := pairing.NewG2().SetBytes(pair.b2)
    tmp := pairing.NewG2().MulZn(Z1, alpha)
    return pairing.NewG2().Sub(Z2, tmp)
}

/*
 * IotaT: GT -> BT
 * BT is in GT^4
 */
func IotaT(pairing *pbc.Pairing, el *pbc.Element) *BTMat{
    mat := new(BTMat)
    mat.el11 = pairing.NewGT().Set1().Bytes() // identity
    mat.el12 = pairing.NewGT().Set1().Bytes() // identity
    mat.el21 = pairing.NewGT().Set1().Bytes() // identity
    mat.el22 = el.Bytes()

    return mat
}


/* IotaPrime1: Zp -> B1
 * IotaPrimt1(z) = zu
 */
func IotaPrime1(pairing *pbc.Pairing, z *pbc.Element, sigma *CommonReferenceString) *BPair {
    pair := new(BPair)
    u1 := pairing.NewG1().SetBytes(sigma.u.u1)
    u2 := pairing.NewG1().SetBytes(sigma.u.u2)
    pair.b1 = pairing.NewG1().MulZn(u1, z).Bytes()
    pair.b2 = pairing.NewG1().MulZn(u2, z).Bytes()
    return pair
}

/*
 RhoPrime1: B1 -> Zp
    = (z2 - alpha * z1)
    // TODO: Convert zP back into z in Zp: z = z*P(P^-1)
 */
func RhoPrime1(pairing *pbc.Pairing, pair *BPair, alpha *pbc.Element) *pbc.Element{
    b1 := pairing.NewG1().SetBytes(pair.b1)
    b2 := pairing.NewG1().SetBytes(pair.b2)

    b2prime := pairing.NewG1().MulZn(b1, alpha)
    zprime := pairing.NewG1().Sub(b2, b2prime)
    // TODO: zprime should be in the group Zp but it is currently in G1 (need to convert)

    return zprime
}

/* IotaPrime2: Zp -> B2
 * IotaPrimt1(z) = zu
 */
func IotaPrime2(pairing *pbc.Pairing, z *pbc.Element, sigma *CommonReferenceString) *BPair {
    pair := new(BPair)
    u1 := pairing.NewG2().SetBytes(sigma.v.u1)
    u2 := pairing.NewG2().SetBytes(sigma.v.u2)
    pair.b1 = pairing.NewG2().MulZn(u1, z).Bytes()
    pair.b2 = pairing.NewG2().MulZn(u2, z).Bytes()
    return pair
}

/*
 RhoPrime2: B2 -> Zp
    = (z2 - alpha * z1)
    // TODO: Convert zP back into z in Zp: z = z*P(P^-1)
 */
func RhoPrime2(pairing *pbc.Pairing, pair *BPair, alpha *pbc.Element) *pbc.Element{
    b1 := pairing.NewG2().SetBytes(pair.b1)
    b2 := pairing.NewG2().SetBytes(pair.b2)

    b2prime := pairing.NewG2().MulZn(b1, alpha)
    zprime := pairing.NewG2().Sub(b2, b2prime)
    // TODO: zprime should be in the group Zp but it is currently in G1 (need to convert)

    return zprime
}