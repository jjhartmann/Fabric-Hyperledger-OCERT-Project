/*
 * Based on elGamal encryption on ECC
 */

package ocert

import (
	"github.com/Nik-U/pbc"
)

/*
 * Generate the pair of public key and secret key used by auditor.
 */
func EKeyGen(sharedParams *SharedParams) (*AuditorPublicKey, *AuditorSecretKey) {
	PKa := new(AuditorPublicKey)
	SKa := new(AuditorSecretKey)

	// Generators g1 & g2 are generated from the groups G1 & G2
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	g1 := pairing.NewG1().SetBytes(sharedParams.G1)
	g2 := pairing.NewG2().SetBytes(sharedParams.G2)


	xa :=pairing.NewZr().Rand()
  SKa.SK=xa.Bytes()

	//produce the public key
	PublicKa :=pairing.NewG1().MulZn(g1,xa)
  PKa.PK=PublicKa.Bytes()

	return PKa, SKa
}

/*
 * Encrypt the client id based on the auditor's public key, the
 * result is the pseudonym of a client, where pseudonym of a client
 * has form (C, D), where both C and D are in G1
 */
func EEnc(sharedParams *SharedParams, PKa *AuditorPublicKey, id *ClientID) *Pseudonym {
	P := new(Pseudonym)

	return P
}

/*
 * Decrypt the client real identiy based on the pseudonym of a client
 */
func EDec(sharedParams *SharedParams, SKa *AuditorSecretKey, P *Pseudonym) *ClientID {
	id := new(ClientID)
	return id
}

/*
 * Rerandomize the client's pseudonym. Given a pseudonym P = (C, D)
 * of a client, this scheme can rerandomize it to a new pseudonym
 * P' = (C', D'), where P' is also in G1 * G1.
 */
// TODO This function may also return the extra information used in validation
func ERerand(sharedParams *SharedParams, P *Pseudonym) *Pseudonym {
	// TODO rerandomize P
	return P
}

/*
 * Given two pseudonyms P and P', validate whether P' is rerandomized
 * from P
 */
// TODO This function may take extra information to do verification
func ERerandVerify(sharedParams *SharedParams, P *Pseudonym, PPrime *Pseudonym) bool {
	return false
}
