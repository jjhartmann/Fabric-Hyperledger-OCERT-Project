/*
 * Used to generate ecert, the X.509 certificate that contains the 
 * signature on client public key and pseudonym
 */

package ocert

import (
    "fmt"
)

/*
 * Convert PKc|P to byte array
 */
func OCertSingedBytes(PKc *ClientPublicKey, P *Pseudonym) ([]byte, error) {
    var result []byte
    result = append(result[:], PKc.PK[:]...)
    Pbytes, err := P.Bytes()
    if err != nil {
        return nil, fmt.Errorf("Failed to convert signed message to bytes")
    }
    result = append(result[:], Pbytes[:]...)
    return result, nil
}