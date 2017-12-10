/*
 *
 * Copyright 2017 Kewei Shi, Jeremy Hartmann, Tuhin Tiwari and Dharvi Verma
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
  "fmt"
	 "../src/ocert"
  "github.com/Nik-U/pbc"
  time2 "time"
  "os"
  "encoding/csv"
  "strconv"
)

func main() {
  fmt.Printf("\nRun Structure Perserving Tests\n")
  fmt.Println(ocert.Stest())

  fmt.Printf("\nRun Proof Tests\n")
  ocert.RunAllPTests(false)

  fmt.Printf("\nRun RMatrix Tests\n")
  ocert.RunAllRTests(false)

  fmt.Printf("\nRun ElGamal ReReand Tests\n")
  ocert.ETestAll(false)

  //Test Key generation from rerandomization
  //fmt.Println(ocert.TestEquation5Verify(true))
  //fmt.Println(ocert.TestElementWiseSubtraction(true, 4, 4))
  //fmt.Println(ocert.Ptest(true))
  //fmt.Println(ocert.TestEquation3Verify(true))

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
  //
  ////zero := pairing.NewZr().Set0()
  ////one := pairing.NewZr().Set1()
  //rand := pairing.NewZr().Rand()
  //res := pairing.NewG1().MulZn(pairing.NewG1().Set0(), rand)
  //fmt.Println(res)

  //fmt.Println("G", g1)
  //Gz := pairing.NewG1().MulZn(g1, zero)
  //fmt.Println("Gz", Gz)
  //Go := pairing.NewG1().MulZn(g1, one)
  //fmt.Println("Go:", Go)
  //fmt.Println("Identiy0:", pairing.NewGT().Set0())
  //fmt.Println("Identiy:1", pairing.NewGT().Set1())


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
