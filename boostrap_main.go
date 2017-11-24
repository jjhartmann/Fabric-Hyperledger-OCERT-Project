package main

import (
  "fmt"
	"./src/ocert"
	"github.com/Nik-U/pbc"
)


func main() {
  //fmt.Println(ocert.Stest())
  fmt.Println(ocert.IotaTest(true))

  // Testing Stuff
  params := pbc.GenerateF(160)
  pairing := params.NewPairing()
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)



}
