# OCERT Hyperledger

## Structure
* **src**: Source code of our main contribution, including **ocert package** that implements the protocol described before, and **ocert chaincode** that is run in **Hyperledger Fabric** network.
    * **ocert**: Ocert package that implements ElGamal rerandomization encryption, structure-preserving signature and non-interactive zero knowledge proof system, as well as ocert main scheme.
    * ***chaincode/ocert.go***: The script that run in **Hyperledger Fabric** network, which is based on ***ocert\_scheme.go***. It implements `Init` and `Invoke` chaincode functions, and provides functions like `GenECert` and `GenOCert` to the client.
* **chaincode-docker-devmode**: This module relies on the **chaincode-docker-devmode** example in **Hyperledger Fabric** tutorial. It starts a docker environment that is used to test a chaincode, including a ***chaincode*** container where the **ocert chaincode** runs and a ***cli*** container where a client can interact with **ocert chaincode**. The docker-compose configuration file is modified so it first builds the images integrating **PBC** library for ***chaincode*** and ***cli*** containers before starting up the whole environment. During starting up, **ocert package** and **ocert chaincode** are copied into ***chaincode*** and ***cli*** containers. Once up, the docker environment is ready for running and testing **ocert chaincode**.
* **docker**: This module contains two ***Dockerfile***s that are used to build the images mentioned in docker-compose configuration file, one for ***chaincode*** container and another for ***cli*** container. 
* **benchmark**: This module contains scripts used in evaluation, including one script ***benchmarkcc.go***. It can run in ***cli*** container, interact with **ocert chaincode** and collect data for benchmark.
* **data**: This module contains three files recording data generated during benchmark. ***genOCertLog.txt*** records the total time to generate one **ocert**; ***genProofLog.txt*** records the time to generate a proof of knowledge in one **ocert** generation; ***verifyProofLog.txt*** records the time to verify a proof of knowledge in one **ocert** generation.
* **benchmark-analysis-tool**: This module contains the script used to evaluate benchmark date.

## Build and Run
To run **ocert chaincode**, please open three **Terminal**s and navigate to **chaincode-docker-devmode** in each **Terminal**
* **Terminal 1**: Start the network
    Setup the whole docker environment by
    ````
    docker-compose -f docker-compose-simple.yaml up
    ````
* **Terminal 2**: Build and start the chaincode
    First enter ***chaincode*** container by
    ````
    docker exec -it chaincode bash
    ````
    Then build **ocert chaincode**
    ````
    cd ocert
    go build ocert.go
    ````
    Last, start chaincode
    ````
    CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc:0 ./ocert
    ````
* **Terminal 3**: Use the chaincode
    First enter ***cli*** container by
    ````
    docker exec -it cli bash
    ````
    Then initialize chaincode
    ````
    peer chaincode install -p chaincodedev/chaincode/ocert -n mycc -v 0
    peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc
    ````
    Now you can play around **ocert chaincode**. For example, you can get the bilinear group used in the protocol
    ````
    peer chaincode query -n mycc -c '{"Args":["sharedParams"]}' -C myc
    ````
    or the public keys
    ````
    peer chaincode query -n mycc -c '{"Args":["get","auditor_pk"]}' -C myc
    peer chaincode query -n mycc -c '{"Args":["get","structure_preserving_vk"]}' -C myc
    peer chaincode query -n mycc -c '{"Args":["get","rsa_pk"]}' -C myc
    ````
    and definitiely, generate **ecert** and **ocert**S
    ````
    peer chaincode query -n mycc -c '{"Args":["genECert", arguments_used_by_GenECert]}' -C myc
    peer chaincode query -n mycc -c '{"Args":["genOCert", arguments_used_by_GenOCert]}' -C myc
    ````
    You should generate `arguments_used_by_GenECert` and `arguments_used_by_GenOCert`, and encode them by `Bytes()` from ***types.go***. Please refer ***benchmarkcc.go*** to use these functions. You also need to use the `SetBytes()` from ***types.go*** to decode the result from these functions.