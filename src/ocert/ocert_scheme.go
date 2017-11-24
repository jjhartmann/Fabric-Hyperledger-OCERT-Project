/*
 * The main ocert scheme, it contains three protocl
 *  - Setup
 *  - GenECert
 *  - GenOCert
 */

package ocert

import (
 	"fmt"
)

// TODO delete
func Put(stub Wrapper, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return nil, fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return []byte(args[1]), nil

}

// TODO delete
func Get(stub Wrapper, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return nil, fmt.Errorf("Asset not found: %s", args[0])
	}
	return value, nil
}

/*
 * The private key used in structure preserving scheme should keep in memory,
 * not publicly on blockchain.
 */
var sharedParams *SharedParams
var SSK *SSigningKey

/*
 * Setup is called by chaincode Init.
 * It takes the auditor id(public key?) as input and stored in blockchain and 
 * generates 3 keypairs.
 *  1. Auditor's key pair (from rerandomization scheme)
 *  2. Key pair to generate ecert (from structure preserving scheme)
 *  3. Key pair to generate ocert (from RSA)
 * All public keys are stored in blockchain, while the private
 * keys are in memory. It returns the Auditor's keypair to the auditor
 */
func Setup(stub Wrapper, args [][]byte) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting the id of an auditor")
	}

	// TODO verify auditor
	err := stub.PutState("auditor_id", args[0])
	if err != nil {
		return nil, fmt.Errorf("Failed to set auditor_id")
	}

	sharedParams = GenerateSharedParams()
	// Generate auditor's keypair
	PKa, SKa := EKeyGen(sharedParams)
	err = stub.PutState("auditor_pk", PKa.PK)
	KPa := new(AuditorKeypair)
	KPa.PK = PKa.PK
	KPa.SK = SKa.SK

	// TODO Generate RSA keypair

	// Generate structure preserving keypair
	VKei, SKei := SKeyGen(sharedParams)
	SSK = SKei
	SVKb, err := VKei.Bytes()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate structure preserving keypair")
	}
	err = stub.PutState("structure_preserving_vk", SVKb)
	if err != nil {
		return nil, fmt.Errorf("Failed to set structure_preserving_vk")
	}

	// TODO ?PSetup

	// Return keypair to the auditor
	KPab, err := KPa.Bytes()
	if err != nil {
		return nil, fmt.Errorf("Failed to get auditor's keypair")
	}
	return KPab, nil
}

/*
 * GenECert is used to generate an ecert of a client
 * It takes the client id and the client's public key, and returns
 * psudonym P and ecert to the client.
 */
func GenECert(stub Wrapper, args [][]byte) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting the id and the PK of client")
	}

	IDc := new(ClientID)
	IDc.ID = args[0]
	PKc := new(ClientPublicKey)
	PKc.PK = args[1]

	// Generate pseudonym P
	valuePKa, err := stub.GetState("auditor_pk")
	if err != nil {
		return nil, fmt.Errorf("Failed to get auditor_pk")
	}
	if valuePKa == nil {
		return nil, fmt.Errorf("Asset not found: auditor_pk")
	}
	PKa := new(AuditorPublicKey)
	PKa.PK = valuePKa

	P := EEnc(sharedParams, PKa, IDc)

	// Generate ecert
	ecert := SSign(sharedParams, SSK, P, PKc)

	reply := new(GenECertReply)
	reply.P, err = P.Bytes()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate pesudonym")
	}
	reply.ecert, err = ecert.Bytes()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate ecert")
	}
	replyBytes, err := reply.Bytes()
		if err != nil {
		return nil, fmt.Errorf("Failed to generate GenECertReply")
	}
	return replyBytes, nil
}

func GenOCert(stub Wrapper, args [][]byte) ([]byte, error) {
	return nil, nil
}