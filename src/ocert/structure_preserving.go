/*
 * Based on "Optimal Structure-Preserving Signatures in Asymmetric Bilinear Groups"
 * by Abe and et al. The key size depends on Ocert paper.
 */

package ocert

import (
)

/*
 * Generate key pair used by orgnization i.
 * SVerificationKey VK is used as ecert verification key VK_e,i for each
 * organization i
 * SSigningKey SK_e,i is used to sign client Pseudonym and ClientPublicKey
 * VK = (U, V, W1, W2, Z) in G1 * G2^4
 * SK = (u, v, w1, w2, z) where u, v, w1, w2, z is randomly picked from group of 
 * units modulo p
 * s.t. U = u * g1, V = v * g2, W1 = w1 * g2, W2 = w2 * g2 and Z = z * g2, where
 * g1 is the generator of group G1, and g2 is the generator of group G2
 */
// TODO It should take some security parameter
func SKeyGen() (*SVerificationKey, *SSigningKey) {
	VK := new(SVerificationKey)
	SK := new(SSigningKey)

	// TODO
	// How can we use PBC library to get generator of each group
	// and order of each element

	return VK, SK
}

/*
 * Signing the pseudonym and public key of a client by the signing key SK_e,i
 * of an organziation i. The output of signing procedure is the ecert.
 * The ecert has format (R, S, T) in G1 * G1 * G2, where g1, g2 are generators
 * of G1 and G2, r is randomly picked from group of units modulo p,
 * SKei = (u, v, w1, w2, z), P = (C, D)
 * R = r * g1
 * S = (z - r * v) * g1 + (-w1) * C + (-w2) * D
 * T = (1 / 6) * (g2 + (-u) * PKc)
 */
func SSign(SKei *SSigningKey, P *Pseudonym, PKc *ClientPublicKey) *Ecert {
	var ecert *Ecert = new(Ecert)

	// Generate R, S, T

	return ecert
}

/*
 * Verifying the signature is signed by the client, and returns a boolean, where
 * g1 and g2 are generators of group G1 and G2 respectively
 * VKei = (U, V, W1, W2, Z),
 * P = (C, D)
 * ecert = (R, S, T),
 * and to verify, test
 * e(R, V) * e(S, g2) * e(C, W1) * e(D, W2) = e(g1, Z) and
 * e(R, T) * e(U, PKc) = e(g1, g2), where e is the pairing operation
 */
func SVerify(VKei *SVerificationKey, P *Pseudonym, PKc *ClientPublicKey, ecert *Ecert) bool {
	return false
}