/*
 * Based on Non-interactive Proof Systems for Bilinear Groups
 */

package ocert

import (
)

/*
 * Set up the proof of knowledge, called by the client
 */
// TODO the parameters
// TODO extra information used in proof may return
func Setup() *ProofOfKnowledge {
	pi := new(ProofOfKnowledge)

	return pi
}

/*
 * Prove a pairing product equaction, return true if such equation holds.
 */
// TODO it may take extra information to prove
func ProvePairingProductEquation(eq *PairingProductEquation) bool {
	return false
}

/*
 * Prove a multi-scalar multiplication equation in G1,
 * return true if such equation holds.
 */
// TODO it may take extra information to prove
func ProveMultiScalarMultiplicationEquationG1(eq *MultiScalarMultiplicationEquationG1) bool {
	return false
}

/*
 * Prove a multi-scalar multiplication equation in G2,
 * return true if such equation holds.
 */
// TODO it may take extra information to prove
func ProveMultiScalarMultiplicationEquationG2(eq *MultiScalarMultiplicationEquationG2) bool {
	return false
}

/*
 * Validate the proof of knowledage, return true if all the equations
 * in the system are proved.
 */
// TODO it may take extra information to prove
func Prove(pi *ProofOfKnowledge) bool {
	return ProveMultiScalarMultiplicationEquationG2(pi.Eq1) && 
		ProveMultiScalarMultiplicationEquationG1(pi.Eq2) &&
		ProveMultiScalarMultiplicationEquationG1(pi.Eq3) &&
		ProvePairingProductEquation(pi.Eq4) &&
		ProvePairingProductEquation(pi.Eq5)
}