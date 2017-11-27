/*
 * Used to generate ecert, the X.509 certificate that contains the 
 * signature on client public key and pseudonym
 */

package ocert

import (
 	"fmt"
	"math/big"
	"time"
	// "crypto"
	"crypto/rsa"
	"crypto/rand"
	// "crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
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

// TODO
/*
 * Sign PKc|P by RSA and generate the X.509 certificate that contains the
 * signature, which is the ocert
 */
func GenOCertHelper(
	PKc *ClientPublicKey, 
	P *Pseudonym, 
	privateKey *rsa.PrivateKey,
	serialNumber *big.Int) ([]byte, error) {

	msg, err := OCertSingedBytes(PKc, P)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		SignatureAlgorithm: x509.SHA256WithRSA,

		PublicKeyAlgorithm: x509.RSA,
		PublicKey: &privateKey.PublicKey,

  		SerialNumber: serialNumber,
  		
  		Issuer: pkix.Name {
  			CommonName: "Test Ocert CA",
  		},
  		Subject: pkix.Name {
  			CommonName: "Test client",
  		},

  		NotBefore: notBefore,
  		NotAfter: notAfter,

  		SubjectKeyId: msg,
  	}

  	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return derBytes, nil
}