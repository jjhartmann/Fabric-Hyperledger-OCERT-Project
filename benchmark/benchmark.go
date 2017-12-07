/*
 * Used to benchmark ocert_scheme
 */

package main

import (
	"fmt"
	"ocert"
	"github.com/Nik-U/pbc"
	"crypto"
	"crypto/x509"
	"crypto/sha256"
	"crypto/rsa"
)

type DB struct {
	DB map[string][]byte
}

func (db *DB) GetState(key string) ([]byte, error) {
	val, exist := db.DB[key]
	if exist {
		return val, nil
	} else {
		return nil, fmt.Errorf("Failed to get state")
	}
}

func (db *DB) PutState(key string, value []byte) error {
	db.DB[key] = value
	return nil
}

func main() {
	db := new(DB)
	db.DB = make(map[string][]byte)

	// Benchmark starts here

	// Setup
	setupArgs := [][]byte{}
	_, err := ocert.Setup(db, setupArgs)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	// Keys
	auditorPKBytesKey := []byte("auditor_pk")
	getArgs := [][]byte{auditorPKBytesKey}
	auditorPKBytes, err := ocert.Get(db, getArgs)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	auditorPK := new(ocert.AuditorPublicKey)
	err = auditorPK.SetBytes(auditorPKBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("[Benchmark] auditor_pk: ")
	fmt.Println(auditorPK)

	rsaPKBytesKey := []byte("rsa_pk")
	getArgs = [][]byte{rsaPKBytesKey}
	rsaPKBytes, err := ocert.Get(db, getArgs)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	rsaPKWrapper := new(ocert.RSAPK)
	err = rsaPKWrapper.SetBytes(rsaPKBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	rsaPK, err := x509.ParsePKIXPublicKey(rsaPKWrapper.PK)
	fmt.Printf("[Benchmark] rsa_pk: ")
	fmt.Println(rsaPK)

	sVKBytesKey := []byte("structure_preserving_vk")
	getArgs = [][]byte{sVKBytesKey}
	sVKBytes, err := ocert.Get(db, getArgs)
		if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	sVK := new(ocert.SVerificationKey)
	err = sVK.SetBytes(sVKBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("[Benchmark] sVK: ")
	fmt.Println(sVK)

	// Biliear groups
	sharedParamsBytes, err := ocert.GetSharedParams(db, setupArgs)
	sharedParams := new(ocert.SharedParams)
	err = sharedParams.SetBytes(sharedParamsBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("[Benchmark] sharedParams: ")
	fmt.Println(sharedParams)
	pairing, _ := pbc.NewPairingFromString(sharedParams.Params)
	H := pairing.NewG2().SetBytes(sharedParams.G2)

	// GenECert
	IDc := new(ocert.ClientID)
	IDc.ID = pairing.NewG1().Rand().Bytes()
	fmt.Printf("[Benchmark] IDc: ")
	fmt.Println(IDc)

	PKc := new(ocert.ClientPublicKey)
	Xc := pairing.NewZr().Rand().Bytes()
	PKc.PK = pairing.NewG2().MulZn(H, pairing.NewZr().SetBytes(Xc)).Bytes()
	fmt.Printf("[Benchmark] PKc: ")
	fmt.Println(PKc)

	ecertRequest := new(ocert.GenECertRequest)
	ecertRequest.IDc = IDc.ID
	ecertRequest.PKc = PKc.PK
	ecertRequestBytes, err := ecertRequest.Bytes()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	ecertArgs := [][]byte{ecertRequestBytes}

	ecertReplyBytes, err := ocert.GenECert(db, ecertArgs)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	ecertReply := new(ocert.GenECertReply)
	err = ecertReply.SetBytes(ecertReplyBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	P := new(ocert.Pseudonym)
	err = P.SetBytes(ecertReply.P)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("[Benchmark] P: ")
	fmt.Println(P)

	ecert := new(ocert.Ecert)
	err = ecert.SetBytes(ecertReply.Ecert)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("[Benchmark] ecert: ")
	fmt.Println(ecert)

	// GenOCert
	newPKc := new(ocert.ClientPublicKey)
	newPKc.PK = pairing.NewG1().Rand().Bytes()
	fmt.Printf("[Benchmark] newPKc: ")
	fmt.Println(newPKc)

	newP, rprime := ocert.ERerand(sharedParams, auditorPK, P)
	fmt.Printf("[Benchmark] newP: ")
	fmt.Println(newP)
	fmt.Printf("[Benchmark] rprime: ")
	fmt.Println(rprime)

	// Proof generation
	vars := new(ocert.ProofVariables)
	vars.PKa = auditorPK
	vars.P = P
	vars.VK = sVK
	vars.RPrime = rprime
	vars.PKc = PKc
	vars.Xc = Xc
	vars.E = ecert

	pi := ocert.PSetup(sharedParams, vars)

	ocertRequest := new(ocert.GenOCertRequest)
	ocertRequest.PKc = newPKc.PK
	ocertRequest.P, err = newP.Bytes()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	ocertRequest.Pi, err = pi.Bytes()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	ocertRequestBytes, err := ocertRequest.Bytes()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	ocertArgs := [][]byte{ocertRequestBytes}

	ocertReplyBytes, err := ocert.GenOCert(db, ocertArgs)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	ocertReply := new(ocert.GenOCertReply)
	err = ocertReply.SetBytes(ocertReplyBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	signature := ocertReply.Sig
	fmt.Printf("[Benchmark] signature: ")
	fmt.Println(signature)

	// Verify signature
	msg, err := ocert.OCertSingedBytes(newPKc, newP)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	hashed := sha256.Sum256(msg)
	err = rsa.VerifyPKCS1v15(rsaPK.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	} else {
		fmt.Println("[Benchmark] ocert verified")
	}
}