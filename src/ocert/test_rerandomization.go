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
func Etest() bool {
	sharedParams := GenerateSharedParams()
	PK, SK := EKeyGen(sharedParams)

	// TODO randomly generate id
	id := new(ClientID)
	P := EEnc(sharedParams, PK, id)

	decryptedId := EDec(sharedParams, SK, P)

	if id != decryptedId {
		return false
	}

	PPrime := ERerand(sharedParams, PK, P)
	return ERerandVerify(sharedParams, P, PPrime)
}

/*
 * Run test b times
 */
func RunETest(b int) {
	for i := 0; i < b; i++ {
		if !Etest() {
			os.Exit(1)
		}
	}
}

func EGenKeyTest() bool{
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
  id2 := EDec(sharedParams, SK, P)
  if verbose {fmt.Println(id2)}


  return reflect.DeepEqual(id, id2)

}




