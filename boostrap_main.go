package main

import (
  "fmt"
	 "./src/ocert"
  "github.com/Nik-U/pbc"
)

func main() {
  fmt.Printf("\nRun Structure Perserving Tests\n")
  fmt.Println(ocert.Stest())

  fmt.Printf("\nRun Proof Tests\n")
  ocert.RunAllPTests(false)

  fmt.Printf("\nRun RMatrix Tests\n")
  ocert.RunAllRTests(false)

  //fmt.Println(ocert.TestEquation1ProofGen(true))

  // Scrap
  sharedParams := ocert.GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  //z := pairing.NewZr().Rand()
  //zP := pairing.NewG1().MulZn(g1, z)
  //
  //fmt.Println(z)
  //fmt.Println(zP)
  //
  //negG1 := pairing.NewG1().Invert(g1)
  //zprime := pairing.NewG1().Mul(zP, negG1)
  //
  //fmt.Println(negG1)
  //fmt.Println("ZPrime:",zprime)
  //
  //
  //z := pairing.NewZr().SetInt32(2)
  //fmt.Println(pairing.NewG1().Add(g1, g1) )
  //fmt.Println( pairing.NewG1().MulZn(g1, z))
  //fmt.Println(pairing.NewG1().Add(g1, g1).Equals(pairing.NewG1().MulZn(g1, z)))
}
