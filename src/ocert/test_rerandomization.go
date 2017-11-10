package ocert

import (
	"os"
)

/*
 * Run a single test
 */
func Etest() bool {
	PK, SK := EKeyGen()

	// TODO randomly generate id
	id := new(ClientID)
	P := EEnc(PK, id)

	decryptedId := EDec(SK, P)

	if id != decryptedId {
		return false
	}

	PPrime := ERerand(P)
	return ERerandVerify(P, PPrime)
}

/*
 * Run test b times
 */
func runETest(b int) {
	for i := 0; i < b; i++ {
		if !Etest() {
			os.Exit(1)
		}
	} 
}