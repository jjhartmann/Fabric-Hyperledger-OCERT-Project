package ocert

import "github.com/Nik-U/pbc"

/*
 * RMatrix: Builds a matrix of random elements from Zp
 */

type RMatrix struct {
  mat [][]*pbc.Element
}

func NewRMatrix(pairing *pbc.Pairing, rows int, cols int) *RMatrix {
  rmat := new(RMatrix)

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
