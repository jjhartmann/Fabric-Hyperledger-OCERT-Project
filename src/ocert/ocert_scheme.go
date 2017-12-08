/*
 * Version 1.0 Prototype for benchmark
 * The main ocert scheme, it contains three protocl
 *  - Setup
 *  - GenECert
 *  - GenOCert
 * It contains the following helper functions
 *  - Get
 *  - GetSharedParams
 *  - GetAuditorKeypair
 */

package ocert

import (
    "fmt"
    "crypto"
    "crypto/rsa"
    "crypto/rand"
    "crypto/sha256"
    "crypto/x509"
    "math/big"
    "time"
    "os"
    "github.com/Nik-U/pbc"
)

/*
 * The private key used in structure preserving scheme should keep in memory,
 * not publicly on blockchain.
 */
var sharedParams *SharedParams
var sSigningKey *SSigningKey
var rsaPrivateKey *rsa.PrivateKey
var serialNumber *big.Int
var auditorKeypair []byte
var consts *ProofConstants

var verifyProofLog *os.File


func getSerialNumber() (*big.Int) {
    serialNumber.Add(serialNumber, big.NewInt(1))
    return serialNumber
}

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

func GetSharedParams(stub Wrapper, args [][]byte) ([]byte, error) {
    if len(args) != 0 {
        return nil, fmt.Errorf("Incorrect arguments. Expecting no arguments")
    }
    value, err := sharedParams.Bytes()
    if err != nil {
        return nil, err
    }
    return value, nil
}

func GetAuditorKeypair(stub Wrapper, args [][]byte)([]byte, error) {
    if len(args) != 0 {
        return nil, fmt.Errorf("Incorrect arguments. Expecting no arguments")
    }

    // TODO We are cheat here, we should verify the request is from the 
    // auditor
    if string(auditorKeypair) == "NoAuditorKeyPair" {
        return nil, fmt.Errorf("NoAuditorKeyPair")
    }
    return auditorKeypair, nil
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
    fmt.Println("[Ocert Scheme] [Setup]")
    if len(args) != 0 {
        return nil, fmt.Errorf("Incorrect arguments. Expecting no arguments")
    }

    var err error;
    verifyProofLog, err = os.Create("/data/verifyProofLog480.txt")
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    auditorKeypair = []byte("NoAuditorKeyPair")
    serialNumber = big.NewInt(0)
    sharedParams = GenerateSharedParams()
    fmt.Printf("[Ocert Scheme] [Setup] sharedParams: ")
    fmt.Println(sharedParams)

    // Generate auditor's keypair
    PKa, SKa := EKeyGen(sharedParams)
    fmt.Printf("[Ocert Scheme] [Setup] auditor_pk: ")
    fmt.Println(PKa)
    PKaBytes, err := PKa.Bytes()
    if err != nil {
        return nil, err
    }
    err = stub.PutState("auditor_pk", PKaBytes)
    if err != nil {
        return nil, err
    }
    KPa := new(AuditorKeypair)
    KPa.PK = PKa.PK
    KPa.SK = SKa.SK

    // Generate RSA keypair
    rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, err
    }
    fmt.Printf("[Ocert Scheme] [Setup] rsa pk: ")
    fmt.Println(&rsaPrivateKey.PublicKey)

    rsaPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
    if err != nil {
        return nil, err
    }
    rsaPK := new(RSAPK)
    rsaPK.PK = rsaPublicKeyBytes
    rsaPKBytes, err := rsaPK.Bytes()
    if err != nil {
        return nil, err
    }
    err = stub.PutState("rsa_pk", rsaPKBytes)
    if err != nil {
        return nil, err
    }

    // Generate structure preserving keypair
    VKei, SKei := SKeyGen(sharedParams)
    fmt.Printf("[Ocert Scheme] [Setup] sVK: ")
    fmt.Println(VKei)
    sSigningKey = SKei
    SVKb, err := VKei.Bytes()
    if err != nil {
        return nil, err
    }
    err = stub.PutState("structure_preserving_vk", SVKb)
    if err != nil {
        return nil, err
    }

    // Setup constants for proof
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    G := pairing.NewG1().SetBytes(sharedParams.G1)
    H := pairing.NewG2().SetBytes(sharedParams.G2)
    consts = new(ProofConstants)
    consts.VK = VKei
    consts.PPrime = nil
    consts.PKa = PKa
    consts.Egh = pairing.NewGT().Pair(G, H).Bytes()
    consts.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VKei.Z)).Bytes()

    // Return keypair to the auditor
    KPab, err := KPa.Bytes()
    if err != nil {
        return nil, err
    }
    auditorKeypair = KPab
    return KPab, nil
}

