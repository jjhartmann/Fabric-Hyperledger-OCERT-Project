package main

import (
  "fmt"
	 "./src/ocert"
  "github.com/Nik-U/pbc"
  time2 "time"
  "os"
  "encoding/csv"
  "strconv"
)

func main() {
  //fmt.Printf("\nRun Structure Perserving Tests\n")
  //fmt.Println(ocert.Stest())

  //fmt.Printf("\nRun Proof Tests\n")
  //ocert.RunAllPTests(false)

  fmt.Println("RMatrix Mult in G1 2x2 2x2     ", ocert.TestRMatrixMultiplicationforElementinG1(true,2, 2, 2, 2))

  //fmt.Printf("\nRun RMatrix Tests\n")
  //ocert.RunAllRTests(false)

  //fmt.Println(ocert.TestEquation1Verify(true))

  // Benchmark
  //ConstructMetricsForProofVerifyEq1(100)

  // Scrap
  //sharedParams := ocert.GenerateSharedParams()
  //pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  //g1 := pairing.NewG1().Rand()
  //g2 := pairing.NewG2().Rand()
  //gt := pairing.NewGT().Pair(g1, g2)
  //_ = gt
  //
  //x := pairing.NewZr().Rand()
  //H := pairing.NewG2().Rand()
  //PK := pairing.NewG2().MulZn(H, x)
  //fmt.Println(PK)
  //
  //identity := pairing.NewZr().Set1()
  //fmt.Println(identity)
  //
  //negId := pairing.NewZr().Neg(identity)
  //fmt.Println(negId)
  //
  //negPK := pairing.NewG2().Neg(PK)
  //fmt.Println(negPK)
  //fmt.Println(pairing.NewG2().Add(PK, negPK))
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

func ConstructMetricsForProofVerifyEq1(n int) {
  sharedParams := ocert.GenerateSharedParams()
  pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
  g1 := pairing.NewG1().Rand()
  g2 := pairing.NewG2().Rand()
  gt := pairing.NewGT().Pair(g1, g2)
  _ = gt

  // Create csv file
  csvfile, err := os.Create("Analysis/result_3.csv")
  if err != nil {
    fmt.Println(err)
    return
  }
  defer csvfile.Close()

  // Create Write
  writer := csv.NewWriter(csvfile)

  // Header for csv
  header := []string{"Sequence", "Function", "Time"}
  writer.Write(header)

  fmt.Printf("\nStarting Benchmark: Proof Generator\n")
  for i := 0; i < n; i++ {
    tmp := time2.Now()
    ocert.TestEquation1ProofGen(false)
    tmpe := time2.Since(tmp)
    fmt.Println("Time: ", i, tmpe.Seconds())
    writer.Write([]string{strconv.Itoa(i), "Proof", strconv.FormatFloat(tmpe.Seconds(), 'f', 6, 64)})
  }

  fmt.Printf("\nStarting Benchmark: Proof and Verify\n")
  for i := 0; i < n; i++ {
    tmp := time2.Now()
    ocert.TestEquation1Verify(false)
    tmpe := time2.Since(tmp)
    fmt.Println("Time: ", i, tmpe.Seconds())
    writer.Write([]string{strconv.Itoa(i), "VerifyProof", strconv.FormatFloat(tmpe.Seconds(), 'f', 6, 64)})
  }

  writer.Flush()
}
