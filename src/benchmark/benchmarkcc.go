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
	"crypto/x509"
)

func parseOut(out []byte) []byte {
	str := string(out)
	str = str[14:]
	return []byte(str)
}

func setup() {
	installCmd := "peer chaincode install -p chaincodedev/chaincode/ocert -n myccc -v 0"
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

func genECert(id *ocert.ClientID, pkc *ocert.ClientPublicKey) {
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
	fmt.Println(out)

	return

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
	fmt.Println(p)
	fmt.Println(ecert)

}

func genOCert(pkc *ocert.ClientPublicKey, p *ocert.Pseudonym) {
	request := new(ocert.GenOCertRequest)
	request.PKc = pkc.PK
	pBytes, err := p.Bytes()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	request.P = pBytes

	requestBytes, err := request.Bytes()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	requestStr := string(requestBytes)
	requestStr = strings.Replace(requestStr, "\"", "\\\"", -1)

	queryCmd := "peer chaincode query -n mycc -c '{\"Args\":[\"genOCert\" ,\"" +
				requestStr + "\"]}' -C myc"
	fmt.Println(queryCmd)

	return

	start := time.Now()
	out, err := exec.Command("sh","-c", queryCmd).Output()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("genOCert: ")
	fmt.Println(elapsed)

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Println(out)

	out = parseOut(out)
	reply := new(ocert.GenOCertReply)
	err = reply.SetBytes(out)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func main () {
	// setup()
	sharedParams := sharedParams()
	fmt.Println(sharedParams)
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	IDc := new(ocert.ClientID)
	IDc.ID = pairing.NewG1().Rand().Bytes()
	PKc := new(ocert.ClientPublicKey)
	PKc.PK = pairing.NewG1().Rand().Bytes()
	// genECert(IDc, PKc)

	start := time.Now()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("proof generation: ")
	fmt.Println(elapsed)

	fmt.Printf("rsa pk: ")
	fmt.Println(rsaPK())
}
