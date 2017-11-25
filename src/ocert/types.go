/*
 * All types used in ocert package, based on Ocert paper.
 */

package ocert

import (
	// "github.com/Nik-U/pbc"
 	"encoding/json"
 	// "bytes"
)

/*
 * Since all schemes are based on same G1, G2 and Gt, this structure
 * holds the parameters shared by all schemes. G1, G2 and Gt here
 * are generators.
 */
// TODO
// We should choose same generator for each group and use these in
// KeyGen and Setup
type SharedParams struct {
	Params string
	G1     []byte
	G2     []byte
}

/*****************************************************************/

// TODO Based on PBC library we choose, we need to modify following
// wrappers
// TODO My suggestion is not to use these wrappers, as this makes us
// life easier. (see structure_preserving.go)
// TODO We need to store all elements in bytes array

// A wrapper of an element in G1
type G1Element struct {
}

// A wrapper of an element in G2
type G2Element struct {
}

// A wrapper of an element in Gt
type GtElement struct {
}

/*****************************************************************/

/*
 * The public key of the auditor. This public key is generated by
 * the rerandomization scheme E, and it is an element in G1.
 */
type AuditorPublicKey struct {
	PK []byte
}

/*
 * The secret key of the auditor.
 */
 type AuditorSecretKey struct {
 	// TODO, the group depends on scheme
 	SK []byte
 }

/*
 * The keypair of the auditor
 */
type AuditorKeypair struct {
	PK []byte
	SK []byte
}

func (KPa *AuditorKeypair) Bytes() ([]byte, error) {
	msg, err := json.Marshal(KPa)
	return msg, err
}

func (KPa *AuditorKeypair) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, KPa)
	return err
}

/*****************************************************************/

/*
 * The identity of a client
 */
// TODO We may use the identity structure of fabric here, and hash
// it to some value in a group in our scheme.
type ClientID struct {
	ID []byte
}

/*
 * The public key of the client. This public key can be generated
 * by any scheme, but it should be an element in G2.
 * TODO (Is it okay for us to hash any key to an element in G2?)
 */
type ClientPublicKey struct {
	PK []byte
}

/*
 * The pseudonym of a client. It is the signature generated by
 * rerandomization scheme E. It is an element in G1 * G1
 */
// TODO do we need a differnt struct for a rerandomized pesudonym
type Pseudonym struct {
	C []byte
	D []byte
}

func (P *Pseudonym) Bytes() ([]byte, error) {
	msg, err := json.Marshal(P)
	return msg, err
}

func (P *Pseudonym) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, P)
	return err
}

/*****************************************************************/

/*
 * SVerificationKey and SSigningKey are the key pairs generated by
 * Structure-preserving scheme S.
 * Each organization i will have one pair.
 * SVerificationKey is used as ecert verification key VK_e,i for each
 * organization i
 * SSigningKey is used to sign client Pseudonym and ClientPublicKey
 */

/*
 * Based on structure-preserving scheme S. The verification key contains
 * 5 elements, U, V, W1, W2 and Z, where only U is an element in G1 and
 * the rest are elements in G2.
 */
type SVerificationKey struct {
	U  []byte
	V  []byte
	W1 []byte
	W2 []byte
	Z  []byte
}

func (VK *SVerificationKey) Bytes() ([]byte, error) {
	msg, err := json.Marshal(VK)
	return msg, err
}

func (VK *SVerificationKey) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, VK)
	return err
}

/*
 * The signing key contains the order of each element in the verification
 * key.
 */
type SSigningKey struct {
	U  []byte
	V  []byte
	W1 []byte
	W2 []byte
	Z  []byte
}

/*
 * Ecert is the signature generated by scheme S. It contains three elements
 * R, S and T, where R and S are in G1 and T is in G2.
 */
type Ecert struct {
	R []byte
	S []byte
	T []byte
}

func (ecert *Ecert) Bytes() ([]byte, error) {
	msg, err := json.Marshal(ecert)
	return msg, err
}

func (ecert *Ecert) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, ecert)
	return err
}

/*****************************************************************/

/*
 * Equations over groups with bilinear map. These equations are used in
 * non-interactive proof systems over bilinear groups.(We should verify
 * these equations)
 */

