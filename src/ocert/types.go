/*
 * All types used in ocert package, based on Ocert paper.
 */

package ocert

import (
 	"fmt"
	// "github.com/Nik-U/pbc"
 	"encoding/json"
 	"bytes"
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

func (P *Pseudonym) Equals(P2 *Pseudonym) bool {
	return  bytes.Equal(P.C, P2.C) &&
		bytes.Equal(P.D, P2.D)
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

func (VK *SVerificationKey) Equals(VK2 *SVerificationKey) bool {
	return  bytes.Equal(VK.U, VK2.U) &&
		bytes.Equal(VK.V, VK2.V) &&
		bytes.Equal(VK.W1, VK2.W1) &&
		bytes.Equal(VK.W2, VK2.W2) &&
		bytes.Equal(VK.Z, VK2.Z)
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

func (bp *BPair) Equals(bp2 *BPair) bool {
	return  bytes.Equal(bp.b1, bp2.b1) &&
		bytes.Equal(bp.b2, bp2.b2)
}

func (bp *BPair) Bytes() ([]byte, error) {
	template := struct {
		B1 []byte
		B2 []byte
	} {
		bp.b1,
		bp.b2,
	}

	msg, err := json.Marshal(template)
	return msg, err
}

func (bp *BPair) SetBytes(msg []byte) error {
	template := new(struct {
		B1 []byte
		B2 []byte
	})

	err := json.Unmarshal(msg, template)
	if err != nil {
		return err
	}

	bp.b1 = template.B1
	bp.b2 = template.B2

	return nil
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
	fmt.Println("\t\t[U]: ")
	for _, element := range crs.U {
		(&element).Print()
	}
	fmt.Println("")

	fmt.Println("\t\t[V]: ")
	for _, element := range crs.U {
		(&element).Print()
	}
	fmt.Println("")

	fmt.Printf("\t\t[u]: ")
	(&crs.u).Print()

	fmt.Printf("\t\t[v]: ")
	(&crs.v).Print()
}

func (crs *CommonReferenceString) Equals(crs2 *CommonReferenceString) bool {
	if len(crs.U) != len(crs2.U) {
		return false
	}

	if len(crs.V) != len(crs2.V) {
		return false
	}

	for i, _ := range crs.U {
		if !(&crs.U[i]).Equals(&crs2.U[i]) {
			return false
		}
	}

	for i, _ := range crs.V {
		if !(&crs.V[i]).Equals(&crs2.V[i]) {
			return false
		}
	}
	return  (&crs.u).Equals(&crs2.u) &&
		(&crs.v).Equals(&crs2.v)
}

func (crs *CommonReferenceString) Bytes() ([]byte, error) {
	var BigU [][]byte
	BigU = make([][]byte, len(crs.U), len(crs.U))
	for i, _ := range crs.U {
		ckBytes, err := (&crs.U[i]).Bytes()
		if err != nil {
			return nil, err
		}
		BigU[i] = ckBytes
	}

	var BigV [][]byte
	BigV = make([][]byte, len(crs.V), len(crs.V))
	for i, _ := range crs.V {
		ckBytes, err := (&crs.V[i]).Bytes()
		if err != nil {
			return nil, err
		}
		BigV[i] = ckBytes
	}

	SmallU, err := (&crs.u).Bytes()
	if err != nil {
		return nil, err
	}

	SmallV, err := (&crs.v).Bytes()
	if err != nil {
		return nil, err
	}

	template := struct {
		BigU   [][]byte
		BigV   [][]byte
		SmallU []byte
		SmallV []byte
	} {
		BigU,
		BigV,
		SmallU,
		SmallV,
	}

	msg, err := json.Marshal(template)
	return msg, err
}

func (crs *CommonReferenceString) SetBytes(msg []byte) error {
	template := new(struct {
		BigU   [][]byte
		BigV   [][]byte
		SmallU []byte
		SmallV []byte
	})

	err := json.Unmarshal(msg, template)
	if err != nil {
		return err
	}

	var U []CommitmentKey
	U = make([]CommitmentKey, len(template.BigU), len(template.BigU))
	for i, _ := range template.BigU {
		ck := new(CommitmentKey)
		err := ck.SetBytes(template.BigU[i])
		if err != nil {
			return err
		}
		U[i] = *ck
	}

	var V []CommitmentKey
	V = make([]CommitmentKey, len(template.BigV), len(template.BigV))
	for i, _ := range template.BigV {
		ck := new(CommitmentKey)
		err := ck.SetBytes(template.BigV[i])
		if err != nil {
			return err
		}
		V[i] = *ck
	}
	
	u := new(CommitmentKey)
	err = u.SetBytes(template.SmallU)
	if err != nil {
		return err
	}

	v := new(CommitmentKey)
	err = v.SetBytes(template.SmallV)
	if err != nil {
		return err
	}

	crs.U = U
	crs.V = V
	crs.u = *u
	crs.v = *v

	return nil
}

type CommitmentKey struct {
	u1 []byte
	u2 []byte
}

func (ck *CommitmentKey) Print() {
	fmt.Println("\t\t\t[CommitmentKey]")
	fmt.Printf("\t\t\t[u1]: ")
	fmt.Println(ck.u1)

	fmt.Printf("\t\t\t[u2]: ")
	fmt.Println(ck.u2)
}

func (ck *CommitmentKey) Equals(ck2 *CommitmentKey) bool {
	return  bytes.Equal(ck.u1, ck2.u1) &&
		bytes.Equal(ck.u2, ck2.u2)
}

func (ck *CommitmentKey) Bytes() ([]byte, error) {
	template := struct {
		U1 []byte
		U2 []byte
	} {
		ck.u1,
		ck.u2,
	}

	msg, err := json.Marshal(template)
	return msg, err
}

func (ck *CommitmentKey) SetBytes(msg []byte) error {
	template := new(struct {
		U1 []byte
		U2 []byte
	})

	err := json.Unmarshal(msg, template)
	if err != nil {
		return err
	}

	ck.u1 = template.U1
	ck.u2 = template.U2

	return nil
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

func (eq *ProofOfEquation) Equals(eq2 *ProofOfEquation) bool {
	if len(eq.Pi) != len(eq2.Pi) {
		return false
	}

	if len(eq.Theta) != len(eq2.Theta) {
		return false
	}

	if len(eq.c) != len(eq2.c) {
		return false
	}

	if len(eq.d) != len(eq2.d) {
		return false
	}

	if len(eq.cprime) != len(eq2.cprime) {
		return false
	}

	if len(eq.dprime) != len(eq2.dprime) {
		return false
	}

	for i, _ := range eq.Pi {
		if !eq.Pi[i].Equals(eq2.Pi[i]) {
			return false
		}
	}

	for i, _ := range eq.Theta {
		if !eq.Theta[i].Equals(eq2.Theta[i]) {
			return false
		}
	}

	for i, _ := range eq.c {
		if !eq.c[i].Equals(eq2.c[i]) {
			return false
		}
	}

	for i, _ := range eq.d {
		if !eq.d[i].Equals(eq2.d[i]) {
			return false
		}
	}

	for i, _ := range eq.cprime {
		if !eq.cprime[i].Equals(eq2.cprime[i]) {
			return false
		}
	}

	for i, _ := range eq.dprime {
		if !eq.dprime[i].Equals(eq2.dprime[i]) {
			return false
		}
	}

	return true
}

func proofOfEquationBytesHelper(in []*BPair) ([][]byte, error) {
	var IN [][]byte
	IN = make([][]byte, len(in), len(in))
	for i, _ := range in {
		bpBytes, err := in[i].Bytes()
		if err != nil {
			return nil, err
		}
		IN[i] = bpBytes
	}
	return IN, nil
}

func (eq *ProofOfEquation) Bytes() ([]byte, error) {
	Pi, err := proofOfEquationBytesHelper(eq.Pi)
	if err != nil {
		return nil, err
	}

	Theta, err := proofOfEquationBytesHelper(eq.Theta)
	if err != nil {
		return nil, err
	}

	c, err := proofOfEquationBytesHelper(eq.c)
	if err != nil {
		return nil, err
	}

	d, err := proofOfEquationBytesHelper(eq.d)
	if err != nil {
		return nil, err
	}

	cprime, err := proofOfEquationBytesHelper(eq.cprime)
	if err != nil {
		return nil, err
	}

	dprime, err := proofOfEquationBytesHelper(eq.dprime)
	if err != nil {
		return nil, err
	}

	template := struct {
		Pi     [][]byte
		Theta  [][]byte
		C      [][]byte
		D      [][]byte
		Cprime [][]byte
		Dprime [][]byte
	} {
		Pi,
		Theta,
		c,
		d,
		cprime,
		dprime,
	}

	msg, err := json.Marshal(template)
	return msg, err
}

func proofOfEquationSetBytesHelper(in [][]byte) ([]*BPair, error) {
	var IN []*BPair
	IN = make([]*BPair, len(in), len(in))
	for i, _ := range in {
		bp := new(BPair)
		err := bp.SetBytes(in[i])
		if err != nil {
			return nil, err
		}
		IN[i] = bp
	}
	return IN, nil
}

func (eq *ProofOfEquation) SetBytes(msg []byte) error {
	template := new(struct {
		Pi     [][]byte
		Theta  [][]byte
		C      [][]byte
		D      [][]byte
		Cprime [][]byte
		Dprime [][]byte
	})

	err := json.Unmarshal(msg, template)
	if err != nil {
		return err
	}

	Pi, err := proofOfEquationSetBytesHelper(template.Pi)
	if err != nil {
		return err
	}

	Theta, err := proofOfEquationSetBytesHelper(template.Theta)
	if err != nil {
		return err
	}

	c, err := proofOfEquationSetBytesHelper(template.C)
	if err != nil {
		return err
	}

	d, err := proofOfEquationSetBytesHelper(template.D)
	if err != nil {
		return err
	}

	cprime, err := proofOfEquationSetBytesHelper(template.Cprime)
	if err != nil {
		return err
	}

	dprime, err := proofOfEquationSetBytesHelper(template.Dprime)
	if err != nil {
		return err
	}

	eq.Pi = Pi
	eq.Theta = Theta
	eq.c = c
	eq.d = d
	eq.cprime = cprime
	eq.dprime = dprime

	return nil
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

func (pi *ProofOfKnowledge) Equals(pi2 *ProofOfKnowledge) bool {
	return pi.Eq1.Equals(pi2.Eq1) &&
		pi.Eq2.Equals(pi2.Eq2) &&
		pi.Eq3.Equals(pi2.Eq3) &&
		pi.Eq4.Equals(pi2.Eq4) &&
		pi.Eq5.Equals(pi2.Eq5) &&
		pi.sigma.Equals(pi2.sigma)
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
	VK     *SVerificationKey // U, V, W1, W2 and Z
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

func (consts *ProofConstants) Equals(consts2 *ProofConstants) bool {
	return consts.VK.Equals(consts2.VK) &&
		consts.PPrime.Equals(consts2.PPrime) &&
		bytes.Equal(consts.Egz, consts2.Egz) &&
		bytes.Equal(consts.PKa.PK, consts2.PKa.PK) &&
		bytes.Equal(consts.Egh, consts2.Egh)
}


func (consts *ProofConstants) Bytes() ([]byte, error) {
	VKbytes, err := consts.VK.Bytes()
	if err != nil {
		return nil, err
	}

	PPrimeBytes, err := consts.PPrime.Bytes()
	if err != nil {
		return nil, err
	}

	template := struct {
		VK     []byte
		PPrime []byte
		Egz    []byte
		PKa    []byte
		Egh    []byte 
	} {
		VKbytes,
		PPrimeBytes,
		consts.Egz,
		consts.PKa.PK,
		consts.Egh,
	}

	msg, err := json.Marshal(template)
	return msg, err
}

func (consts *ProofConstants) SetBytes(msg []byte) error {
	template := new(struct {
		VK     []byte
		PPrime []byte
		Egz    []byte
		PKa    []byte
		Egh    []byte 
	})

	err := json.Unmarshal(msg, template)
	if err != nil {
		return err
	}

	VK := new(SVerificationKey)
	err = VK.SetBytes(template.VK)
	if err != nil {
		return err
	}
	consts.VK = VK

	PPrime := new(Pseudonym)
	err = PPrime.SetBytes(template.PPrime)
	if err != nil {
		return err
	}
	consts.PPrime = PPrime

	consts.Egz = template.Egz

	PKa := new(AuditorPublicKey)
	PKa.PK = template.PKa
	consts.PKa = PKa

	consts.Egh = template.Egh

	return nil
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
