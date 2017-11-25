package main

import (
  "fmt"
	 "./src/ocert"
	"github.com/Nik-U/pbc"
)

func main() {
  //fmt.Println(ocert.Stest())
  //fmt.Println(ocert.IotaRhoTest(true))

  sharedParams := ocert.GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  alpha := pairing.NewZr().Rand()
  sigma := ocert.CreateCommonReferenceString(sharedParams, alpha)
  _ = sigma

  // Test IotaPrim: Zp -> B1
  z := pairing.NewZr().Rand()
  B := ocert.IotaPrime1(pairing, z, sigma)
  fmt.Println(B)


}
