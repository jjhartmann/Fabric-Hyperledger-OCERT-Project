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
  return ret && TestRMatrixStructure(verbose, xR)
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
  return ret && TestRMatrixStructure(verbose, LR)
}

func TestRMatrixBPairScalar(verbose bool, rows int, cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  R := NewRMatrix(pairing, rows, cols)
  B := BPair{pairing.NewG1().Rand().Bytes(), pairing.NewG1().Rand().Bytes()}
  Rb := R.MulBScalarinB1(pairing, B)

  ret1 := len(Rb) == len(R.mat) && len(Rb[0]) == len(R.mat[0])
  if(verbose) {fmt.Println("Equality on length: ", ret1)}

  ret2 := true
  for i := 0; i < R.rows; i++ {
    for j := 0; j < R.cols; j++ {
      b1 := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(B.b1), R.mat[i][j])
      b2 := pairing.NewG1().MulZn(pairing.NewG1().SetBytes(B.b2), R.mat[i][j])

      ret2 = ret2 &&
        pairing.NewG1().SetBytes(Rb[i][j].b1).Equals(b1) &&
        pairing.NewG1().SetBytes(Rb[i][j].b2).Equals(b2)
    }
  }
  return ret2 && ret1 && TestRMatrixStructure(verbose, R)
}

func TestRMatrixInversion(verbose bool, rows int, cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  R := NewRMatrix(pairing, rows, cols)
  Ri := R.InvertMatrix()

  ret := true
  for i := 0; i < rows; i++ {
    for j := 0; j < cols; j++ {
      ret = R.mat[i][j].Equals(Ri.mat[j][i])
    }
  }
  return ret && TestRMatrixStructure(verbose, R) && TestRMatrixStructure(verbose, Ri)
}


func TestRMatrixMultiplicationforElementinG2(verbose bool, r_rows int, r_cols int, x_rows int, x_cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  ones := NewOnesMatrix(pairing, r_rows, r_cols)
  X := NewRMatrixinG2(pairing, x_rows, x_cols)

  if verbose {
    fmt.Println("X:")
    for i := 0; i < len(X.mat); i++ {
      for j := 0; j < len(X.mat[0]); j++ {
        fmt.Printf("%s", X.mat[i][j])
      }
      fmt.Println()
    }
  }

  mat := ones.MultElementArrayG2(pairing, X.mat)
  if verbose {
    fmt.Println("RET:")
    for i := 0; i < mat.rows; i++ {
      for j := 0; j < mat.cols; j++ {
        fmt.Printf("%s", mat.mat[i][j])
      }
      fmt.Println()
    }
  }

  ret1 := len(mat.mat) == len(ones.mat) && len(mat.mat[0]) == len(X.mat[0])

  return TestRMatrixStructure(verbose, mat) && ret1

}

func TestRMatrixMultiplicationforElementinG1(verbose bool, r_rows int, r_cols int, x_rows int, x_cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  ones := NewOnesMatrix(pairing, r_rows, r_cols)
  X := NewRMatrixinG1(pairing, x_rows, x_cols)

  if verbose {
    fmt.Println("X:")
    for i := 0; i < len(X.mat); i++ {
      for j := 0; j < len(X.mat[0]); j++ {
        fmt.Printf("%s", X.mat[i][j])
      }
      fmt.Println()
    }
  }

  mat := ones.MultElementArrayG1(pairing, X.mat)
  if verbose {
    fmt.Println("RET:")
    for i := 0; i < mat.rows; i++ {
      for j := 0; j < mat.cols; j++ {
        fmt.Printf("%s", mat.mat[i][j])
      }
      fmt.Println()
    }
  }

  ret1 := len(mat.mat) == len(ones.mat) && len(mat.mat[0]) == len(X.mat[0])

  return TestRMatrixStructure(verbose, mat) && ret1

}

