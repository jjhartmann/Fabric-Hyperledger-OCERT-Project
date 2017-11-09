package ocert

import (
	"os"
)

/*
 * Run a single test
 */
func test() bool {
	pi := Setup()
	return Prove(pi)
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