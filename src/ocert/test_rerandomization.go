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

package ocert

import (
    "os"
    "github.com/Nik-U/pbc"
    "reflect"
    "fmt"
)

/*
 * Run a single test
 */
func Etest(verbose bool) bool {
    sharedParams := GenerateSharedParams()
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    PK, SK := EKeyGen(sharedParams)

    // Random id in G1 (will be calcualted from hyperledger)
    if verbose {fmt.Println("Generate Client ID")}
    id := new(ClientID)
    id.ID = pairing.NewG1().Rand().Bytes()
    if verbose {fmt.Println(id)}


    // Generate pseudonym
    if verbose { fmt.Println("Generate Pseudonym")}
    P := EEnc(sharedParams, PK, id)
    if verbose {fmt.Println(P)}

    // Decrypt client id
    if verbose { fmt.Println("Decrypt id")}
    decryptID := EDec(sharedParams, SK, P)
    if verbose {fmt.Println(decryptID)}


    if !reflect.DeepEqual(id, decryptID) {
        if verbose {fmt.Println("Deep Equal Failed on ID")}
        return false
    }

    // Generate Rerand
    if verbose {fmt.Println("Genearte Rerand")}
    Pprime, _ := ERerand(sharedParams, PK, P)
    if verbose {fmt.Println(Pprime)}

    return ERerandVerify(sharedParams, SK, P, Pprime)
}

/*
 * Run test b times
 */
func RunETest(b int) {
    for i := 0; i < b; i++ {
        if !Etest(false) {
            os.Exit(1)
        }
    }
}

func EGenKeyTest(verbose bool) bool{
    sharedParams := GenerateSharedParams()
    PK, SK := EKeyGen(sharedParams)

    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)

    SKa:=pairing.NewZr().SetBytes(SK.SK)
    PKa:=pairing.NewG1().MulZn(g1,SKa)

    return PKa.Equals(pairing.NewG1().SetBytes(PK.PK))

}

func  ETestEncDec(verbose bool) bool {
    sharedParams := GenerateSharedParams()
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    PK, SK := EKeyGen(sharedParams)

    // Random id in G1 (will be calcualted from hyperledger)
    if verbose {fmt.Println("Generate Client ID")}
    id := new(ClientID)
    id.ID = pairing.NewG1().Rand().Bytes()
    if verbose {fmt.Println(id)}


    // Generate pseudonym
    if verbose { fmt.Println("Generate Pseudonym")}
    P := EEnc(sharedParams, PK, id)
    if verbose {fmt.Println(P)}

    // Decrypt client id
    if verbose { fmt.Println("Decrypt id")}
    decryptID := EDec(sharedParams, SK, P)
    if verbose {fmt.Println(decryptID)}


    return reflect.DeepEqual(id, decryptID)
}


func ETestRerandVerify(verbose bool) bool {
    sharedParams := GenerateSharedParams()
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    PK, SK := EKeyGen(sharedParams)

    // Random id in G1 (will be calcualted from hyperledger)
    if verbose {fmt.Println("Generate Client ID")}
    id := new(ClientID)
    id.ID = pairing.NewG1().Rand().Bytes()
    if verbose {fmt.Println(id)}


    // Generate pseudonym
    if verbose { fmt.Println("Generate Pseudonym")}
    P := EEnc(sharedParams, PK, id)
    if verbose {fmt.Println(P)}

    // Generate Rerand
    if verbose {fmt.Println("Genearte Rerand")}
    Pprime, _ := ERerand(sharedParams, PK, P)
    if verbose {fmt.Println(Pprime)}

    if verbose {fmt.Println("Verify Rerand")}
    return ERerandVerify(sharedParams, SK, P, Pprime)
}


func ETestAll(verbose bool) {
    fmt.Println("KeyGen:         ", EGenKeyTest(verbose))
    fmt.Println("Enc and Dec:    ", ETestEncDec(verbose))
    fmt.Println("Rerand Verify:  ", ETestRerandVerify(verbose))
}