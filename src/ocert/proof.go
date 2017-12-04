/*
 * Based on Non-interactive Proof Systems for Bilinear Groups
 */

package ocert

import (
	"github.com/Nik-U/pbc"
	"reflect"
	_ "fmt"
)

/*
 * Set up the proof of knowledge, called by the client. It takes a system
 * of equations(e.g. pairing product equations and multi-scalar multiplication
 * equations) and outputs proof (e.g pi and theta ...)
 */
func PSetup(sharedParams *SharedParams, eqs *SystemOfEquations, vars *ProofVariables) *ProofOfKnowledge {
	pi := new(ProofOfKnowledge)

	// TODO setup proof of eq1
	// TODO setup proof of eq2
	// TODO setup proof of eq3
	// TODO setup proof of eq4
	// TODO setup proof of eq5

	return pi
}

/*
 * Validate the proof of knowledage, return true if all the equations
 * in the system hold.
 */
// TODO it may take extra information to prove
func PProve(sharedParams *SharedParams, pi *ProofOfKnowledge, consts *ProofConstants) bool {
	// TODO validate eq1
	// TODO validate eq2
	// TODO validate eq3
	// TODO validate eq4
	// TODO validate eq5

	return false
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
func ProveEquation1(pairing *pbc.Pairing, xc *pbc.Element, H *pbc.Element, PKc *pbc.Element, sigma *Sigma) *ProofOfEquation{
	proof := new(ProofOfEquation)

	// Create commitment in B1 for Xc
	cprime, _, R := CreateCommitmentPrimeOnG1(pairing, []*pbc.Element{xc}, sigma)

	// Create commitment in B2 for PKc
	d, _, S := CreateCommitmentOnG2(pairing, []*pbc.Element{PKc}, sigma)

	// Convert parameters to B groups
	Hi := Iota2(pairing, H)
	PKci := Iota2(pairing, PKc)
	neg1 := IotaPrime1(pairing, pairing.NewZr().SetInt32(-1), sigma)
	xci := IotaPrime1(pairing, xc, sigma)

	// Random samples from Zp
	if R.cols != 1 && R.rows != 1 {
		panic("Issues in conversion and creation of samples in Zp for R")
	}
	r := R.mat[0][0]
	gammaMat := NewRMatrix(pairing, 1, 1)
	gamma := gammaMat.mat[0][0]
	rgamma := pairing.NewZr().Mul(r, gamma)

	if S.rows != 1 && S.cols != 2 {
		panic("Issues in conversion and creation of samples in Zp for S")
	}

	/////////////////////////////////////
	// Pi: In G2

	// Multiply Scalar from Zn on B elements
	Hir := Hi.MulScalarInG2(pairing, r)           // r*ι_2(H)
	PKcirl := PKci.MulScalarInG2(pairing, rgamma) // r*gamma*ι_2(PKc)

	// Create Phi := (r*lambada*S - T)
	T := NewRMatrix(pairing, 2, 1)
	Srl := S.MulScalarZn(pairing, rgamma)  // r*gamma*S
	Ti := T.InvertMatrix()                 // T' invert
	Phi := Srl.ElementWiseSub(pairing, Ti) // S - T'

	// Multiple Phi by commitment keys
	Vphi := Phi.MulCommitmentKeysG2(pairing, sigma.V) // (r*gamma*S - T')V (commitment key in G2)
	if len(Vphi) > 1{
		panic("VPhi Should have len == 1")
	}

	// Construct Pi (Hir + PKcirl + Vphi)
	HPKcir := Hir.AddinG2(pairing, PKcirl)
	pi := HPKcir.AddinG2(pairing, Vphi[0])


	////////////////////////////////////////////
	// Theta: In G1
	_ = neg1
	_ = xci

	Si := S.InvertMatrix()                        // S Invert = S'
	Sneg := Si.MulBScalarinB1(pairing, *neg1)     // S'*ι'_1(-1)
	// +
	Sl := Si.MulScalarZn(pairing, gamma)    // S'*gamma
	Sxc := Sl.MulBScalarinB1(pairing, *xci) // S'*gamma*ι'_1(Xc)
	// +
	Tu := T.MulCommitmentKeysG1(pairing, []CommitmentKey{sigma.U[0]}) // Tu_1

	if len(Sneg) != len(Sxc) && len(Sxc) != len(Tu) {
		panic("All section lengths need to be equivalent")
	}

	// Construct theta
	theta := []*BPair{}

	for i := 0; i < len(Sneg); i++ {
		Snegxc := Sneg[i][0].AddinG1(pairing, Sxc[i][0])
		tmpB := Snegxc.AddinG1(pairing, Tu[i])
		theta = append(theta, tmpB)
	}

	// Collect elements
	proof.Gamma = gammaMat
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
 *	 C' is a Constant
 *
 * Proof:
 *    Pi     := R'*ι'_2(1) + R'*lambda*ι'_2(rprime)  + (R'*lambda*s - T')*v_1
 *    Theta  := s*ι_1(C) + s*lambda*ι_1(G) + TU
 */
func ProveEquation2(pairing *pbc.Pairing, rprime *pbc.Element, G *pbc.Element, C *pbc.Element, sigma *Sigma) *ProofOfEquation{
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
 * Verifiy Equation 1
 *
 */
func VerifyEquation1(pairing *pbc.Pairing, proof *ProofOfEquation, H *pbc.Element, tau *pbc.Element, sigma *Sigma) bool {

	// Construct LHS
	neg1 := IotaPrime1(pairing, pairing.NewZr().SetInt32(-1), sigma)
	Fid := FMap(pairing, neg1, proof.d[0])
	// +
	Hi := Iota2(pairing, H)
	FcH := FMap(pairing, proof.cprime[0], Hi)
	// +
	gamma := proof.Gamma.mat[0][0]
	dgamma := proof.d[0].MulScalarInG2(pairing, gamma)
	Fcd := FMap(pairing, proof.cprime[0], dgamma)
	// =
	tmp1 := Fid.AddinGT(pairing, FcH)
	LHS := tmp1.AddinGT(pairing, Fcd)


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
func VerifyEquation2(pairing *pbc.Pairing, proof *ProofOfEquation, G *pbc.Element, tau *pbc.Element, sigma *Sigma) bool {

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
 * Create Common refernce string sigma.
 * sigma = (u1, u2, v1, v2)
 * u1 and u1 are in B1
 * v1 and v2 are in B2
 *
 * u1 = (O, P)
 * u2 = t * u1
 */
func CreateCommonReferenceString(sharedParams *SharedParams, alpha *pbc.Element) *Sigma {
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

	// Proof should use different generators then what is stored in the
	// params. Florian: Using the same generators could case security issues
	// due to the discrete logrihtm problem
	// Since the groups are cyclic, this should not matter.
	g1 := pairing.NewG1().Rand()
	g2 := pairing.NewG2().Rand()
	sigma := new(Sigma)

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
func CreateCommitmentOnG1(pairing *pbc.Pairing, chi []*pbc.Element, sigma *Sigma) ([]*BPair, []*BPair, *RMatrix){

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
func CreateCommitmentPrimeOnG1(pairing *pbc.Pairing, x []*pbc.Element, sigma *Sigma) ([]*BPair, []*BPair, *RMatrix){

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
func CreateCommitmentOnG2(pairing *pbc.Pairing, Y []*pbc.Element, sigma *Sigma) ([]*BPair, []*BPair, *RMatrix){

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
func CreateCommitmentPrimeOnG2(pairing *pbc.Pairing, y []*pbc.Element, sigma *Sigma) ([]*BPair, []*BPair, *RMatrix){

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
func IotaHat(pairing *pbc.Pairing, Z *pbc.Element, sigma *Sigma) *BTMat {
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
func IotaHat2(pairing *pbc.Pairing, Z *pbc.Element, sigma *Sigma) *BTMat {
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
func IotaPrime1(pairing *pbc.Pairing, z *pbc.Element, sigma *Sigma) *BPair {
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
func IotaPrime2(pairing *pbc.Pairing, z *pbc.Element, sigma *Sigma) *BPair {
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