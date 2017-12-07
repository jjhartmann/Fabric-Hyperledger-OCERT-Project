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
	// pi.Print()
	// pi2.Print()
	// pi23.Print()
	fmt.Println(pi.Equals(pi2))
	fmt.Println(pi.Equals(pi3))

	fmt.Println("-----------")
	crs := pi.sigma
	fmt.Println(crs)
	crsBytes, err := crs.Bytes()
	fmt.Println(err)
	fmt.Println(string(crsBytes))
	crs2 := new(CommonReferenceString)
	err = crs2.SetBytes(crsBytes)
	fmt.Println(err)
	fmt.Println(crs.Equals(crs2))
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
	
	// consts.Print()
	// consts2.Print()
	// fmt.Println(consts.Equals(consts2))

	consts3 := new(ProofConstants)
	consts3.VK = VK
	consts3.PKa = PKa
	consts3.Egh = pairing.NewGT().Pair(G, H).Bytes()
	consts3.Egz = pairing.NewGT().Pair(G, pairing.NewG2().SetBytes(VK.Z)).Bytes()
	consts3.PPrime = P
	// consts3.Print()
	// fmt.Println(consts.Equals(consts3))
	// fmt.Println(consts2.Equals(consts3))

	constsBytes, err := consts.Bytes()
	fmt.Println(err)
	fmt.Println(constsBytes)
	consts4 := new(ProofConstants)
	err = consts4.SetBytes(constsBytes)
	// consts4.Print()
	fmt.Println(err)
	fmt.Println(consts.Equals(consts4))

	// result := PProve(sharedParams, pi, consts)
	// fmt.Println(result)
	// result3 := PProve(sharedParams, pi3, consts)
	// fmt.Println(result3)

}