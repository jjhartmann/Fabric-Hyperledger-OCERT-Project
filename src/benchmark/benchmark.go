/*
 * Used to benchmark ocert_scheme
 */

package main

import (
	"fmt"
	"ocert"
	"github.com/Nik-U/pbc"
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
	setupArgs := [][]byte{}
	_, err := ocert.Setup(db, setupArgs)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	auditorPKBytesKey := []byte("auditor_pk")
	getArgs := [][]byte{auditorPKBytesKey}
	auditorPKBytes, err := ocert.Get(db, getArgs)
	auditorPK := new(ocert.AuditorPublicKey)
	err = auditorPK.SetBytes(auditorPKBytes)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	fmt.Printf("[Benchmark] auditor_pk: ")
	fmt.Println(auditorPK)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

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

	IDc := new(ocert.ClientID)
	IDc.ID = pairing.NewG1().Rand().Bytes()
	fmt.Printf("[Benchmark] IDc: ")
	fmt.Println(IDc)

	PKc := new(ocert.ClientPublicKey)
	PKc.PK = pairing.NewG1().Rand().Bytes()
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
}