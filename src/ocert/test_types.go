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

import(
    "fmt"
    "github.com/Nik-U/pbc"
)

func RunTypesTest() {
    fmt.Println("RunTypesTest")

    sharedParams := GenerateSharedParams()
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

    // Shared generators
    G := pairing.NewG1().SetBytes(sharedParams.G1)
    H := pairing.NewG2().SetBytes(sharedParams.G2)

    // Verifcation key
    VK, SK := SKeyGen(sharedParams)

    // Auditor Key pair
    PKa, _ := EKeyGen(sharedParams)

    // Pseudonym Generation
    clientID := new(ClientID)
    clientID.ID = pairing.NewG1().Rand().Bytes()
    P := EEnc(sharedParams, PKa, clientID)

    // Rerandomization
    Pprime, rprime := ERerand(sharedParams, PKa, P)

    // Client Keypair
    PKc := new(ClientPublicKey)
    Xc := pairing.NewZr().Rand().Bytes()
    PKc.PK = pairing.NewG2().MulZn(H, pairing.NewZr().SetBytes(Xc)).Bytes()

    // Ecert Generation
    ecert := SSign(sharedParams, SK, P, PKc)

    // Construct Var
    vars := new(ProofVariables)
    vars.PKa = PKa
    vars.P = P
    vars.VK = VK
    vars.RPrime = rprime
    vars.PKc = PKc
    vars.Xc = Xc
    vars.E = ecert

    pi := PSetup(sharedParams, vars)
    pi2 := pi
    pi3 := PSetup(sharedParams, vars)
    pi.Print()
    pi2.Print()
    pi3.Print()
    fmt.Println(pi.Equals(pi2))
    fmt.Println(pi.Equals(pi3))

    fmt.Println("-----------")
    crs := pi.sigma
    crs.Print()
    crsBytes, err := crs.Bytes()
    fmt.Println(err)
    fmt.Println(crsBytes)
    crs2 := new(CommonReferenceString)
    err = crs2.SetBytes(crsBytes)
    fmt.Println(err)
    fmt.Println(crs.Equals(crs2))
    fmt.Println("-----------")
    eq := pi.Eq1
    eq.Print()
    eqBytes, err := eq.Bytes()
    fmt.Println(err)
    fmt.Println(eqBytes)
    eq2 := new(ProofOfEquation)
    err = eq2.SetBytes(eqBytes)
    fmt.Println(err)
    fmt.Println(eq.Equals(eq2))
    fmt.Println("-----------")
    piBytes, err := pi.Bytes()
    fmt.Println(err)
    fmt.Println(piBytes)
    pi4 := new(ProofOfKnowledge)
    err = pi4.SetBytes(piBytes)
    fmt.Println(err)
    fmt.Println(pi.Equals(pi4))
    fmt.Println("-----------")

    // Create constants for verify
    consts := new(ProofConstants)
    consts.VK = VK
    consts.PKa = PKa
    consts.Egh = pairing.NewGT().Pair(G, H).Bytes()
    consts.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VK.Z)).Bytes()
    consts.PPrime = Pprime

    consts2 := new(ProofConstants)
    consts2.VK = VK
    consts2.PKa = PKa
    consts2.Egh = pairing.NewGT().Pair(G, H).Bytes()
    consts2.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VK.Z)).Bytes()
    consts2.PPrime = Pprime
    
    consts.Print()
    consts2.Print()
    fmt.Println(consts.Equals(consts2))

    consts3 := new(ProofConstants)
    consts3.VK = VK
    consts3.PKa = PKa
    consts3.Egh = pairing.NewGT().Pair(G, H).Bytes()
    consts3.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VK.Z)).Bytes()
    consts3.PPrime = P
    consts3.Print()
    fmt.Println(consts.Equals(consts3))
    fmt.Println(consts2.Equals(consts3))

    constsBytes, err := consts.Bytes()
    fmt.Println(err)
    fmt.Println(constsBytes)
    consts4 := new(ProofConstants)
    err = consts4.SetBytes(constsBytes)
    consts4.Print()
    fmt.Println(err)
    fmt.Println(consts.Equals(consts4))
}