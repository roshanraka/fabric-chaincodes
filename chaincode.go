/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	s "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var i int

//Entity - Structure for an entity like user
type Entity struct {
	Type    string  `json:"type"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Points  int     `json:"points"`
}

//Product - Structure for products used in buy coffee
type Product struct {
	Name   string  `json:"name"`
	Points int     `json:"points"`
	Amount float64 `json:"amount"`
	Entity string  `json:"entity"`
	Qty    int     `json:"qty"`
}

//TxnTopup - User transactions for adding points or balance
type TxnTopup struct {
	Initiator string `json:"initiator"`
	Remarks   string `json:"remarks"`
	ID        string `json:"id"`
	Time      string `json:"time"`
	Value     string `json:"value"`
	Asset     string `json:"asset"`
}

//TxnGoods - User transaction details for buying coffee
type TxnGoods struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Remarks  string `json:"remarks"`
	ID       string `json:"id"`
	Time     string `json:"time"`
	Value    string `json:"value"`
	Asset    string `json:"asset"`
}


//Chaincode  - struct consisting of all the chaincode funcs
type Chaincode struct {
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	key1 := args[0] //customer
	
	cust := Entity{
		Type:    "customer",
		Name:    key1,
		Balance: 3000,
		Points:  3000,
	}
	fmt.Println(cust)
	bytes, err := json.Marshal(cust)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key1, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	
	// Initialize the collection of  keys for products and various transactions
}

// Invoke isur entry point to invoke a chaincode function
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions/transactions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}  else if function == "add" {
		return t.add(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *Chaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("running write()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. expecting 4")
	}

	//writing a new customer to blockchain
	typeOf := args[0]
	name := args[1]
	balance, err := strconv.ParseFloat(args[2], 64)
	points, err := strconv.Atoi(args[3])
	entity := Entity{
		Type:    typeOf,
		Name:    name,
		Balance: balance,
		Points:  points,
	}
	fmt.Println(entity)
	bytes, err := json.Marshal(entity)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(name, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	return nil, nil
}

// read - query function to read key/value pair
func (t *Chaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	customer := Entity{}
	err = json.Unmarshal(bytes, &customer)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}
	bytes, err = json.Marshal(customer)
	if err != nil {
		fmt.Println("Error marshaling customer")
		return nil, errors.New("Error marshaling customer")
	}

	fmt.Println(bytes)
	return bytes, nil
}


func (t *Chaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("add is running ")

	if len(args) != 3 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 3 for add")
	}

	asset := args[0] //points or balance
	key := args[1]   //Entity ex: customer
	//amt, err := strconv.Atoi(args[1]) // points to be issued

	// GET the state of entity from the ledger
	bytes, err := stub.GetState(key)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key)
	}

	entity := Entity{}
	err = json.Unmarshal(bytes, &entity)
	if err != nil {
		fmt.Println("Error Unmarshaling entity Bytes")
		return nil, errors.New("Error Unmarshaling entity Bytes")
	}

	// Perform the addition of assests
	if asset == "points" {
		amt, err := strconv.Atoi(args[2])
		if err == nil {
			entity.Points = entity.Points + amt
			fmt.Println("entity Points = ", entity.Points)
		}
	} else {
		amt, err := strconv.ParseFloat(args[2], 64)
		if err == nil {
			entity.Balance = entity.Balance + amt
			fmt.Println("entity Balance = ", entity.Balance)
		}
	}

	// Write the state back to the ledger
	bytes, err = json.Marshal(entity)
	if err != nil {
		fmt.Println("Error marshaling entity")
		return nil, errors.New("Error marshaling entity")
	}
	err = stub.PutState(key, bytes)
	if err != nil {
		return nil, err
	}

	ID := stub.GetTxID()
	blockTime, err := stub.GetTxTimestamp()
	args = append(args, ID)
	args = append(args, blockTime.String())
	

	return nil, nil
}













