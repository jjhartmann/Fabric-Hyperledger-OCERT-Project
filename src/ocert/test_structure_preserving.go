package ocert

import (
	"os"
)

/*
 * Run a single test
 */
func test() bool {
	VK, SK := SKeyGen()

	// TODO randomly generate P and PKc
	P := new(Pseudonym)
	PKc := new(ClientPublicKey)

	ecert := SSign(SK, P, PKc)

	return SVerify(VK, P, PKc, ecert)
}

/*
 * Run test b times
 */
func runTest(b int) {
	for i := 0; i < b; i++ {
		if !test() {
			os.Exit(1)
		}
	} 
}