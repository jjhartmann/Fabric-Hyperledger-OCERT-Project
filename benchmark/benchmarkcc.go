/*
 * Used to benchmark ocert chaincode
 */

package main

import (
    "os/exec"
    "fmt"
    "ocert"
    "github.com/Nik-U/pbc"
    "strings"
    "time"
    "crypto"
    "crypto/x509"
    "crypto/sha256"
    "crypto/rsa"
    "os"
)

func parseOut(out []byte) []byte {
    str := string(out)
    str = str[14:]
    return []byte(str)
}

func setup() {
    installCmd := "peer chaincode install -p chaincodedev/chaincode/ocert -n mycc -v 0"
    _, err := exec.Command("sh","-c", installCmd).Output()
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    instantiateCmd := "peer chaincode instantiate -n mycc -v 0 -c '{\"Args\":[]}' -C myc"
    _, err = exec.Command("sh","-c", instantiateCmd).Output()
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
}

func auditorPK() *ocert.AuditorPublicKey {
    queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"get\", \"auditor_pk\"]}' -C myc"
    out, err := exec.Command("sh","-c", queryCmd).Output()

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    
    auditorPK := new(ocert.AuditorPublicKey)
    err = auditorPK.SetBytes(parseOut(out))

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    return auditorPK
}

func sVK() *ocert.SVerificationKey {
    queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"get\", \"structure_preserving_vk\"]}' -C myc"
    out, err := exec.Command("sh","-c", queryCmd).Output()

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    
    sVK := new(ocert.SVerificationKey)
    err = sVK.SetBytes(parseOut(out))

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    return sVK
}

func sharedParams() *ocert.SharedParams {
    queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"sharedParams\"]}' -C myc"
    out, err := exec.Command("sh","-c", queryCmd).Output()

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    
    sharedParams := new(ocert.SharedParams)
    err = sharedParams.SetBytes(parseOut(out))

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    return sharedParams
}

func rsaPK() (interface{}) {
    queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"get\",\"rsa_pk\"]}' -C myc"
    out, err := exec.Command("sh","-c", queryCmd).Output()

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    rsaPK := new(ocert.RSAPK)
    err = rsaPK.SetBytes(parseOut(out))

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    pk, err := x509.ParsePKIXPublicKey(rsaPK.PK)
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    return pk
}

func genECert(id *ocert.ClientID, pkc *ocert.ClientPublicKey) (*ocert.Pseudonym, *ocert.Ecert){
    request := new(ocert.GenECertRequest)
    request.IDc = id.ID
    request.PKc = pkc.PK
    requestBytes, err := request.Bytes()
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    requestStr := string(requestBytes)
    requestStr = strings.Replace(requestStr, "\"", "\\\"", -1)

    queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"genECert\" ,\"" +
                requestStr + "\"]}' -C myc"

    out, err := exec.Command("sh","-c", queryCmd).Output()

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    out = parseOut(out)
    reply := new(ocert.GenECertReply)
    err = reply.SetBytes(out)
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    p := new(ocert.Pseudonym)
    err = p.SetBytes(reply.P)
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    ecert := new(ocert.Ecert)
    err = ecert.SetBytes(reply.Ecert)
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    return p, ecert
}

func genOCert(sharedParams *ocert.SharedParams, 
              p *ocert.Pseudonym, 
              auditorPK *ocert.AuditorPublicKey,
              vars *ocert.ProofVariables) (*ocert.ClientPublicKey, *ocert.Pseudonym, []byte){
    fmt.Println("[Benchmarkcc] [genOCert]------------------------------------------------")
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

    // New client public key and pseudonym
    newPKc := new(ocert.ClientPublicKey)
    newPKc.PK = pairing.NewG2().Rand().Bytes()
    fmt.Printf("[Benchmarkcc] newPKc: ")
    fmt.Println(newPKc)

    newP, rprime := ocert.ERerand(sharedParams, auditorPK, p)
    fmt.Printf("[Benchmarkcc] newP: ")
    fmt.Println(newP)
    fmt.Printf("[Benchmarkcc] rprime: ")
    fmt.Println(rprime)

    vars.RPrime = rprime

    // Proof generation
    start := time.Now()

    pi := ocert.PSetup(sharedParams, vars)
    
    end := time.Now()
    elapsed := end.Sub(start)
    fmt.Printf("[Benchmarkcc] ProofGeneration time: ")
    fmt.Println(elapsed)
    genProofLog.WriteString("genProof: " + elapsed.String() + "\n")

    fmt.Printf("[Benchmarkcc] pi: ")
    pi.Print()

    request := new(ocert.GenOCertRequest)
    request.PKc = newPKc.PK

    pBytes, err := newP.Bytes()
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    request.P = pBytes

    piBytes, err := pi.Bytes()
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    request.Pi = piBytes

    requestBytes, err := request.Bytes()
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    requestStr := string(requestBytes)
    requestStr = strings.Replace(requestStr, "\"", "\\\"", -1)

    queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"genOCert\" ,\"" +
                requestStr + "\"]}' -C myc"

    out, err := exec.Command("sh","-c", queryCmd).Output()

    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    out = parseOut(out)
    reply := new(ocert.GenOCertReply)
    err = reply.SetBytes(out)
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    return newPKc, newP, reply.Sig
}

