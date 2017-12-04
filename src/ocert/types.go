/*
 * All types used in ocert package, based on Ocert paper.
 */

package ocert

import (
	// "github.com/Nik-U/pbc"
 	"encoding/json"
 	// "bytes"
	"github.com/Nik-U/pbc"
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

func (s *SharedParams) Bytes() ([]byte, error) {
	msg, err := json.Marshal(s)
	return msg, err
}

func (s *SharedParams) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, s)
	return err
}

/*****************************************************************/

/*
 * The public key of the auditor. This public key is generated by
 * the rerandomization scheme E, and it is an element in G1.
 */
type AuditorPublicKey struct {
	PK []byte
}

func (k *AuditorPublicKey) Bytes() ([]byte, error) {
	msg, err := json.Marshal(k)
	return msg, err
}

func (k *AuditorPublicKey) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, k)
	return err
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


////////////////////////////////////////////////////////////////////////////////////
/*
 * Holds a pair of elements form G1 or G2 in the group B1 or B2
 * Iota: G -> B
 */
type BPair struct {
	b1 []byte
	b2 []byte
}

// BPair Helper function to add elements
func (l BPair) AddinG1(pairing *pbc.Pairing, r *BPair) *BPair{
  ret := new(BPair)

	// Convert to Groups
  lb1 := pairing.NewG1().SetBytes(l.b1)
  lb2 := pairing.NewG1().SetBytes(l.b2)
  rb1 := pairing.NewG1().SetBytes(r.b1)
  rb2 := pairing.NewG1().SetBytes(r.b2)

  ret.b1 = pairing.NewG1().Add(lb1, rb1).Bytes()
  ret.b2 = pairing.NewG1().Add(lb2, rb2).Bytes()

  return ret
}
func (l BPair) AddinG2(pairing *pbc.Pairing, r *BPair) *BPair{
  ret := new(BPair)

  // Convert to Groups
  lb1 := pairing.NewG2().SetBytes(l.b1)
  lb2 := pairing.NewG2().SetBytes(l.b2)
  rb1 := pairing.NewG2().SetBytes(r.b1)
  rb2 := pairing.NewG2().SetBytes(r.b2)

  ret.b1 = pairing.NewG2().Add(lb1, rb1).Bytes()
  ret.b2 = pairing.NewG2().Add(lb2, rb2).Bytes()

  return ret
}
func (l BPair) SubinG1(pairing *pbc.Pairing, r *BPair) *BPair{
  ret := new(BPair)

  // Convert to Groups
  Bp1 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(l.b1),
                            pairing.NewG1().SetBytes(l.b1))
  Bp2 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(l.b2),
                             pairing.NewG1().SetBytes(l.b2))

  ret.b1 = Bp1.Bytes()
  ret.b2 = Bp2.Bytes()

  return ret
}
func (l BPair) MulScalarInG1(pairing *pbc.Pairing, r *pbc.Element) *BPair {
  pair := new (BPair)

  lb1 := pairing.NewG1().SetBytes(l.b1)
  lb2 := pairing.NewG1().SetBytes(l.b2)

  pair.b1 = pairing.NewG1().MulZn(lb1, r).Bytes()
  pair.b2 = pairing.NewG1().MulZn(lb2, r).Bytes()

  return pair
}
func (l BPair) MulScalarInG2(pairing *pbc.Pairing, r *pbc.Element) *BPair {
  pair := new (BPair)

  lb1 := pairing.NewG2().SetBytes(l.b1)
  lb2 := pairing.NewG2().SetBytes(l.b2)

  pair.b1 = pairing.NewG2().MulZn(lb1, r).Bytes()
  pair.b2 = pairing.NewG2().MulZn(lb2, r).Bytes()

  return pair
}
////////////////////////////////////////////////////////////////////////////////////
// BTMAT
// Row x Col
type BTMat struct {
	el11 []byte
	el12 []byte
	el21 []byte
	el22 []byte
}
// BTMat Functions
func (btmat *BTMat) AddinGT(pairing *pbc.Pairing, rb *BTMat) *BTMat{
  BT := new(BTMat)

  BT.el11 = pairing.NewGT().Add(pairing.NewGT().SetBytes(btmat.el11),
    pairing.NewGT().SetBytes(rb.el11)).Bytes()
  BT.el12 = pairing.NewGT().Add(pairing.NewGT().SetBytes(btmat.el12),
    pairing.NewGT().SetBytes(rb.el12)).Bytes()
  BT.el21 = pairing.NewGT().Add(pairing.NewGT().SetBytes(btmat.el21),
    pairing.NewGT().SetBytes(rb.el21)).Bytes()
  BT.el22 = pairing.NewGT().Add(pairing.NewGT().SetBytes(btmat.el22),
    pairing.NewGT().SetBytes(rb.el22)).Bytes()

  return BT
}

////////////////////////////////////////////////////////////////////////////////////
/*
 * Sigma: the common reference string (CRS) for NIWI ProofOfEquation
 */
type Sigma struct {
	U []CommitmentKey // U[0].u1 holds generator G1
	V []CommitmentKey // V[0].u1 holds generator G2
	u CommitmentKey // For group G1 (used in Iota Prime)
	v CommitmentKey // For group G2
}

type CommitmentKey struct {
	u1 []byte
	u2 []byte
}
func (ck CommitmentKey) ConvertToBPair() *BPair {
  B := new(BPair)
  B.b1 = ck.u1
  B.b2 = ck.u2
  return B
}


////////////////////////////////////////////////////////////////////////////////////
// Proof Structs
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

type ProofOfEquation struct {
	Pi     []*BPair
	Theta  []*BPair
  c      []*BPair
  d      []*BPair
  cprime []*BPair
  dprime []*BPair
  Gamma  *RMatrix
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
 * Request to and reply from main scheme (chaincode)
 */

type GenECertRequest struct {
	IDc []byte
	PKc []byte
}

func (request *GenECertRequest) Bytes() ([]byte, error) {
	msg, err := json.Marshal(request)
	return msg, err
}

func (request *GenECertRequest) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, request)
	return err
}

type GenECertReply struct {
	P []byte
	Ecert []byte
}

func (reply *GenECertReply) Bytes() ([]byte, error) {
	msg, err := json.Marshal(reply)
	return msg, err
}

func (reply *GenECertReply) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, reply)
	return err
}

type GenOCertRequest struct {
	PKc []byte
	P []byte
	// TODO pi
}

func (request *GenOCertRequest) Bytes() ([]byte, error) {
	msg, err := json.Marshal(request)
	return msg, err
}

func (request *GenOCertRequest) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, request)
	return err
}

type GenOCertReply struct {
	Sig []byte
}

func (reply *GenOCertReply) Bytes() ([]byte, error) {
	msg, err := json.Marshal(reply)
	return msg, err
}

func (reply *GenOCertReply) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, reply)
	return err
}

type RSAPK struct {
	PK []byte
}

func (pk *RSAPK) Bytes() ([]byte, error) {
	msg, err := json.Marshal(pk)
	return msg, err
}

func (pk *RSAPK) SetBytes(msg []byte) error {
	err := json.Unmarshal(msg, pk)
	return err
}