func TestRMatrixMultiplicationforElementinZr(verbose bool, r_rows int, r_cols int, x_rows int, x_cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  ones := NewOnesMatrix(pairing, r_rows, r_cols)
  X := NewRMatrix(pairing, x_rows, x_cols)

  if verbose {
    fmt.Println("X:")
    for i := 0; i < len(X.mat); i++ {
      for j := 0; j < len(X.mat[0]); j++ {
        fmt.Printf("%s, ", X.mat[i][j])
      }
      fmt.Println()
    }
  }

  mat := ones.MultElementArrayZr(pairing, X.mat)
  if verbose {
    fmt.Println("RET:")
    for i := 0; i < mat.rows; i++ {
      for j := 0; j < mat.cols; j++ {
        fmt.Printf("%s, ", mat.mat[i][j])
      }
      fmt.Println()
    }
  }

  ret1 := len(mat.mat) == len(ones.mat) && len(mat.mat[0]) == len(X.mat[0])

  return TestRMatrixStructure(verbose, mat) && ret1

}


func TestRMatrixMultiplicationforBPairMatrixinG2(verbose bool, r_rows int, r_cols int, x_rows int, x_cols int) bool {
  sharedParams := GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  ones := NewOnesMatrix(pairing, r_rows, r_cols)
  X := new(BMatrix)

  for i := 0; i < x_rows; i++ {
    elementRow := []*BPair{}
    for j := 0; j < x_cols; j++ {
      el := new(BPair)
      el.b1 = pairing.NewG2().Rand().Bytes()
      el.b2 =  pairing.NewG2().Rand().Bytes()
      elementRow = append(elementRow, el)
    }
    X.mat = append(X.mat, elementRow)
  }

  if verbose {
    fmt.Println("X:")
    for i := 0; i < len(X.mat); i++ {
      for j := 0; j < len(X.mat[0]); j++ {
        fmt.Printf("[%d, %d]: %v, ",i, j, X.mat[i][j])
      }
      fmt.Println()
    }
  }

  BMat := ones.MultBPairMatrixG2(pairing, X)
  if verbose {
    fmt.Println("RET:")
    for i := 0; i < len(BMat.mat); i++ {
      for j := 0; j < len(BMat.mat[0]); j++ {
        fmt.Printf("[%d, %d]: %v, ", BMat.mat[i][j])
      }
      fmt.Println()
    }
  }

  ret1 := len(BMat.mat) == len(ones.mat) && len(BMat.mat[0]) == len(X.mat[0])
  return ret1 && TestBMatrixStructure(verbose, BMat)

}

func TestRMatrixStructure(verbose bool, R *RMatrix) bool{
   if (verbose) {fmt.Println("Testing R Matrix Structure")}
   ret1 := R.rows == len(R.mat)
   if (verbose){fmt.Println("Rows Equality: ", ret1)}
   ret2 := R.cols == len(R.mat[0])
   if (verbose){fmt.Println("Cols Equality: ", ret2)}
   return ret1 && ret2
}

func TestBMatrixStructure(verbose bool, R *BMatrix) bool{
  if (verbose) {fmt.Println("Testing R Matrix Structure")}
  ret1 := R.rows == len(R.mat)
  if (verbose){fmt.Println("Rows Equality: ", ret1)}
  ret2 := R.cols == len(R.mat[0])
  if (verbose){fmt.Println("Cols Equality: ", ret2)}
  return ret1 && ret2
}

/*
 Run all Matrix Tests
 */
