# Simplyfi-softtech



Preliminary Simple Chaincode for handling token distribution on completion of tasks given to an employee by his manager and to get the Champion for a closing sprint.



peer chaincode install -n task -v 1.0 -p github.com/chaincode/task/

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n task -v 1.0 -c '{"Args":["Viswa", "0"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

peer chaincode query -C $CHANNEL_NAME -n task -c '{"Args":["getEntity","Viswa"]}'

peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C $CHANNEL_NAME -n task -c '{"Args":["createEntity","Roshan","chain code writing","0"]}'

peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C $CHANNEL_NAME -n task -c '{"Args":["taskCompletion","Roshan","chain code writing","1"]}'

peer chaincode query -C $CHANNEL_NAME -n task -c '{"Args":["getEntity","Roshan"]}'

peer chaincode query -C $CHANNEL_NAME -n task -c '{"Args":["getChampion",""]}'
