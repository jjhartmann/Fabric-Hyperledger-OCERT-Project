package main

import (
  "fmt"
	//"./src/ocert"
	"github.com/Nik-U/pbc"
)

type B1_Pair struct {
  b1 []byte
  b2 []byte
}

func iota(pairing *pbc.Pairing, el *pbc.Element) B1_Pair {
  var pair B1_Pair

  pair.b1 = el.Set0().Bytes()
  fmt.Printf("Pair_B1 = %s", pair.b1)

  pair.b2 = el.Bytes()
  fmt.Printf("Pair_B2 = %s", pair.b1)

  return pair

}

func main() {
  //fmt.Println(ocert.Stest())


  // Testing Stuff
  params := pbc.GenerateF(160)
  pairing := params.NewPairing()
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)


  fmt.Printf("g1 = %s\n", g1)
  fmt.Printf("g1_2 = %s\n", pairing.NewG1().Rand())
  fmt.Printf("g2 = %s\n", g2)
  fmt.Printf("gt = %s\n", gt)
  fmt.Println()

  gg := pairing.NewG1().Add(g1, g1)
  fmt.Printf("gg = %s\n", gg)




}
