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

/*
 * Based on "Optimal Structure-Preserving Signatures in Asymmetric Bilinear Groups"
 * by Abe and et al. The key size depends on Ocert paper.
 */

package ocert

import (
    "github.com/Nik-U/pbc"
)

/*
 * Generate key pair used by orgnization i.
 * SVerificationKey VK is used as ecert verification key VK_e,i for each
 * organization i
 * SSigningKey SK_e,i is used to sign client Pseudonym and ClientPublicKey
 * VK = (U, V, W1, W2, Z) in G1 * G2^4
 * SK = (u, v, w1, w2, z) where u, v, w1, w2, z is randomly picked from group of 
 * units modulo p
 * s.t. U = u * g1, V = v * g2, W1 = w1 * g2, W2 = w2 * g2 and Z = z * g2, where
 * g1 is the generator of group G1, and g2 is the generator of group G2
 */
func SKeyGen(sharedParams *SharedParams) (*SVerificationKey, *SSigningKey) {
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)
    g2 := pairing.NewG2().SetBytes(sharedParams.G2)

    VK := new(SVerificationKey)
    SK := new(SSigningKey)

    u := pairing.NewZr().Rand()
    v := pairing.NewZr().Rand()
    w1 := pairing.NewZr().Rand()
    w2 := pairing.NewZr().Rand()
    z := pairing.NewZr().Rand()
    
    SK.U = u.Bytes()
    SK.V = v.Bytes()
    SK.W1 = w1.Bytes()
    SK.W2 = w2.Bytes()
    SK.Z = z.Bytes()

    U := pairing.NewG1().MulZn(g1, u)
    V := pairing.NewG2().MulZn(g2, v)
    W1 := pairing.NewG2().MulZn(g2, w1)
    W2 := pairing.NewG2().MulZn(g2, w2)
    Z := pairing.NewG2().MulZn(g2, z)

    VK.U = U.Bytes()
    VK.V = V.Bytes()
    VK.W1 = W1.Bytes()
    VK.W2 = W2.Bytes()
    VK.Z = Z.Bytes()

    return VK, SK
}

/*
 * Signing the pseudonym and public key of a client by the signing key SK_e,i
 * of an organziation i. The output of signing procedure is the ecert.
 * The ecert has format (R, S, T) in G1 * G1 * G2, where g1, g2 are generators
 * of G1 and G2, r is randomly picked from group of units modulo p,
 * SKei = (u, v, w1, w2, z), P = (C, D)
 * R = r * g1
 * S = (z - r * v) * g1 + (-w1) * C + (-w2) * D
 * T = (1 / 6) * (g2 + (-u) * PKc)
 */
func SSign(sharedParams *SharedParams, SKei *SSigningKey, P *Pseudonym, PKc *ClientPublicKey) *Ecert {
    ecert := new(Ecert)
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)
    g2 := pairing.NewG2().SetBytes(sharedParams.G2)

    u := pairing.NewZr().SetBytes(SKei.U)
    v := pairing.NewZr().SetBytes(SKei.V)
    w1 := pairing.NewZr().SetBytes(SKei.W1)
    w2 := pairing.NewZr().SetBytes(SKei.W2)
    z := pairing.NewZr().SetBytes(SKei.Z)

    C := pairing.NewG1().SetBytes(P.C)
    D := pairing.NewG1().SetBytes(P.D)

    N := pairing.NewG2().SetBytes(PKc.PK)

    // Generate R
    r := pairing.NewZr().Rand()
    R := pairing.NewG1().MulZn(g1, r)
    ecert.R = R.Bytes()

    // Generate S
    s0 := pairing.NewZr().Mul(r, v)
    s0.Sub(z, s0)
    S0 := pairing.NewG1().MulZn(g1, s0)

    negW1 := pairing.NewZr().Neg(w1)
    S1 := pairing.NewG1().MulZn(C, negW1)

    negW2 := pairing.NewZr().Neg(w2)
    S2 := pairing.NewG1().MulZn(D, negW2)

    S := pairing.NewG1().Add(S0, S1)
    S.Add(S, S2)
    ecert.S = S.Bytes()

    // Generate T
    negU := pairing.NewZr().Neg(u)
    T0 := pairing.NewG2().MulZn(N, negU)
    T := pairing.NewG2().Add(g2, T0)

    invR := pairing.NewZr().Invert(r)
    T.MulZn(T, invR)
    ecert.T = T.Bytes()

    return ecert
}

/*
 * Verifying the signature is signed by the client, and returns a boolean, where
 * g1 and g2 are generators of group G1 and G2 respectively
 * VKei = (U, V, W1, W2, Z),
 * P = (C, D)
 * ecert = (R, S, T),
 * and to verify, test
 * e(R, V) * e(S, g2) * e(C, W1) * e(D, W2) = e(g1, Z) and
 * e(R, T) * e(U, PKc) = e(g1, g2), where e is the pairing operation
 */
func SVerify(sharedParams *SharedParams, VKei *SVerificationKey, P *Pseudonym, PKc *ClientPublicKey, ecert *Ecert) bool {
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)
    g2 := pairing.NewG2().SetBytes(sharedParams.G2)

    U := pairing.NewG1().SetBytes(VKei.U)
    V := pairing.NewG2().SetBytes(VKei.V)
    W1 := pairing.NewG2().SetBytes(VKei.W1)
    W2 := pairing.NewG2().SetBytes(VKei.W2)
    Z := pairing.NewG2().SetBytes(VKei.Z)

    R := pairing.NewG1().SetBytes(ecert.R)
    S := pairing.NewG1().SetBytes(ecert.S)
    T := pairing.NewG2().SetBytes(ecert.T)

    C := pairing.NewG1().SetBytes(P.C)
    D := pairing.NewG1().SetBytes(P.D)

    N := pairing.NewG2().SetBytes(PKc.PK)

    // Verify e1 * e2 * e3 * e4 = e5
    e1 := pairing.NewGT().Pair(R, V)
    e2 := pairing.NewGT().Pair(S, g2)
    e3 := pairing.NewGT().Pair(C, W1)
    e4 := pairing.NewGT().Pair(D, W2)
    e5 := pairing.NewGT().Pair(g1, Z)
    
    LHS1 := pairing.NewGT().Mul(e1, e2)
    LHS1.Mul(LHS1, e3)
    LHS1.Mul(LHS1, e4)

    if !LHS1.Equals(e5) {
        return false
    }

    // Verify e6 * e7 = e8
    e6 := pairing.NewGT().Pair(R, T)
    e7 := pairing.NewGT().Pair(U, N)
    e8 := pairing.NewGT().Pair(g1, g2)

    LHS2 := pairing.NewGT().Mul(e6, e7)

    if !LHS2.Equals(e8) {
        return false
    }

    return true
}