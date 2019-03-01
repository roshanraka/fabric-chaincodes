package main

import (
		"encoding/json"
		"fmt"
		"strconv"
		"github.com/hyperledger/fabric/core/chaincode/shim"
		"github.com/hyperledger/fabric/protos/peer"
		)

//Entity - Structure for an entity account
type Entity struct {
	Name   	string          `json:"name"` //key
	Tasks 	map[string]bool `json:"tasks"`
	Tokens  int             `json:"tokens"`
}

var champion string
var max int

// TasksMgmt - chaincode struct
type TasksMgmt struct {
}

func main() {
	err := shim.Start(new(TasksMgmt))
	if err != nil {
		fmt.Printf("Error starting TasksMgmt chaincode: %s", err)
	}
}


func (t * TasksMgmt) Init(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetStringArgs()

	name := args[0]
	tasks := map[string]float64{}
	tokens := strconv.Atoi(args[2])

	entity := Entity{
		Name:   name,
		Tasks:   tasks,
		Tokens: tokens,
	}

	_, err := t.saveInBlockchain(stub, name, entity)
	if err != nil {
		return shim.Error(err.Error())
	}

	max = 0
	return shim.Success(nil)
}

func (t *TasksMgmt) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions/transactions
	if function == "taskCompletion" {
		return t.taskCompletion(stub, args)
	} else if function == "createEntity" {
		return t.createEntity(stub, args)
	} else if function == "addTask" {
		return t.addTask(stub, args)
	}

	if function == "getEntity" {
		return t.read(stub, args)
	} else if function == "getChampion" {
		return t.getChampion(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)

	return shim.Error("Received unknown function invocation: " + function)
}

func (t * TasksMgmt) taskCompletion(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
		name := args[0]
		task := args[1]
		tokens := strconv.Atoi(args[2])

		entity, err := t.getEntity(stub, name)
		if err != nil {
			return shim.Error(err.Error())
		}
		entity.Tasks[task] = true
		entity.Tokens += tokens

		_, err := t.saveInBlockchain(stub, name, entity)
		if err != nil {
			return shim.Error(err.Error())
		}

		if entity.Tokens > max {
			max = entity.Tokens
			champion = name
		}
		return shim.Success(nil)
}

func (t * TasksMgmt) addTask(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
		name := args[0]
		task := args[1]

		entity, err := t.getEntity(stub, name)
		if err != nil {
			return shim.Error(err.Error())
		}
		entity.Tasks[task] = false

		_, err := t.saveInBlockchain(stub, name, entity)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
}

func (t * TasksMgmt) getChampion(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	bytes, err := stub.GetState(champion)
	if err != nil {
		return shim.Error("Failed to get state of " + champion)
	}

	return shim.Success(bytes)
}

func (t * TasksMgmt) createEntity(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	name := args[0]
	task := args[1]
	tasks := map[string]float64{"task" : false}

	tokens := args[2])

	entity := Entity{
		Name:   name,
		Tasks:   tasks,
		Tokens: tokens,
	}

	_, err := t.saveInBlockchain(stub, name, entity)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t * TasksMgmt) saveInBlockchain(stub shim.ChaincodeStubInterface, key string, value interface{}) ([]byte, error)) {
	fmt.Println("saveInBlockchain")
	
	bytes, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Error marsalling")
		return shim.Error("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return shim.Error(err.Error())
	}
	return bytes, nil
}


func (t * TasksMgmt) getEntity(stub shim.ChaincodeStubInterface, key string) (Entity, error) {

	// GET the state of entity from the ledger
	bytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get state of " + key)
	}

	entity := Entity{}
	err = json.Unmarshal(bytes, &entity)
	if err != nil {
		fmt.Println("Error Unmarshaling entity Bytes")
		return shim.Error("Error Unmarshaling entity Bytes")
	}

	return entity, nil
}