var genOCertLog *os.File
var genProofLog *os.File
var genECertLog *os.File

func main () {
    var err error
    genOCertLog, err = os.Create("/data/genOCertLog320.txt")
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    genProofLog, err = os.Create("/data/genProofLog320.txt")
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    genECertLog, err = os.Create("/data/genECertLog320.txt")
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }

    // Setup
    setup()

    // Wait for chaincode Init call finish
    time.Sleep(3000 * time.Millisecond)

    // Keys
    auditorPK := auditorPK()
    fmt.Printf("[Benchmarkcc] auditorPK: ")
    fmt.Println(auditorPK)

    rsaPK := rsaPK()
    fmt.Printf("[Benchmarkcc] rsa_pk: ")
    fmt.Println(rsaPK)

    sVK := sVK()
    fmt.Printf("[Benchmarkcc] sVK: ")
    fmt.Println(sVK)

    // Bilinear group
    sharedParams := sharedParams()
    fmt.Printf("[Benchmarkcc] sharedParams: ")
    fmt.Println(sharedParams)
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    H := pairing.NewG2().SetBytes(sharedParams.G2)

    // Benchmark GenECert, We don't need to generate multiple ecerts in real use case,
    // one ecert is enought
    for i := 0; i < 100; i++ {
        start := time.Now()
        
        IDc := new(ocert.ClientID)
        IDc.ID = pairing.NewG1().Rand().Bytes()
        fmt.Printf("[Benchmarkcc] IDc: ")
        fmt.Println(IDc)

        PKc := new(ocert.ClientPublicKey)
        Xc := pairing.NewZr().Rand().Bytes()
        PKc.PK = pairing.NewG2().MulZn(H, pairing.NewZr().SetBytes(Xc)).Bytes()
        fmt.Printf("[Benchmarkcc] PKc: ")
        fmt.Println(PKc)
        P, ecert := genECert(IDc, PKc)
        fmt.Printf("[Benchmarkcc] P: ")
        fmt.Println(P)
        fmt.Printf("[Benchmarkcc] ecert: ")
        fmt.Println(ecert)

        vars := new(ocert.ProofVariables)
        vars.PKa = auditorPK
        vars.P = P
        vars.VK = sVK
        vars.RPrime = nil
        vars.PKc = PKc
        vars.Xc = Xc
        vars.E = ecert

        end := time.Now()
        elapsed := end.Sub(start)
        fmt.Printf("[Benchmarkcc] genECert: ")
        fmt.Println(elapsed)
        genECertLog.WriteString("genECert: " + elapsed.String() + "\n")
    }


    // GenECert
    IDc := new(ocert.ClientID)
    IDc.ID = pairing.NewG1().Rand().Bytes()
    fmt.Printf("[Benchmarkcc] IDc: ")
    fmt.Println(IDc)

    PKc := new(ocert.ClientPublicKey)
    Xc := pairing.NewZr().Rand().Bytes()
    PKc.PK = pairing.NewG2().MulZn(H, pairing.NewZr().SetBytes(Xc)).Bytes()
    fmt.Printf("[Benchmarkcc] PKc: ")
    fmt.Println(PKc)
    P, ecert := genECert(IDc, PKc)
    fmt.Printf("[Benchmarkcc] P: ")
    fmt.Println(P)
    fmt.Printf("[Benchmarkcc] ecert: ")
    fmt.Println(ecert)

    vars := new(ocert.ProofVariables)
    vars.PKa = auditorPK
    vars.P = P
    vars.VK = sVK
    vars.RPrime = nil
    vars.PKc = PKc
    vars.Xc = Xc
    vars.E = ecert

    for i := 0; i < 100; i++ {
        fmt.Printf("iteration %d\n", i)
        // GenOCert
        start := time.Now()
        
        newPKc, newP, signature := genOCert(sharedParams, P, auditorPK, vars)
        
        end := time.Now()
        elapsed := end.Sub(start)
        fmt.Printf("[Benchmarkcc] genOCert: ")
        fmt.Println(elapsed)
        genOCertLog.WriteString("genOCert: " + elapsed.String() + "\n")
        
        fmt.Printf("[Benchmarkcc] signature: ")
        fmt.Println(signature)

        // Verify signature
        msg, err := ocert.OCertSingedBytes(newPKc, newP)
        if err != nil {
            fmt.Println(err)
            panic(err.Error())
        }
        hashed := sha256.Sum256(msg)
        err = rsa.VerifyPKCS1v15(rsaPK.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature)
        if err != nil {
            fmt.Println(err)
            panic(err.Error())
        } else {
            fmt.Println("[Benchmarkcc] ocert verified")
        }
    }
}
