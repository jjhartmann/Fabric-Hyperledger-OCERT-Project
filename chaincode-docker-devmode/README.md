# Build and Run
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