package ocert

import (
	"os"
)

/*
 * Run a single test
 */
func Ptest() bool {
	pi := Setup()
	return Prove(pi)
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