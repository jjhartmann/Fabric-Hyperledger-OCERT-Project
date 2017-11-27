package main

import (
	"os/exec"
    "fmt"
    "ocert"
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


func main () {
    // setup()

    fmt.Println(sharedParams())
}
