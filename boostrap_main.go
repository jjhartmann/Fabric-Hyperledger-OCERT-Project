package main

import (
  "fmt"
	"./src/ocert"
	"github.com/Nik-U/pbc"
)

type BPair struct {
  b1 []byte
  b2 []byte
}

func Iota1(pairing *pbc.Pairing, el *pbc.Element) *BPair {
  pair := new(BPair)

  pair.b1 = pairing.NewG1().Set0().Bytes()
  //fmt.Printf("Pair_B1 = %s", pair.b1)

  pair.b2 = el.Bytes()
  //fmt.Printf("Pair_B2 = %s", pair.b1)

  return pair

}

func Rho1(pairing *pbc.Pairing, pair *BPair, alpha *pbc.Element) *pbc.Element {
  Z1 := pairing.NewG1().SetBytes(pair.b1)
  Z2 := pairing.NewG1().SetBytes(pair.b2)
  tmp := pairing.NewG1().MulZn(Z1, alpha)
  return pairing.NewG1().Sub(Z2, tmp)
}

func main() {
  //fmt.Println(ocert.Stest())
  fmt.Println(ocert.IotaTest(true))


  //// Testing Stuff
  //params := pbc.GenerateF(160)
  //pairing := params.NewPairing()
  //g1 := pairing.NewG1().Rand()
  //g2 := pairing.NewG2().Rand()
  //gt := pairing.NewGT().Pair(g1, g2)
  //
  //
  //fmt.Printf("g1 = %s\n", g1)
  //fmt.Printf("g1_2 = %s\n", pairing.NewG1().Rand())
  //fmt.Printf("g2 = %s\n", g2)
  //fmt.Printf("gt = %s\n", gt)
  //fmt.Println()
  //
  //gg := pairing.NewG1().Add(g1, g1)
  //fmt.Printf("gg = %s\n", gg)
  //
  //
  //fmt.Println("\n/////////////////////////////\n" +
  //               "Testing Iota and p conversion\n")
  //
  //Z := pairing.NewG1().Rand()
  //b1pair := Iota1(pairing, Z)
  //
  //tb1 := pairing.NewG1().SetBytes(b1pair.b1)
  //tb2 := pairing.NewG1().SetBytes(b1pair.b2)
  //
  //fmt.Printf("Z = %s\n", Z)
  //fmt.Printf("B1.b1 = %s\n", tb1)
  //fmt.Printf("B1.b1 = %s\n", tb2)
  //
  //alpha := pairing.NewZr().Rand()
  //Zprime := Rho1(pairing, b1pair, alpha)
  //fmt.Printf("Zprime = %s\n", Zprime)
  //ret := Zprime.Equals(Z)
  //fmt.Println("Test ==", ret)

}
