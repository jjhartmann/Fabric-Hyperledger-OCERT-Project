package ocert

import (
	"os"
)

/*
 * Run a single test
 */
func Ptest() bool {
	sharedParams := GenerateSharedParams()
	eqs := new(SystemOfEquations)
	vars := new(ProofVariables)
	pi := PSetup(sharedParams, eqs, vars)

	consts := new(ProofConstants)
	return PProve(sharedParams, pi, consts)
}

/*
 * Run test b times
 */
func RunPTest(b int) {
	for i := 0; i < b; i++ {
		if !Ptest() {
			os.Exit(1)
		}
	} 
}