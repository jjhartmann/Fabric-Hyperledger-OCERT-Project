/*
 * Based on Non-interactive Proof Systems for Bilinear Groups
 */

package ocert

import "github.com/Nik-U/pbc"

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