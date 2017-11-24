/*
 * A wrapper of 
 */

package ocert

type Wrapper interface {
	GetState(key string) ([]byte, error)
	PutState(key string, value []byte) error
}