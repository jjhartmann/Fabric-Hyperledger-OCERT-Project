package ocert

import (
  "github.com/Nik-U/pbc"
  "fmt"
)

/*
 * RMatrix: Builds a matrix of random elements from Zp
 */

type BMatrix struct {
  mat [][]*BPair
  rows int
  cols int
  invert bool
}

type RMatrix struct {
  mat [][]*pbc.Element
  rows int
  cols int
  invert bool
}

func NewRMatrix(pairing *pbc.Pairing, rows int, cols int) *RMatrix {
  rmat := new(RMatrix)
  rmat.rows = rows
  rmat.cols = cols
  for i := 0; i < rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < cols; j++ {
      el := pairing.NewZr().Rand()
      elementRow = append(elementRow, el)
    }
    rmat.mat = append(rmat.mat, elementRow)
  }
  return rmat
}


func NewRMatrixinG2(pairing *pbc.Pairing, rows int, cols int) *RMatrix {
  rmat := new(RMatrix)
  rmat.rows = rows
  rmat.cols = cols
  for i := 0; i < rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < cols; j++ {
      el := pairing.NewG2().Rand()
      elementRow = append(elementRow, el)
    }
    rmat.mat = append(rmat.mat, elementRow)
  }
  return rmat
}

func NewRMatrixinG1(pairing *pbc.Pairing, rows int, cols int) *RMatrix {
  rmat := new(RMatrix)
  rmat.rows = rows
  rmat.cols = cols
  for i := 0; i < rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < cols; j++ {
      el := pairing.NewG1().Rand()
      elementRow = append(elementRow, el)
    }
    rmat.mat = append(rmat.mat, elementRow)
  }
  return rmat
}

func NewOnesMatrix(pairing *pbc.Pairing, rows int, cols int) *RMatrix {
  rmat := new(RMatrix)
  rmat.rows = rows
  rmat.cols = cols
  for i := 0; i < rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < cols; j++ {
      el := pairing.NewZr().Set1()
      elementRow = append(elementRow, el)
    }
    rmat.mat = append(rmat.mat, elementRow)
  }
  return rmat
}

func NewIdentiyMatrix(pairing *pbc.Pairing, rows int, cols int) *RMatrix {
  rmat := new(RMatrix)
  rmat.rows = rows
  rmat.cols = cols
  for i := 0; i < rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < cols; j++ {
      el := pairing.NewZr().SetInt32(0)
      if i == j {
        el = pairing.NewZr().Set1()
      }
      elementRow = append(elementRow, el)
    }
    rmat.mat = append(rmat.mat, elementRow)
  }
  return rmat
}

// Group Type
// 1: G1
// 2: G2
// 3: GT
func (rmat *RMatrix) PrintAll(pairing *pbc.Pairing) {
  for i := 0 ;i < rmat.rows; i ++ {
    for j := 0; j < rmat.cols; j++ {
      fmt.Printf("%s, ", rmat.mat[i][j])
    }
    fmt.Println()
  }
}


func (rmat *RMatrix) MultBPairMatrixG2(pairing *pbc.Pairing, X *BMatrix) *BMatrix {
  if len(X.mat) != len(rmat.mat[0]) {
    panic("Matrix elements need to be compatiable")
  }

  retMat := new(BMatrix)
  retMat.rows = len(rmat.mat)
  retMat.cols = len(X.mat[0])

  for i := 0; i < rmat.rows; i++ {
    elementRow := []*BPair{}
    for j := 0; j < len(X.mat[0]); j++ {
      el := new(BPair)
      el.b1 = pairing.NewG2().Set1().Bytes()
      el.b2 = pairing.NewG2().Set1().Bytes()
      for k := 0; k < len(X.mat); k++ {
        tmpX := X.mat[k][j]
        tmp := tmpX.MulScalarInG2(pairing, rmat.mat[i][k])
        el = el.AddinG2(pairing, tmp)
      }
      elementRow = append(elementRow, el)
    }
    retMat.mat = append(retMat.mat, elementRow)
  }
  return retMat
}


