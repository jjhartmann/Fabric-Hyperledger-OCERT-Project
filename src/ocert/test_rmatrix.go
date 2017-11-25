package ocert

import (
  "os"
  "fmt"
  "github.com/Nik-U/pbc"
)

func TestRMatrixGen(verbose bool) bool{
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  if (verbose) {fmt.Println("Starting RMatrix Test")}
  rows := 3
  cols := 2
  rmat := NewRMatrix(pairing, rows, cols)

  if (verbose) {
      for i := 0; i < rows; i++ {
        fmt.Printf("\n")
        for j := 0; j < cols; j++ {
          fmt.Printf("%s \t", rmat.mat[i][j])
        }
      }
      fmt.Println()
  }


  if (verbose) {fmt.Println("Creating CRS Sigma")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS

  if (verbose) {fmt.Println("Testing Muliplication on Commitment Key")}
  Ru := rmat.MulCommitmentKeysG1(pairing, sigma.U)

  if (verbose) {
    for i := 0; i < len(Ru); i++ {
      fmt.Printf("%s\t", pairing.NewG1().SetBytes(Ru[i].b1))
      fmt.Printf("%s\t\n", pairing.NewG1().SetBytes(Ru[i].b2))
    }
    fmt.Println()
  }

  return true

}






/*
 * Run test b times
 */
func RunRMatTest(b int) {
  for i := 0; i < b; i++ {
    if !TestRMatrixGen(false) {
      os.Exit(1)
    }
  }
}