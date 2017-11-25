package ocert

import (
  "os"
  "fmt"
  "github.com/Nik-U/pbc"
)

func TestRMatrix(verbose bool) bool{
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if (verbose) {fmt.Println("Starting RMatrix Test")}
  rows := 5
  cols := 2
  rmat := NewRMatrix(pairing, rows, cols)

  if (verbose) {
      for i := 0; i < rows; i++ {
        fmt.Printf("\n")
        for j := 0; j < cols; j++ {
          fmt.Printf("%s \t", rmat.mat[i][j])
        }
      }
  }

  return true

}






/*
 * Run test b times
 */
func RunRMatTest(b int) {
  for i := 0; i < b; i++ {
    if !TestRMatrix(false) {
      os.Exit(1)
    }
  }
}