func (rmat *RMatrix) MultBPairMatrixG1(pairing *pbc.Pairing, X *BMatrix) *BMatrix {
  if len(X.mat) != len(rmat.mat[0]) {
    panic("Matrix elements need to be compatiable")
  }

  retMat := new(BMatrix)
  retMat.rows = len(rmat.mat)
  retMat.cols = len(X.mat[0])

  for i := 0; i < rmat.rows; i++ {
    elementRow := []*BPair{}
    for j := 0; j < len(X.mat[0]); j++ {
      el := new(BPair)
      el.b1 = pairing.NewG1().Set1().Bytes()
      el.b2 = pairing.NewG1().Set1().Bytes()
      for k := 0; k < len(X.mat); k++ {
        tmpX := X.mat[k][j]
        tmp := tmpX.MulScalarInG1(pairing, rmat.mat[i][k])
        el = el.AddinG1(pairing, tmp)
      }
      elementRow = append(elementRow, el)
    }
    retMat.mat = append(retMat.mat, elementRow)
  }
  return retMat
}

func (rmat *RMatrix) MultElementArrayZr(pairing *pbc.Pairing, X [][]*pbc.Element) *RMatrix {
  if len(X) != len(rmat.mat[0]) {
    panic("Matrix elements need to be compatiable")
  }

  retMat := new(RMatrix)
  retMat.rows = len(rmat.mat)
  retMat.cols = len(X[0])

  for i := 0; i < rmat.rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < len(X[0]); j++ {
      el := pairing.NewZr().Set1()
      el := pairing.NewZr().Set0()
      for k := 0; k < len(X); k++ {
        tmp := pairing.NewZr().Mul(X[k][j], rmat.mat[i][k])
        el = pairing.NewZr().Add(el, tmp)
      }
      elementRow = append(elementRow, el)
    }
    retMat.mat = append(retMat.mat, elementRow)
  }
  return retMat
}





func (rmat *RMatrix) MultElementArrayG2(pairing *pbc.Pairing, X [][]*pbc.Element) *RMatrix {
  if len(X) != len(rmat.mat[0]) {
    panic("Matrix elements need to be compatiable")
  }

  retMat := new(RMatrix)
  retMat.rows = len(rmat.mat)
  retMat.cols = len(X[0])

  for i := 0; i < rmat.rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < len(X[0]); j++ {
      el := pairing.NewG2().Set1()
      for k := 0; k < len(X); k++ {
        tmp := pairing.NewG2().MulZn(X[k][j], rmat.mat[i][k])
        el = pairing.NewG2().Add(el, tmp)
      }
      elementRow = append(elementRow, el)
    }
    retMat.mat = append(retMat.mat, elementRow)
  }
  return retMat
}


func (rmat *RMatrix) MultElementArrayG1(pairing *pbc.Pairing, X [][]*pbc.Element) *RMatrix {
  if len(X) != len(rmat.mat[0]) {
    panic("Matrix elements need to be compatiable")
  }

  retMat := new(RMatrix)
  retMat.rows = len(rmat.mat)
  retMat.cols = len(X[0])

  for i := 0; i < rmat.rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < len(X[0]); j++ {
      el := pairing.NewG1().Set1()
      for k := 0; k < len(X); k++ {
        tmp := pairing.NewG1().MulZn(X[k][j], rmat.mat[i][k])
        el = pairing.NewG1().Add(el, tmp)
      }
      elementRow = append(elementRow, el)
    }
    retMat.mat = append(retMat.mat, elementRow)
  }
  return retMat
}

// TODO: A lot of these fucntions could be optimised!
func (rmat *RMatrix) InvertMatrix() *RMatrix{
  R := new(RMatrix)
  R.rows = rmat.cols
  R.cols = rmat.rows

  for j := 0; j < rmat.cols; j++{
    elrow := []*pbc.Element{}
    for i := 0; i < rmat.rows; i++ {
      elrow = append(elrow, rmat.mat[i][j])
    }
    R.mat = append(R.mat, elrow)
  }
  return R
}

func (rmat *RMatrix) ElementWiseSub(pairing *pbc.Pairing, L *RMatrix) *RMatrix {
  if rmat.cols != L.cols || rmat.rows != L.rows {
    panic("Rows and Cols need to be equivalent")
  }

  R := NewRMatrix(pairing, L.rows, L.cols)
  for i := 0; i < L.rows; i++{
    for j := 0; j < L.cols; j++ {
      R.mat[i][j] = pairing.NewZr().Sub(rmat.mat[i][j], L.mat[i][j])
    }
  }
  return R
}

