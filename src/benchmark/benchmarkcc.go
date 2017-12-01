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
}

func main () {
	// setup()
	fmt.Println(sharedParams())
	sharedParams := sharedParams()
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	IDc := new(ocert.ClientID)
	IDc.ID = pairing.NewG1().Rand().Bytes()
	PKc := new(ocert.ClientPublicKey)
	PKc.PK = pairing.NewG1().Rand().Bytes()
	genECert(IDc, PKc)
}
