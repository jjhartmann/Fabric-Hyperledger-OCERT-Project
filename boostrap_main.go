package main

import (
  "fmt"
	 "./src/ocert"
  //"github.com/Nik-U/pbc"
)

func main() {
  //fmt.Println(ocert.Stest())
  //fmt.Println(ocert.IotaRhoTest(true))
  fmt.Println(ocert.TestIotaRhoPrime(true))

  //sharedParams := ocert.GenerateSharedParams()
  //pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  //g1 := pairing.NewG1().Rand()
  //g2 := pairing.NewG2().Rand()
  //gt := pairing.NewGT().Pair(g1, g2)
  //_ = gt
  //
  //
  //z := pairing.NewZr().SetInt32(2)
  //fmt.Println(pairing.NewG1().Add(g1, g1) )
  //fmt.Println( pairing.NewG1().MulZn(g1, z))
  //fmt.Println(pairing.NewG1().Add(g1, g1).Equals(pairing.NewG1().MulZn(g1, z)))
}
