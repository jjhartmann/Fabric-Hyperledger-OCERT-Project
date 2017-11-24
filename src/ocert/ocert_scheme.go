/*
 * The main ocert scheme, it contains three protocl
 *  - Setup
 *  - GenECert
 *  - GenOCert
 */

package ocert

import (
 	"fmt"
)

func Put(stub Wrapper, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return nil, fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return []byte(args[1]), nil

}

func Get(stub Wrapper, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return nil, fmt.Errorf("Asset not found: %s", args[0])
	}
	return value, nil
}