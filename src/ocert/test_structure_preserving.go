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
    "fmt"
    "github.com/Nik-U/pbc"
    "math/rand"
)

/*
 * Run a single test
 */
// TODO add benchmark in future
func Stest() bool {
    fmt.Println("[Structure Preserving] Start test")
    sharedParams := GenerateSharedParams()
    VK, SK := SKeyGen(sharedParams)

    P := new(Pseudonym)
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    P.C = pairing.NewG1().Rand().Bytes()
    P.D = pairing.NewG1().Rand().Bytes()
    PKc := new(ClientPublicKey)
    PKc.PK = pairing.NewG2().Rand().Bytes()
    ecert := SSign(sharedParams, SK, P, PKc)

    if SVerify(sharedParams, VK, P, PKc, ecert) {
        fmt.Println("[Structure Preserving] Verify a ecert successfully")
    } else {
        fmt.Println("[Structure Preserving] Cannot verify a ecert")
        return false
    }

    seed := rand.Intn(3)
    if seed == 0 {
        fmt.Println("[Structure Preserving] Modify C")
        P.C = pairing.NewG1().Rand().Bytes()
    } else if seed == 1 {
        fmt.Println("[Structure Preserving] Modify D")
        P.D = pairing.NewG1().Rand().Bytes()
    } else {
        fmt.Println("[Structure Preserving] Modify PKc")
        PKc.PK = pairing.NewG2().Rand().Bytes()
    }
    
    if !SVerify(sharedParams, VK, P, PKc, ecert) {
        fmt.Println("[Structure Preserving] Reject a ecert correctly")
    } else {
        fmt.Println("[Structure Preserving] Fail to reject a false ecert")
        return false
    }

    fmt.Println("[Structure Preserving] Pass test")
    return true
}

/*
 * Run test b times
 */
func RunSTest(b int) {
    for i := 0; i < b; i++ {
        if !Stest() {
            os.Exit(1)
        }
    } 
}