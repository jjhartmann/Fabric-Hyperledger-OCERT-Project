/*
 * Chaincode to generate ocerts
 * This chaincode is only used in benchmark
 */

package main

import (
	"fmt"
	"ocert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type OcertAsset struct {
}

/*
 * Call ocert.Setup
 */
func (t *OcertAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetArgs()
	result, err := ocert.Setup(stub, args)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

/*
 * ocert chaincode provides the following functions
 *  -
 */
func (t *OcertAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	allargs := stub.GetArgs()
	if len(allargs) < 1 {
		return shim.Error("Incorrect arguments")
	}

	fn := ""
	args := [][]byte{}
	if len(allargs) >= 1 {
		fn = string(allargs[0])
		args = allargs[1:]
	}

	var result []byte
	var err error
	if fn == "get" {
		result, err = ocert.Get(stub, args)
	} else if fn == "put" {
		result, err = ocert.Put(stub, args)
	} else if fn == "sharedParams" {
		result, err = ocert.GetSharedParams(stub, args)
	} else {
		return shim.Error("Unknown functions")
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success(result)
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(OcertAsset)); err != nil {
		fmt.Printf("Error starting OcertAsset chaincode: %s", err)
	}
}