func (rmat *RMatrix) MulBScalarinB1(pairing *pbc.Pairing, B BPair) [][]*BPair {
  Rb := [][]*BPair{}

  for i := 0; i < rmat.rows; i++ {
    pairRow := []*BPair{}
    for j := 0; j < rmat.cols; j++ {
      pair := new(BPair)
      pair.b1 = pairing.NewG1().MulZn(pairing.NewG1().SetBytes(B.b1), rmat.mat[i][j]).Bytes()
      pair.b2 = pairing.NewG1().MulZn(pairing.NewG1().SetBytes(B.b2), rmat.mat[i][j]).Bytes()
      pairRow = append(pairRow, pair)
    }
    Rb = append(Rb, pairRow)
  }
  return Rb
}

func (rmat *RMatrix) MulBScalarinB2(pairing *pbc.Pairing, B BPair) [][]*BPair {
  Rb := [][]*BPair{}

  for i := 0; i < rmat.rows; i++ {
    pairRow := []*BPair{}
    for j := 0; j < rmat.cols; j++ {
      pair := new(BPair)
      pair.b1 = pairing.NewG2().MulZn(pairing.NewG2().SetBytes(B.b1), rmat.mat[i][j]).Bytes()
      pair.b2 = pairing.NewG2().MulZn(pairing.NewG2().SetBytes(B.b2), rmat.mat[i][j]).Bytes()
      pairRow = append(pairRow, pair)
    }
    Rb = append(Rb, pairRow)
  }
  return Rb
}

func (rmat *RMatrix) MulScalarZn(pairing *pbc.Pairing, r *pbc.Element) *RMatrix {
  R := new(RMatrix)
  R.rows = rmat.rows
  R.cols = rmat.cols

  for i := 0; i < rmat.rows; i++ {
    elementRow := []*pbc.Element{}
    for j := 0; j < rmat.cols; j++ {
      el := pairing.NewZr().Mul(rmat.mat[i][j], r)
      elementRow = append(elementRow, el)
    }
    R.mat = append(R.mat, elementRow)
  }
  return R
}

func (rmat *RMatrix) MulCommitmentKeysG1(pairing *pbc.Pairing, U []CommitmentKey) []*BPair {
  rows := rmat.rows
  cols := len(U)
  Ru := []*BPair{}
  if (rmat.cols != len(U) ){
    panic("Error Occured in MulCommitmentKeys: CommitmentKeys incompatiable")
    return Ru
  }

  for i := 0; i < rows; i++ {

    // The BPair in B1
    B1 := pairing.NewG1().Set1()
    B2 := pairing.NewG1().Set1()

    for j := 0; j < cols; j++ {
      // Get pair (P, Q) form commitment key
      P := pairing.NewG1().SetBytes(U[j].u1)
      Q := pairing.NewG1().SetBytes(U[j].u2)

      // Get random r in Zp
      r := rmat.mat[i][j]

      // Multiple r by P and Q
      Pr := pairing.NewG1().MulZn(P, r)
      Qr := pairing.NewG1().MulZn(Q, r)

      // Add to Bpairs
      B1 = pairing.NewG1().Add(B1, Pr)
      B2 = pairing.NewG1().Add(B2, Qr)
    }

    // Append to BPair
    tmp := new(BPair)
    tmp.b1 = B1.Bytes()
    tmp.b2 = B2.Bytes()

    Ru = append(Ru, tmp)
  }

  return Ru
}

func (rmat *RMatrix) MulCommitmentKeysG2(pairing *pbc.Pairing, V []CommitmentKey) []*BPair {
  rows := rmat.rows
  cols := len(V)
  Rv := []*BPair{}
  if (rmat.cols != len(V) ){
    panic("Error Occured in MulCommitmentKeys: CommitmentKeys incompatiable")
    return Rv
  }

  for i := 0; i < rows; i++ {

    // The BPair in B1
    B1 := pairing.NewG2().Set1()
    B2 := pairing.NewG2().Set1()

    for j := 0; j < cols; j++ {
      // Get pair (P, Q) form commitment key
      P := pairing.NewG2().SetBytes(V[j].u1)
      Q := pairing.NewG2().SetBytes(V[j].u2)

      // Get random r in Zp
      r := rmat.mat[i][j]

      // Multiple r by P and Q
      Pr := pairing.NewG2().MulZn(P, r)
      Qr := pairing.NewG2().MulZn(Q, r)

      // Add to Bpairs
      B1 = pairing.NewG2().Add(B1, Pr)
      B2 = pairing.NewG2().Add(B2, Qr)
    }

    // Append to BPair
    tmp := new(BPair)
    tmp.b1 = B1.Bytes()
    tmp.b2 = B2.Bytes()

    Rv = append(Rv, tmp)
  }

  return Rv
}