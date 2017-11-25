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
func CreateCommonReferenceString(sharedParams *SharedParams) *Sigma {
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	g1 := pairing.NewG1().SetBytes(sharedParams.G1)
	g2 := pairing.NewG2().SetBytes(sharedParams.G2)
	sigma := new(Sigma)

	// Create commit keys for u1 and u2 on G1
	alpha := pairing.NewZr().Rand()

	p1 := g1.Bytes()
	q1 := pairing.NewG1().MulZn(g1, alpha)

	t := pairing.NewZr().Rand()
	p2 := pairing.NewG1().MulZn(g1, t)
	q2 := pairing.NewG1().MulZn(q1, t)

	sigma.U = []CommitmentKey{
		CommitmentKey{p1, q1.Bytes()},
		CommitmentKey{p2.Bytes(), q2.Bytes()},
	}

	// Create commit keys v1 and v2 on G2
	alpha2 := pairing.NewZr().Rand()

	p12 := g2.Bytes()
	q12 := pairing.NewG1().MulZn(g2, alpha2)

	t2 := pairing.NewZr().Rand()
	p22 := pairing.NewG1().MulZn(g2, t2)
	q22 := pairing.NewG1().MulZn(q1, t2)

	sigma.V = []CommitmentKey{
		CommitmentKey{p12, q12.Bytes()},
		CommitmentKey{p22.Bytes(), q22.Bytes()},
	}

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