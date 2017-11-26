package ocert

import (
  "github.com/Nik-U/pbc"
  "fmt"
)

/*
 * RMatrix: Builds a matrix of random elements from Zp
 */

type RMatrix struct {
  mat [][]*pbc.Element
  rows int
  cols int
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

func (rmat *RMatrix) ElementWiseSub(pairing *pbc.Pairing, L *RMatrix) *RMatrix {
  if rmat.cols != L.cols || rmat.rows != L.rows {
    fmt.Errorf("Rows and Cols need to be equivalent")
    return nil
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

func (rmat *RMatrix) MulScalarZn(pairing *pbc.Pairing, r *pbc.Element) *RMatrix {
  R := new(RMatrix)
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
    fmt.Errorf("Error Occured in MulCommitmentKeys: CommitmentKeys incompatiable\n%s", U)
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
    fmt.Errorf("Error Occured in MulCommitmentKeys: CommitmentKeys incompatiable\n%s", V)
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