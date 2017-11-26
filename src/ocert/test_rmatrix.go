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
    fmt.Println()
  }


  if (verbose) {fmt.Println("Creating CRS Sigma")}
  alpha := pairing.NewZr().Rand() // Secret Key
  sigma := CreateCommonReferenceString(sharedParams, alpha) // CRS


  if (verbose) {fmt.Println("Testing Muliplication on Commitment Key in G1")}
  Ru := rmat.MulCommitmentKeysG1(pairing, sigma.U)

  if (verbose) {
    for i := 0; i < len(Ru); i++ {
      fmt.Printf("%s\t", pairing.NewG1().SetBytes(Ru[i].b1))
      fmt.Printf("%s\t\n", pairing.NewG1().SetBytes(Ru[i].b2))
    }
    fmt.Println()
    fmt.Println()
  }


  if (verbose) {fmt.Println("Testing Muliplication on Commitment Key in G2")}
  Rv := rmat.MulCommitmentKeysG2(pairing, sigma.V)

  if (verbose) {
    for i := 0; i < len(Rv); i++ {
      fmt.Printf("%s\t", pairing.NewG2().SetBytes(Rv[i].b1))
      fmt.Printf("%s\t\n", pairing.NewG2().SetBytes(Rv[i].b2))
    }
    fmt.Println()
  }

  return true

}

func TestRMatrixMulSclarInZn(verbose bool, rows int, cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  R := NewRMatrix(pairing, rows, cols)
  x := pairing.NewZr().Rand()
  xR := R.MulScalarZn(pairing, x)

  ret := true
  for i := 0; i < R.rows; i++ {
    for j := 0; j < R.cols; j++ {
      tmp := pairing.NewZr().Mul(R.mat[i][j], x)
      ret = tmp.Equals(xR.mat[i][j]) && ret
    }
  }
  return ret
}

func TestElementWiseSubtraction(verbose bool, rows int, cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  L := NewRMatrix(pairing, rows, cols)
  R := NewRMatrix(pairing, rows, cols)
  LR := L.ElementWiseSub(pairing, R)

  ret := true
  for i := 0; i < L.rows; i++ {
    for j := 0; j < R.cols; j++ {
      temp := pairing.NewZr().Sub(L.mat[i][j], R.mat[i][j])
      ret = ret && LR.mat[i][j].Equals(temp)
    }
  }
  return ret
}

/*
 Run all Matrix Tests
 */
func RunAllRTests(verbose bool) {
  fmt.Println("RMatrix Generator          ", TestRMatrixGen(verbose))
  fmt.Println("RMatrix Scalar Mul 1x1     ", TestRMatrixMulSclarInZn(true, 1, 1))
  fmt.Println("RMatrix Scalar Mul 10x1    ", TestRMatrixMulSclarInZn(true, 10, 1))
  fmt.Println("RMatrix Scalar Mul 10x10   ", TestRMatrixMulSclarInZn(true, 10, 10))
  fmt.Println("RMatrix EW Subtract 1x1    ", TestElementWiseSubtraction(true, 1, 1))
  fmt.Println("RMatrix EW Subtract 10x1   ", TestElementWiseSubtraction(true, 10, 1))
  fmt.Println("RMatrix EW Subtract 10x10  ", TestElementWiseSubtraction(true, 10, 10))
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