docker-compose -f docker-compose-simple.yaml up

docker exec -it chaincode bash
	CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc:0 ./ocert

docker exec -it cli bash
	peer chaincode install -p chaincodedev/chaincode/ocert -n mycc -v 0
	peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

	peer chaincode query -n mycc -c '{"Args":["get","structure_preserving_vk"]}' -C myc
	peer chaincode query -n mycc -c '{"Args":["get","rsa_pk"]}' -C myc
	peer chaincode query -n mycc -c '{"Args":["get","auditor_pk"]}' -C myc
	peer chaincode query -n mycc -c '{"Args":["sharedParams"]}' -C myc

