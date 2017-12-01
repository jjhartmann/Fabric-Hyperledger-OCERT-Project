package ocert

import (
	"os"
	"github.com/Nik-U/pbc"
)

/*
 * Run a single test
 */
func Etest() bool {
	sharedParams := GenerateSharedParams()
	PK, SK := EKeyGen(sharedParams)

	// TODO randomly generate id
	id := new(ClientID)
	P := EEnc(sharedParams, PK, id)

	decryptedId := EDec(sharedParams, SK, P)

	if id != decryptedId {
		return false
	}

	PPrime := ERerand(sharedParams, P)
	return ERerandVerify(sharedParams, P, PPrime)
}

/*
 * Run test b times
 */
func RunETest(b int) {
	for i := 0; i < b; i++ {
		if !Etest() {
			os.Exit(1)
		}
	}
}

func EGenKeyTest() bool{
	sharedParams := GenerateSharedParams()
	PK, SK := EKeyGen(sharedParams)

	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	g1 := pairing.NewG1().SetBytes(sharedParams.G1)

	SKa:=pairing.NewZr().SetBytes(SK.SK)
  PKa:=pairing.NewG1().MulZn(g1,SKa)

	return PKa.Equals(pairing.NewG1().SetBytes(PK.PK))

}