func RunAllRTests(verbose bool) {
  fmt.Println("RMatrix Generator                     ", TestRMatrixGen(verbose))
  fmt.Println("RMatrix Scalar Mul 1x1                ", TestRMatrixMulSclarInZn(false, 1, 1))
  fmt.Println("RMatrix Scalar Mul 10x1               ", TestRMatrixMulSclarInZn(false, 10, 1))
  fmt.Println("RMatrix Scalar Mul 10x10              ", TestRMatrixMulSclarInZn(false, 10, 10))
  fmt.Println("RMatrix EW Subtract 1x1               ", TestElementWiseSubtraction(false, 1, 1))
  fmt.Println("RMatrix EW Subtract 10x1              ", TestElementWiseSubtraction(false, 10, 1))
  fmt.Println("RMatrix EW Subtract 10x10             ", TestElementWiseSubtraction(false, 10, 10))
  fmt.Println("RMatrix BPair Scalar 1x1              ", TestRMatrixBPairScalar(false, 1, 1))
  fmt.Println("RMatrix BPair Scalar 10x1             ", TestRMatrixBPairScalar(false, 10, 1))
  fmt.Println("RMatrix BPair Scalar 10x10            ", TestRMatrixBPairScalar(false, 10, 10))
  fmt.Println("RMatrix Inversion 1x1                 ", TestRMatrixInversion(false, 1, 1))
  fmt.Println("RMatrix Inversion 10x1                ", TestRMatrixInversion(false, 10, 1))
  fmt.Println("RMatrix Inversion 3x7                 ", TestRMatrixInversion(false, 3, 7))
  fmt.Println("RMatrix Inversion 4x8                 ", TestRMatrixInversion(false, 4, 8))
  fmt.Println("RMatrix Mult in G2 1x1 1x1            ", TestRMatrixMultiplicationforElementinG2(false,1, 1, 1, 1))
  fmt.Println("RMatrix Mult in G2 2x1 1x2            ", TestRMatrixMultiplicationforElementinG2(false,2, 1, 1, 2))
  fmt.Println("RMatrix Mult in G2 1x2 2x1            ", TestRMatrixMultiplicationforElementinG2(false,1, 2, 2, 1))
  fmt.Println("RMatrix Mult in G2 2x2 2x1            ", TestRMatrixMultiplicationforElementinG2(false,2, 2, 2, 1))
  fmt.Println("RMatrix Mult in G2 1x2 2x2            ", TestRMatrixMultiplicationforElementinG2(false,1, 2, 2, 2))
  fmt.Println("RMatrix Mult in G2 2x2 2x2            ", TestRMatrixMultiplicationforElementinG2(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult in G2 2x2 2x2            ", TestRMatrixMultiplicationforElementinG2(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult in G1 1x1 1x1            ", TestRMatrixMultiplicationforElementinG1(false,1, 1, 1, 1))
  fmt.Println("RMatrix Mult in G1 2x1 1x2            ", TestRMatrixMultiplicationforElementinG1(false,2, 1, 1, 2))
  fmt.Println("RMatrix Mult in G1 1x2 2x1            ", TestRMatrixMultiplicationforElementinG1(false,1, 2, 2, 1))
  fmt.Println("RMatrix Mult in G1 2x2 2x1            ", TestRMatrixMultiplicationforElementinG1(false,2, 2, 2, 1))
  fmt.Println("RMatrix Mult in G1 1x2 2x2            ", TestRMatrixMultiplicationforElementinG1(false,1, 2, 2, 2))
  fmt.Println("RMatrix Mult in G1 2x2 2x2            ", TestRMatrixMultiplicationforElementinG1(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult in G1 2x2 2x2            ", TestRMatrixMultiplicationforElementinG1(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult in Zr 1x1 1x1            ", TestRMatrixMultiplicationforElementinZr(false,1, 1, 1, 1))
  fmt.Println("RMatrix Mult in Zr 2x1 1x2            ", TestRMatrixMultiplicationforElementinZr(false,2, 1, 1, 2))
  fmt.Println("RMatrix Mult in Zr 1x2 2x1            ", TestRMatrixMultiplicationforElementinZr(false,1, 2, 2, 1))
  fmt.Println("RMatrix Mult in Zr 2x2 2x1            ", TestRMatrixMultiplicationforElementinZr(false,2, 2, 2, 1))
  fmt.Println("RMatrix Mult in Zr 1x2 2x2            ", TestRMatrixMultiplicationforElementinZr(false,1, 2, 2, 2))
  fmt.Println("RMatrix Mult in Zr 2x2 2x2            ", TestRMatrixMultiplicationforElementinZr(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult in Zr 2x2 2x2            ", TestRMatrixMultiplicationforElementinZr(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult BPair Mat G2 1x1 1x1     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,1, 1, 1, 1))
  fmt.Println("RMatrix Mult BPair Mat G2 2x1 1x2     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,2, 1, 1, 2))
  fmt.Println("RMatrix Mult BPair Mat G2 1x2 2x1     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,1, 2, 2, 1))
  fmt.Println("RMatrix Mult BPair Mat G2 2x2 2x1     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,2, 2, 2, 1))
  fmt.Println("RMatrix Mult BPair Mat G2 1x2 2x2     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,1, 2, 2, 2))
  fmt.Println("RMatrix Mult BPair Mat G2 2x2 2x2     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,2, 2, 2, 2))
  fmt.Println("RMatrix Mult BPair Mat G2 2x2 2x2     ", TestRMatrixMultiplicationforBPairMatrixinG2(false,2, 2, 2, 2))

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