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
 *  - genECert
 *  - genOCert
 * and
 *  - sharedParams
 *  - get
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
    } else if fn == "sharedParams" {
        result, err = ocert.GetSharedParams(stub, args)
    } else if fn == "auditorKeypair" {
        result, err = ocert.GetAuditorKeypair(stub, args)
    } else if fn == "genECert" {
        result, err = ocert.GenECert(stub, args)
    } else if fn == "genOCert" {
        result, err = ocert.GenOCert(stub, args)
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
