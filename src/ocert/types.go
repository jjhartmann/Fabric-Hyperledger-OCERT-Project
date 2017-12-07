/*
 * All types used in ocert package, based on Ocert paper.
 */

package ocert

import (
 	"fmt"
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
// We assume client id is an element from G1 for simplicity.
type ClientID struct {
	ID []byte
}

/*
 * The public key of the client. This public key can be generated
 * by any scheme, but it should be an element in G2.
 */
type ClientPublicKey struct {
	PK []byte
}

/*
 * The pseudonym of a client. It is the signature generated by
 * rerandomization scheme E. It is an element in G1 * G1
 */
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

func (bp *BPair) Print() {
	fmt.Println("\t\t\t[BPair]")
	fmt.Printf("\t\t\t[b1]: ")
	fmt.Println(bp.b1)

	fmt.Printf("\t\t\t[b2]: ")
	fmt.Println(bp.b2)
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

// TODO: Why doesn't this work!
//func (tmp BPair) SubinG1(pairing *pbc.Pairing, l *BPair, r *BPair) *BPair{
//  ret := new(BPair)
//
//  // Convert to Groups
//  Bp1 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(l.b1),
//                            pairing.NewG1().SetBytes(r.b1))
//  Bp2 := pairing.NewG1().Sub(pairing.NewG1().SetBytes(l.b2),
//                             pairing.NewG1().SetBytes(r.b2))
//
//  ret.b1 = Bp1.Bytes()
//  ret.b2 = Bp2.Bytes()
//
//  return ret
//}
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
 * CommonReferenceString: the common reference string (CRS) for NIWI ProofOfEquation
 */
type CommonReferenceString struct {
	U []CommitmentKey // U[0].u1 holds generator G1
	V []CommitmentKey // V[0].u1 holds generator G2
	u CommitmentKey // For group G1 (used in Iota Prime)
	v CommitmentKey // For group G2
}

func (crs *CommonReferenceString) Print() {
	fmt.Println("\t[CommonReferenceString]")
	fmt.Printf("\t\t[U]: ")
	fmt.Println(crs.U)

	fmt.Printf("\t\t[V]: ")
	fmt.Println(crs.V)

	fmt.Printf("\t\t[u]: ")
	fmt.Println(crs.u)

	fmt.Printf("\t\t[v]: ")
	fmt.Println(crs.v)
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

type ProofOfEquation struct {
	Pi     []*BPair
	Theta  []*BPair
	c      []*BPair
	d      []*BPair
	cprime []*BPair
	dprime []*BPair
}

func (eq *ProofOfEquation) Print() {
	fmt.Println("\t[ProofOfEquation]")
	fmt.Println("\t\t[Pi]: ")
	for _, element := range eq.Pi {
		element.Print()
	}
	fmt.Println("")

	fmt.Println("\t\t[Theta]: ")
	for _, element := range eq.Theta {
		element.Print()
	}
	fmt.Println("")

	fmt.Println("\t\t[c]: ")
	for _, element := range eq.c {
		element.Print()
	}
	fmt.Println("")

	fmt.Println("\t\t[d]: ")
	for _, element := range eq.d {
		element.Print()
	}
	fmt.Println("")

	fmt.Println("\t\t[cprime]: ")
	for _, element := range eq.cprime {
		element.Print()
	}
	fmt.Println("")

	fmt.Println("\t\t[dprime]: ")
	for _, element := range eq.dprime {
		element.Print()
	}
	fmt.Println("")
}

type ProofOfKnowledge struct {
	Eq1 *ProofOfEquation
	Eq2 *ProofOfEquation
	Eq3 *ProofOfEquation
	Eq4 *ProofOfEquation
	Eq5 *ProofOfEquation
	sigma *CommonReferenceString
}

func (pi *ProofOfKnowledge) Print() {
	fmt.Println("[ProofOfKnowledge]-------------------")
	fmt.Printf("\t[Eq1]: **********************")
	pi.Eq1.Print()

	fmt.Printf("\t[Eq2]: **********************")
	pi.Eq2.Print()

	fmt.Printf("\t[Eq3]: **********************")
	pi.Eq3.Print()

	fmt.Printf("\t[Eq4]: **********************")
	pi.Eq4.Print()

	fmt.Printf("\t[Eq5]: **********************")
	pi.Eq5.Print()

	fmt.Printf("\t[sigma]: **********************")
	pi.sigma.Print()
	fmt.Println("-------------------")
	fmt.Println("")
}

/*
 * The proof generate proof of knowledge by using these variables
 * as witness
 */
type ProofVariables struct {
	P      *Pseudonym
	PKc    *ClientPublicKey
	PKa    *AuditorPublicKey
	E      *Ecert
	VK     *SVerificationKey
	Xc     []byte // This is the client private key
	RPrime []byte
}

/*
 * The proof takes these constant to validate that 5 equations hold
 */
type ProofConstants struct {
	// g1, g2 and e(g1, g2) are from sharedParams
	VK    *SVerificationKey // U, V, W1, W2 and Z
	PPrime *Pseudonym        // C' and D'
	Egz    []byte            // e(g1, Z)
	PKa    *AuditorPublicKey
	Egh    []byte            // e(G, H)
}

func (consts *ProofConstants) Print() {
	fmt.Println("[ProofConstants]-------------------")
	fmt.Printf("\t[VK]: ")
	fmt.Println(consts.VK)

	fmt.Printf("\t[PPrime]: ")
	fmt.Println(consts.PPrime)

	fmt.Printf("\t[Egs]: ")
	fmt.Println(consts.Egz)

	fmt.Printf("\t[PKa]: ")
	fmt.Println(consts.PKa)

	fmt.Printf("\t[Egh]: ")
	fmt.Println(consts.Egh)
	fmt.Println("-------------------")
	fmt.Println("")
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
