/*
 * Based on elGamal encryption on ECC
 */

package ocert

import (
    "github.com/Nik-U/pbc"
)

/*
 * Generate the pair of public key and secret key used by auditor.
 */
func EKeyGen(sharedParams *SharedParams) (*AuditorPublicKey, *AuditorSecretKey) {
    PKa := new(AuditorPublicKey)
    SKa := new(AuditorSecretKey)

    // Generators g1 & g2 are generated from the groups G1 & G2
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)


    xa :=pairing.NewZr().Rand()
    SKa.SK=xa.Bytes()

    //produce the public key
    PublicKa :=pairing.NewG1().MulZn(g1,xa)
    PKa.PK = PublicKa.Bytes()

    return PKa, SKa
}

/*
 * Encrypt the client id based on the auditor's public key, the
 * result is the pseudonym of a client, where pseudonym of a client
 * has form (C, D), where both C and D are in G1
 */
func EEnc(sharedParams *SharedParams, PKa *AuditorPublicKey, id *ClientID) *Pseudonym {
    P := new(Pseudonym)

    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)
    r := pairing.NewZr().Rand()

    C := pairing.NewG1().MulZn(g1,r)

    PK := pairing.NewG1().SetBytes(PKa.PK)
    D := pairing.NewG1().MulZn(PK,r)

    //Convert the id into an element in G1
    Cid := pairing.NewG1().SetBytes(id.ID)

    //Add the id value to D
    D = pairing.NewG1().Add(D,Cid)

    P.C = C.Bytes()
    P.D = D.Bytes()

    return P
}

/*
 * Decrypt the client real identiy based on the pseudonym of a client
 */
func EDec(sharedParams *SharedParams, SKa *AuditorSecretKey, P *Pseudonym) *ClientID {
    id := new(ClientID)
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

    C := pairing.NewG1().SetBytes(P.C)
    D := pairing.NewG1().SetBytes(P.D)
    x := pairing.NewZr().SetBytes(SKa.SK)

    xC := pairing.NewG1().MulZn(C, x)
    id.ID = pairing.NewG1().Sub(D, xC).Bytes()
    return id
}

/*
 * Rerandomize the client's pseudonym. Given a pseudonym P = (C, D)
 * of a client, this scheme can rerandomize it to a new pseudonym
 * P' = (C', D'), where P' is also in G1 * G1.
 */
func ERerand(sharedParams *SharedParams, PKa *AuditorPublicKey, P *Pseudonym) (*Pseudonym, []byte) {
    // TODO rerandomize P
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
    g1 := pairing.NewG1().SetBytes(sharedParams.G1)

    PPrime := new(Pseudonym)

    //Generate rprime
    rprime := pairing.NewZr().Rand()

    //Getting C & D from the Pseudonym P that has been passed
    C := pairing.NewG1().SetBytes(P.C)
    D := pairing.NewG1().SetBytes(P.D)

    //Multiplying rprime with the generator g1
    tempC := pairing.NewG1().MulZn(g1,rprime)

    //Adding the product of rprime & g1 and C (from Pseudonym P) to Cprime
    Cprime := pairing.NewG1().Add(C,tempC)

    //To find Dprime, using D from Pseudonym & product of rprime & Public Key
    PK := pairing.NewG1().SetBytes(PKa.PK)
    tempD := pairing.NewG1().MulZn(PK,rprime)
    Dprime := pairing.NewG1().Add(D,tempD)

    PPrime.C = Cprime.Bytes()
    PPrime.D = Dprime.Bytes()
    return PPrime, rprime.Bytes()
}

/*
 * Given two pseudonyms P and P', validate whether P' is rerandomized
 * from P
 */
func ERerandVerify(sharedParams *SharedParams, SKa *AuditorSecretKey, P *Pseudonym, PPrime *Pseudonym) bool {
    pairing, _ := pbc.NewPairingFromString(sharedParams.Params)

    // Retrieve Original P
    C := pairing.NewG1().SetBytes(P.C)
    D := pairing.NewG1().SetBytes(P.D)
    x := pairing.NewZr().SetBytes(SKa.SK)

    // Retrieve P prime
    Cprime := pairing.NewG1().SetBytes(PPrime.C)
    Dprime := pairing.NewG1().SetBytes(PPrime.D)

    rG := pairing.NewG1().Sub(Cprime, C)
    xrG := pairing.NewG1().MulZn(rG, x)

    return xrG.Equals(pairing.NewG1().Sub(Dprime, D))
}