// TODO optimize the equation, we may not need generic equaction any more,
// we can just keep a record of the 5 specific equactions used in our
// system by defining type Eq1 struct {}, type Eq2 struct {}, ...
/*
 * Pairing product equation has form
 * e(G1_1, G2_1) * e(G1_2, G2_2) * ... * (G1_l, G2_l) = RHS,
 * where RHS is in Gt
 */
type PairingProductEquationPair struct {
	X []byte
	Y []byte
}

type PairingProductEquation struct {
	Length uint32
	LHS    []PairingProductEquationPair
	RHS    []byte
}

/*
 * Multi-scalar multiplication equation in G1 has form
 * a_1 * G1_1 + a_2 * G1_2 + ... + a_l * G1_l = RHS,
 * where RHS is in G1
 */
type MultiScalarMultiplicationEquationG1Pair struct {
	A []byte
	X []byte
}

type MultiScalarMultiplicationEquationG1 struct {
	Length uint32
	LHS    []MultiScalarMultiplicationEquationG1Pair
	RHS    []byte
}

/*
 * Multi-scalar multiplication equation in G2 has form
 * a_1 * G2_1 + a_2 * G2_2 + ... + a_l * G2_l = RHS,
 * where RHS is in G2
 */
type MultiScalarMultiplicationEquationG2Pair struct {
	A []byte
	Y []byte
}

type MultiScalarMultiplicationEquationG2 struct {
	Length uint32
	LHS    []MultiScalarMultiplicationEquationG2Pair
	RHS    []byte
}

/*
 * This is the system of equations that used in the non-interactive proof.
 * It is used to prove the knowledge of a client's public key, the knowledge
 * of a rerandomization of a client's pseduonym and the knowledge of a valid
 * signatrue over the client's public key and pseudonym
 * Eq1: x_C * g2 + (-1) PK_c = 0
 * Eq2: C + r' * g1 = C'
 * Eq3: D + r' PK_a = D'
 * Eq4: e(R, V) * e(S, H) * e(C, W1) * e(D, W2) = e(G, Z)
 * Eq5: e(R, T) * e(U, PK_c) = e(G, H)
 */
type SystemOfEquations struct {
	Eq1 *MultiScalarMultiplicationEquationG2
	Eq2 *MultiScalarMultiplicationEquationG1
	Eq3 *MultiScalarMultiplicationEquationG1
	Eq4 *PairingProductEquation
	Eq5 *PairingProductEquation
}

/*****************************************************************/
/*
 * This is a proof generated by giving a equation over biliner group
 */

/*
 * Holds a pair of elements form G1 or G2 in the group B1 or B2
 * Iota: G -> B
 */
type BPair struct {
	b1 []byte
	b2 []byte
}

/*
 * Sigma: the common reference string (CRS) for NIWI ProofOfEquation
 */
type Sigma struct {
	U []CommitmentKey
	V []CommitmentKey
}

type CommitmentKey struct {
	u1 []byte
	u2 []byte
}

type ProofString struct {
	U1 []byte
	U2 []byte
	V1 []byte
	V2 []byte
}

type C struct {
}

type D struct {
}

type Pi struct {
}

type Theta struct {
}

type ProofOfEquation struct {
	C      *C
	D      *D
	CPrime *C
	DPrime *D
	Pi     *Pi
	Theta  *Theta
}

type ProofOfKnowledge struct {
	Eq1 *ProofOfEquation
	Eq2 *ProofOfEquation
	Eq3 *ProofOfEquation
	Eq4 *ProofOfEquation
	Eq5 *ProofOfEquation
}

/*
 * The proof generate proof of knowledge by using these variables
 * as witness
 */
type ProofVariables struct {
	P      *Pseudonym
	PKc    *ClientPublicKey
	E      *Ecert
	Xc     []byte // This is the client private key
	RPrime []byte
}

/*
 * The proof takes these constant to validate that 5 equations hold
 */
type ProofConstants struct {
	// g1, g2 and e(g1, g2) are from sharedParams
	PKi    *SVerificationKey // U, V, W1, W2 and Z
	PPrime *Pseudonym        // C' and D'
	Egz    []byte            // e(g1, Z)
	PKa    *AuditorPublicKey
}

/*****************************************************************/
/*
 * Result from main scheme(chaincode)
 */

type GenECertReply struct {
	P []byte
	ecert []byte
}

func (reply *GenECertReply) Bytes() ([]byte, error) {
	msg, err := json.Marshal(reply)
	return msg, err
}

func (reply *GenECertReply) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, reply)
	return err
}