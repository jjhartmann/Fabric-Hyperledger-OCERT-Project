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
	pi.Print()

	// Create constants for verify
	consts := new(ProofConstants)
	consts.VK = VK
	consts.PKa = PKa
	consts.Egh = pairing.NewGT().Pair(G, H).Bytes()
	consts.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VK.Z)).Bytes()
	consts.PPrime = Pprime

	consts.Print()

	result := PProve(sharedParams, pi, consts)
	fmt.Println(result)
}