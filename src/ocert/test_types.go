package ocert

import(
	"fmt"
	"github.com/Nik-U/pbc"
)

func RunTypesTest() {
	fmt.Println("RunTypesTest")

	sharedParams := GenerateSharedParams()
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

	_ = pairing
}