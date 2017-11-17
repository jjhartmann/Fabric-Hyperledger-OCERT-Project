package ocert

import (
	"os"
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