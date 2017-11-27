package main

import (
  "fmt"
	 "./src/ocert"
  "github.com/Nik-U/pbc"
  time2 "time"
)

func main() {
  fmt.Printf("\nRun Structure Perserving Tests\n")
  fmt.Println(ocert.Stest())

  fmt.Printf("\nRun Proof Tests\n")
  ocert.RunAllPTests(false)

  fmt.Printf("\nRun RMatrix Tests\n")
  ocert.RunAllRTests(false)

  // Benchmark
  BenchMarkEq1(10)

  // Scrap
  //sharedParams := ocert.GenerateSharedParams()
  //pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  //g1 := pairing.NewG1().Rand()
  //g2 := pairing.NewG2().Rand()
  //gt := pairing.NewGT().Pair(g1, g2)
  //_ = gt

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


func BenchMarkEq1(n int) {
  sharedParams := ocert.GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  fmt.Printf("\nStarting Benchmark: Proof Generator\n")
  start := time2.Now()
  for i := 0; i < n; i++ {
    tmp := time2.Now()
    ocert.TestEquation1ProofGen(false)
    tmpe := time2.Since(tmp)
    fmt.Println("Time: ", i, tmpe.Seconds())
  }
  elapsed := time2.Since(start)
  avgtimeProofGen := elapsed.Seconds()/float64(n)
  fmt.Println("Avg Time: ", avgtimeProofGen)


  fmt.Printf("\nStarting Benchmark: Proof and Verify\n")
  start = time2.Now()
  for i := 0; i < n; i++ {
    tmp := time2.Now()
    ocert.TestEquation1Verify(false)
    tmpe := time2.Since(tmp)
    fmt.Println("Time: ", i, tmpe.Seconds())
  }
  elapsed = time2.Since(start)
  avgtimeProofVerify := elapsed.Seconds()/float64(n)
  fmt.Println("Avg Time: ", avgtimeProofVerify)

  fmt.Printf("\nSummary Statistics:\n")
  fmt.Println("Proof Generation:    ", avgtimeProofGen)
  fmt.Println("Verify Proof:        ", avgtimeProofVerify - avgtimeProofGen)
  fmt.Println("Total:               ", avgtimeProofVerify)


}