/*
 * GenECert is used to generate an ecert of a client
 * It takes the client id and the client's public key, and returns
 * psudonym P and ecert to the client.
 */
func GenECert(stub Wrapper, args [][]byte) ([]byte, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("Incorrect arguments.")
    }

    request := new(GenECertRequest)
    err := request.SetBytes(args[0])
    if err != nil {
        return nil, err
    }

    IDc := new(ClientID)
    IDc.ID = request.IDc
    PKc := new(ClientPublicKey)
    PKc.PK = request.PKc

    fmt.Println("[Ocert Scheme] [GenECert]")
    fmt.Printf("[Ocert Scheme] [GenECert] IDc: ")
    fmt.Println(IDc)
    fmt.Printf("[Ocert Scheme] [GenECert] PKc: ")
    fmt.Println(PKc)

    // Generate pseudonym P
    valuePKa, err := stub.GetState("auditor_pk")
    if err != nil {
        return nil, err
    }
    if valuePKa == nil {
        return nil, fmt.Errorf("Asset not found: auditor_pk")
    }
    PKa := new(AuditorPublicKey)
    err = PKa.SetBytes(valuePKa)
    if err != nil {
        return nil, err
    }

    P := EEnc(sharedParams, PKa, IDc)
    fmt.Printf("[Ocert Scheme] [GenECert] P: ")
    fmt.Println(P)

    // Generate ecert
    ecert := SSign(sharedParams, sSigningKey, P, PKc)
    fmt.Printf("[Ocert Scheme] [GenECert] ecert: ")
    fmt.Println(ecert)

    reply := new(GenECertReply)
    reply.P, err = P.Bytes()
    if err != nil {
        return nil, err
    }
    reply.Ecert, err = ecert.Bytes()
    if err != nil {
        return nil, err
    }
    replyBytes, err := reply.Bytes()
    if err != nil {
        return nil, err
    }
    return replyBytes, nil
}

/*
 * GenOCert is used to generate an ocert of a client
 * It takes a client's public key, a client's pseudonym and the 
 * proof of knowledge, and returns the ocert to the client 
 */
func GenOCert(stub Wrapper, args [][]byte) ([]byte, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("Incorrect arguments.")
    }

    request := new(GenOCertRequest)
    err := request.SetBytes(args[0])
    if err != nil {
        return nil, err
    }

    PKc := new(ClientPublicKey)
    PKc.PK = request.PKc
    P := new(Pseudonym)
    err = P.SetBytes(request.P)
    if err != nil {
        return nil, err
    }
    pi := new(ProofOfKnowledge)
    err = pi.SetBytes(request.Pi)
    if err != nil {
        return nil, err
    }

    fmt.Println("[Ocert Scheme] [GenOCert]")
    fmt.Printf("[Ocert Scheme] [GenOert] PKc: ")
    fmt.Println(PKc)
    fmt.Printf("[Ocert Scheme] [GenOCert] P: ")
    fmt.Println(P)
    fmt.Printf("[Ocert Scheme] [GenOCert] pi: ")
    pi.Print()

    // Verify proof of knowledge
    start := time.Now()

    consts.PPrime = P
    result := PProve(sharedParams, pi, consts)
    consts.PPrime = nil

    end := time.Now()
    elapsed := end.Sub(start)
    fmt.Printf("[Ocert Scheme] [GenOCert] proof verfication time: ")
    fmt.Println(elapsed)
    fmt.Printf("[Ocert Scheme] [GenOCert] proof verfication result: ")
    fmt.Println(result)
    if !result {
        return nil, fmt.Errorf("Proof verfication fails")
    }
    verifyProofLog.WriteString("verifyProof: " + elapsed.String() + "\n")

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
    fmt.Printf("[Ocert Scheme] [GenOCert] signature: ")
    fmt.Println(signature)

    reply := new(GenOCertReply)
    reply.Sig = signature
    replyBytes, err := reply.Bytes()
    if err != nil {
        return nil, err
    }
    return replyBytes, nil
}