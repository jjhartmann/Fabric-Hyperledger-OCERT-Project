/*
 * Based on elGamal encryption on ECC
 */

package ocert

import (
)

/*
 * Generate the pair of public key and secret key used by auditor.
 */
// TODO it should take some security parameter
func EKeyGen(sharedParams *SharedParams) (*AuditorPublicKey, *AuditorSecretKey) {
	PKa := new(AuditorPublicKey)
	SKa := new(AuditorSecretKey)

	return PKa, SKa
}

/*
 * Encrypt the client id based on the auditor's public key, the 
 * result is the pseudonym of a client, where pseudonym of a client
 * has form (C, D), where both C and D are in G1
 */
func EEnc(PKa *AuditorPublicKey, id *ClientID) *Pseudonym {
	P := new(Pseudonym)

	return P
}

/*
 * Decrypt the client real identiy based on the pseudonym of a client
 */
func EDec(SKa *AuditorSecretKey, P *Pseudonym) *ClientID {
	id := new(ClientID)
	return id
}

/*
 * Rerandomize the client's pseudonym. Given a pseudonym P = (C, D)
 * of a client, this scheme can rerandomize it to a new pseudonym
 * P' = (C', D'), where P' is also in G1 * G1.
 */
// TODO This function may also return the extra information used in validation
func ERerand(P *Pseudonym) *Pseudonym {
	// TODO rerandomize P
	return P
}

/*
 * Given two pseudonyms P and P', validate whether P' is rerandomized
 * from P 
 */
// TODO This function may take extra information to do verification
func ERerandVerify(P *Pseudonym, PPrime *Pseudonym) bool {
	return false
}
