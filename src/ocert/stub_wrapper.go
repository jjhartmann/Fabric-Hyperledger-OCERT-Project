/*
 * A wrapper of shim.ChaincodeStubInterface, so we
 * can test part of chaincode locally, without starting
 * the whole Hyperledger Fabric network
 */

package ocert

import (
    "fmt"
)

type Wrapper interface {
    GetState(key string) ([]byte, error)
    PutState(key string, value []byte) error
}

func Put(stub Wrapper, args [][]byte) ([]byte, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("Incorrect arguments. Expecting a key and a value")
    }

    err := stub.PutState(string(args[0]), args[1])
    if err != nil {
        return nil, fmt.Errorf("Failed to set asset: %s", args[0])
    }
    return []byte(args[1]), nil

}