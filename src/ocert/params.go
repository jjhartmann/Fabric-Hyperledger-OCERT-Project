package ocert

import (
    "github.com/Nik-U/pbc"
)

/*
 * Randomly generate a paired group and the corresponding
 * generator of each group
 */
func GenerateSharedParams() *SharedParams {
    sharedParams := new(SharedParams)
    params := pbc.GenerateF(160)
    pairing := params.NewPairing()
    g1 := pairing.NewG1().Rand()
    g2 := pairing.NewG2().Rand()
    sharedParams.Params = params.String()
    sharedParams.G1 = g1.Bytes()
    sharedParams.G2 = g2.Bytes()
    return sharedParams
}
