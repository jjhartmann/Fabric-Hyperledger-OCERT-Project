# Structure 
* **ocert**: Ocert package that implements ElGamal rerandomization encryption, structure-preserving signature and non-interactive zero knowledge proof system, as well as ocert main scheme.
    * **types.go**: This file contains all the `struct`s used in our protocol, as well as printing(`Print()`), equality(`Equals()`), encoding(`Bytes()`) and decoding(`SetBytes()`) functions for `struct`s.
    * **test\_types.go**: Test for printing, equality, encoding and decoding function for proof of knowledge $$\pi$$.
    * **params.go**: `GenerateSharedParams()` in this file is used to generate the bilinear group used through all schemes.
    * **proof.go**: 
    * **test\_proof.go**:
    * **rmatrix.go**:
    * **test\_rmatrix.go**:
    * **rerandomization.go**:
    * **test\_rerandomization.go**:
    * **structure\_preserving.go**: The structure-preserving signature is implemented in this file, including three algorithms `SKeyGen()`, `SSign()` and `SVerify()`.
    * **test\_structure\_preserving.go**: Test for structure-preserving signature scheme.
    * **stub\_wrapper.go**: An abstract interface that wraps **Hyperledger Fabric** `shim. ChaincodeStubInterface`, so we can test our main protocol without starting the network.
    * **certificate.go**: The function used to generate X.509 certificate should be included in this file in future. 
    * **ocert\_scheme.go**: Main protocol, including three main algorithms `Setup()`, `GenECert()` and `GenOCert()`, and some helper functions like `Get()`, `GetSharedParams()` and `GetAuditorKeypair()`.
* ***chaincode/ocert.go***: The script that run in **Hyperledger Fabric** network, which is based on ***ocert\_scheme.go***. It implements `Init` and `Invoke` chaincode functions, and provides functions like `GenECert` and `GenOCert` to the client.