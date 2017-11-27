/*
 * Version 1.0 Prototype for benchmark
 * The main ocert scheme, it contains three protocl
 *  - Setup
 *  - GenECert
 *  - GenOCert
 * It contains the following helper functions
 *  - Get
 *  - GetSharedParams
 */

package ocert

import (
 	"fmt"
 	"crypto"
 	"crypto/rsa"
 	"crypto/rand"
 	"crypto/sha256"
)

// TODO delete
func Put(stub Wrapper, args [][]byte) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(string(args[0]), args[1])
	if err != nil {
		return nil, fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return []byte(args[1]), nil

}

// TODO delete
func Get(stub Wrapper, args [][]byte) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(string(args[0]))
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
var rsaPrivateKey *rsa.PrivateKey

func GetSharedParams(stub Wrapper, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting no arguments")
	}
	value, err := sharedParams.Bytes()
	if err != nil {
		return nil, err
	}
	return value, nil
}

/*
 * Setup is called by chaincode Init.
 * It generates 3 keypairs.
 *  1. Auditor's key pair (from rerandomization scheme)
 *  2. Key pair to generate ecert (from structure preserving scheme)
 *  3. Key pair to generate ocert (from RSA)
 * All public keys are stored in blockchain, while the private
 * keys are in memory. It returns the Auditor's keypair to the auditor
 */
func Setup(stub Wrapper, args [][]byte) ([]byte, error) {
	fmt.Println("Setup")
	if len(args) != 0 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting no arguments")
	}

	sharedParams = GenerateSharedParams()
	// Generate auditor's keypair
	PKa, SKa := EKeyGen(sharedParams)
	err := stub.PutState("auditor_pk", PKa.PK)
	KPa := new(AuditorKeypair)
	KPa.PK = PKa.PK
	KPa.SK = SKa.SK

	// Generate RSA keypair
	rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Generate structure preserving keypair
	VKei, SKei := SKeyGen(sharedParams)
	SSK = SKei
	SVKb, err := VKei.Bytes()
	if err != nil {
		return nil, err
	}
	err = stub.PutState("structure_preserving_vk", SVKb)
	if err != nil {
		return nil, err
	}

	// TODO ?PSetup

	// TODO generate certificate of CA

	// Return keypair to the auditor
	KPab, err := KPa.Bytes()
	if err != nil {
		return nil, err
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

/*
 * GenOCert is used to generate an ocert of a client
 * It takes a client's public key, a client's pseudonym and the 
 * proof of knowledge, and returns the ocert to the client 
 */
func GenOCert(stub Wrapper, args [][]byte) ([]byte, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting the PK, P, and PoK of client")
	}

	PKc := new(ClientPublicKey)
	PKc.PK = args[0]
	P := new(Pseudonym)
	err := P.SetBytes(args[1])
	if err != nil {
		return nil, err
	}

	// TODO verify proof of knowledge

	// TODO generate X.509 certificate
	msg, err := OCertSingedBytes(PKc, P)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(msg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}