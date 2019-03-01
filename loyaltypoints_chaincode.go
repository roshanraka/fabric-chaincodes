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

//Entity - Structure for an entity like user, merchant, bank
type Entity struct {
	Type    string  `json:"type"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Points  int     `json:"points"`
}

//Product - Structure for products used in buy goods
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

//TxnGoods - User transaction details for buying goods
type TxnGoods struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Remarks  string `json:"remarks"`
	ID       string `json:"id"`
	Time     string `json:"time"`
	Value    string `json:"value"`
	Asset    string `json:"asset"`
}

//TxnEncash - details of requests from merchant to encash points
type TxnEncash struct {
	Key       string `json:"key"`
	ID        string `json:"id"`
	Initiator string `json:"initiator"`
	Bank      string `json:"bank"`
	Points    int    `json:"points"`
	Amount    int    `json:"amount"`
	Remarks   string `json:"remarks"`
	Time      string `json:"time"`
}

//LoyaltyChaincode  - struct consisting of all the chaincode funcs
type LoyaltyChaincode struct {
}

func main() {
	err := shim.Start(new(LoyaltyChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *LoyaltyChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	key1 := args[0] //customer
	key2 := args[1] //merchant
	key3 := args[2] //bank

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

	merch := Entity{
		Type:    "merchant",
		Name:    key2,
		Balance: 6000,
		Points:  6000,
	}
	fmt.Println(merch)
	bytes, err = json.Marshal(merch)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key2, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	bank := Entity{
		Type:    "bank",
		Name:    key3,
		Balance: 10000,
		Points:  10000,
	}
	fmt.Println(bank)
	bytes, err = json.Marshal(bank)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key3, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	// Initialize the collection of  keys for products and various transactions
	fmt.Println("Initializing keys collection")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err = stub.PutState("Products", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize Products key collection")
	}
	err = stub.PutState("TxnTopup", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize TxnTopUp key collection")
	}
	err = stub.PutState("TxnGoods", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize TxnGoods key collection")
	}
	err = stub.PutState("TxnEncash", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize TxnEncash key collection")
	}

	fmt.Println("Initialization complete")

	t.addProduct(stub, []string{"Speakers", "2000", "200", key2, "100"})
	t.addProduct(stub, []string{"Headphones", "3000", "300",  key2, "100"})
	t.addProduct(stub, []string{"Mobile", "4000", "400", key2, "100"})
	t.addProduct(stub, []string{"BagPack", "1000", "100", key2, "100"})

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *LoyaltyChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions/transactions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "buyGoods" {
		return t.buyGoods(stub, args)
	} else if function == "add" {
		return t.add(stub, args)
	} else if function == "encashMerchant" {
		return t.encashMerchant(stub, args)
	} else if function == "approve" {
		return t.approve(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *LoyaltyChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	} else if function == "getAllProducts" {
		return t.getAllProducts(stub)
	} else if function == "getAllTxnTopup" {
		return t.getAllTxnTopup(stub)
	} else if function == "getAllTxnGoods" {
		return t.getAllTxnGoods(stub)
	} else if function == "getAllTxnEncash" {
		return t.getAllTxnEncash(stub)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *LoyaltyChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

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
func (t *LoyaltyChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *LoyaltyChaincode) buyGoods(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("buyGoods is running ")

	if len(args) != 6 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 6 for buy goods")
	}
	asset := args[0] //points or balance
	key1 := args[1]  //Entity1 ex: customer
	key2 := args[2]  //Entity2 ex: merchant
	key3 := args[3]  //Product Entity
	qty, err := strconv.Atoi(args[4])

	bytes, err := stub.GetState(key1)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key1)
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	customer := Entity{}
	err = json.Unmarshal(bytes, &customer)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}

	bytes, err = stub.GetState(key2)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key2)
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	merchant := Entity{}
	err = json.Unmarshal(bytes, &merchant)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}
	bytes, err = stub.GetState(key3)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key2)
	}
	product := Product{}
	err = json.Unmarshal(bytes, &product)
	if err != nil {
		fmt.Println("Error Unmarshaling product bytes")
		return nil, errors.New("Error Unmarshaling product Bytes")
	}
	if product.Entity == merchant.Name && product.Qty >= qty {
		// Perform the transfer
		if s.Compare(asset, "points") == 0 {
			fmt.Println("points transfer")
			//X, err := strconv.Atoi(args[3])
			if customer.Points >= product.Points*qty {
				customer.Points = customer.Points - product.Points*qty
				merchant.Points = merchant.Points + product.Points*qty
				product.Qty -= qty
				args[4] = strconv.Itoa(product.Points * qty)
				fmt.Printf("customer Points = %d, merchant Points = %d\n", customer.Points, merchant.Points)
			} else {
				return nil, errors.New("Insufficient points to buy goods")
			}
		} else {
			fmt.Println("balance to be added")
			//X, err := strconv.ParseFloat(args[3], 64)
			if customer.Balance >= product.Amount*float64(qty) {
				customer.Balance = customer.Balance - product.Amount*float64(qty)
				merchant.Balance = merchant.Balance + product.Amount*float64(qty)
				product.Qty -= qty
				args[4] = strconv.FormatFloat(product.Amount*float64(qty), 'f', -1, 64)
				fmt.Printf("customer Balance = %f, merchant Balance = %f\n", customer.Balance, merchant.Balance)
			} else {
				return nil, errors.New("Insufficient balance to buy goods")
			}
		}
		//product.Entity = customer.Name
		// Write the customer/entity1 state back to the ledger
		bytes, err = json.Marshal(customer)
		if err != nil {
			fmt.Println("Error marshaling customer")
			return nil, errors.New("Error marshaling customer")
		}
		err = stub.PutState(key1, bytes)
		if err != nil {
			return nil, err
		}

		// Write the merchant/entity2 state back to the ledger]
		bytes, err = json.Marshal(merchant)
		if err != nil {
			fmt.Println("Error marshaling customer")
			return nil, errors.New("Error marshaling customer")
		}
		err = stub.PutState(key2, bytes)
		if err != nil {
			return nil, err
		}
		// Write the product state back to the ledger
		bytes, err = json.Marshal(product)
		if err != nil {
			fmt.Println("Error marshaling customer")
			return nil, errors.New("Error marshaling customer")
		}
		err = stub.PutState(key3, bytes)
		if err != nil {
			return nil, err
		}

		args = append(args, stub.GetTxID())
		blockTime, err := stub.GetTxTimestamp()
		if err != nil {
			return nil, err
		}
		args = append(args, blockTime.String())
		t.putTxnGoods(stub, args)
	}

	return nil, nil
}

func (t *LoyaltyChaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

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
			fmt.Println("entity Points = ", entity.Points)
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
	t.putTxnTopup(stub, args)

	return nil, nil
}

func (t *LoyaltyChaincode) encashMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("encashMerchant is running ")

	if len(args) != 3 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 3 for encashMerchant")
	}

	points, err := strconv.Atoi(args[2])
	blockTime, err := stub.GetTxTimestamp()
	//time.Unix(blockTime.Seconds, 0)

	i++
	key := "encash" + strconv.Itoa(i)
	txn := TxnEncash{
		Key:       key,
		ID:        stub.GetTxID(),
		Initiator: args[0],
		Bank:      args[1],
		Points:    points,
		Amount:    points / 100,
		Remarks:   "New Request for Encashment",
		Time:      blockTime.String(),
	}

	bytes, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling encashMerchant")
		return nil, errors.New("Error marshaling encashMerchant")
	}

	err = stub.PutState(key, bytes)
	if err != nil {
		return nil, err
	}

	return t.appendKey(stub, "TxnEncash", key)
}

func (t *LoyaltyChaincode) approve(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("approve is running ")

	if len(args) != 1 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 1 for approve")
	}

	bytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state of encash key")
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	txn := TxnEncash{}
	err = json.Unmarshal(bytes, &txn)
	if err != nil {
		fmt.Println("Error Unmarshaling TxnEncash")
		return nil, errors.New("Error Unmarshaling TxnEncash")
	}

	bytes, err = stub.GetState(txn.Initiator)
	if err != nil {
		return nil, errors.New("Failed to get state of " + txn.Initiator)
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	merchant := Entity{}
	err = json.Unmarshal(bytes, &merchant)
	if err != nil {
		fmt.Println("Error Unmarshaling merchant encash")
		return nil, errors.New("Error Unmarshaling encash merchant")
	}

	bytes, err = stub.GetState(txn.Bank)
	if err != nil {
		return nil, errors.New("Failed to get state of " + txn.Bank)
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	bank := Entity{}
	err = json.Unmarshal(bytes, &bank)
	if err != nil {
		fmt.Println("Error Unmarshaling bank encash")
		return nil, errors.New("Error Unmarshaling encash bank")
	}

	// Perform encashment
	bank.Points = bank.Points + txn.Points
	merchant.Points = merchant.Points - txn.Points
	bank.Balance = bank.Balance - float64(txn.Amount)
	merchant.Balance = merchant.Balance + float64(txn.Amount)

	// Write the merchant/entity1 state back to the ledger
	bytes, err = json.Marshal(merchant)
	if err != nil {
		fmt.Println("Error marshaling merchant")
		return nil, errors.New("Error marshaling merchant")
	}
	err = stub.PutState(txn.Initiator, bytes)
	if err != nil {
		return nil, err
	}

	// Write the bank/entity2 state back to the ledger]
	bytes, err = json.Marshal(bank)
	if err != nil {
		fmt.Println("Error marshaling bank")
		return nil, errors.New("Error marshaling bank")
	}
	err = stub.PutState(txn.Bank, bytes)
	if err != nil {
		return nil, err
	}

	blockTime, err := stub.GetTxTimestamp()
	txn.Time = blockTime.String()
	txn.ID = stub.GetTxID()
	txn.Remarks = "Encashment Completed"
	// Write the TxnEncash state back to the ledger
	bytes, err = json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling TxnGoods")
		return nil, errors.New("Error marshaling TxnGoods")
	}
	err = stub.PutState(txn.Key, bytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *LoyaltyChaincode) addProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding product information")
	
	if len(args) != 5 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 5 for addProduct")
	}
	amt, err := strconv.ParseFloat(args[2], 64)
	points, err := strconv.Atoi(args[1])
	qty, err := strconv.Atoi(args[4])

	product := Product{
		Name:   args[0],
		Points: points,
		Amount: amt,
		Entity: args[3],
		Qty:    qty,
	}

	bytes, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling product")
		return nil, errors.New("Error marshaling product")
	}

	err = stub.PutState(product.Name, bytes)
	if err != nil {
		return nil, err
	}

	bytes, err = stub.GetState("Products")
	if err != nil {

	}
	var keys []string
	err = json.Unmarshal(bytes, &keys)
	if err != nil {

	}
	keys = append(keys, args[0])
	bytes, err = json.Marshal(keys)
	if err != nil {
		fmt.Println("Error marshaling product keys")
		return nil, errors.New("Error marshaling product keys")
	}
	err = stub.PutState("Products", bytes)
	if err != nil {

	}

	return nil, nil
}

func (t *LoyaltyChaincode) getAllProducts(stub shim.ChaincodeStubInterface) ([]byte, error) {

	fmt.Println("getAllProducts is running ")

	var products []Product

	// Get list of all the keys - Products
	keysBytes, err := stub.GetState("Products")
	if err != nil {
		fmt.Println("Error retrieving Products")
		return nil, errors.New("Error retrieving Products")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling Products")
		return nil, errors.New("Error unmarshalling Products")
	}

	// Get each product from "Products" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var product Product
		err = json.Unmarshal(bytes, &product)
		if err != nil {
			fmt.Println("Error retrieving product " + value)
			return nil, errors.New("Error retrieving product " + value)
		}

		fmt.Println("Appending product " + value)
		products = append(products, product)
	}

	bytes, err := json.Marshal(products)
	if err != nil {
		fmt.Println("Error marshaling product")
		return nil, errors.New("Error marshaling product")
	}
	return bytes, nil
}

func (t *LoyaltyChaincode) putTxnTopup(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("putTxnTopup is running ")

	if len(args) != 5 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 5 for putTxnTopup")
	}
	txn := TxnTopup{
		Initiator: args[1],
		Remarks:   args[0] + " addedd",
		ID:        args[3],
		Time:      args[4],
		Value:     args[2],
		Asset:     args[0],
	}

	bytes, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling TxnTopup")
		return nil, errors.New("Error marshaling TxnTopup")
	}

	err = stub.PutState(txn.ID, bytes)
	if err != nil {
		return nil, err
	}

	return t.appendKey(stub, "TxnTopup", txn.ID)
}

func (t *LoyaltyChaincode) getAllTxnTopup(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("getAllTxnTopup is running ")

	var txns []TxnTopup

	// Get list of all the keys - TxnTopup
	keysBytes, err := stub.GetState("TxnTopup")
	if err != nil {
		fmt.Println("Error retrieving TxnTopup keys")
		return nil, errors.New("Error retrieving TxnTopup keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling TxnTopup key")
		return nil, errors.New("Error unmarshalling TxnTopup keys")
	}

	// Get each product txn "TxnTopup" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var txn TxnTopup
		err = json.Unmarshal(bytes, &txn)
		if err != nil {
			fmt.Println("Error retrieving txn " + value)
			return nil, errors.New("Error retrieving txn " + value)
		}

		fmt.Println("Appending txn" + value)
		txns = append(txns, txn)
	}

	bytes, err := json.Marshal(txns)
	if err != nil {
		fmt.Println("Error marshaling txns topup")
		return nil, errors.New("Error marshaling txns topup")
	}
	return bytes, nil
}

func (t *LoyaltyChaincode) putTxnGoods(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("putTxnGoods is running ")

	if len(args) != 8 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 8 for putTxnGoods")
	}
	txn := TxnGoods{
		Sender:   args[1],
		Receiver: args[2],
		Remarks:  args[3] + " - " + args[5],
		ID:       args[6],
		Time:     args[7],
		Value:    args[4],
		Asset:    args[0],
	}

	bytes, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling TxnGoods")
		return nil, errors.New("Error marshaling TxnGoods")
	}

	err = stub.PutState(txn.ID, bytes)
	if err != nil {
		return nil, err
	}

	return t.appendKey(stub, "TxnGoods", txn.ID)
}

func (t *LoyaltyChaincode) getAllTxnGoods(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("getAllTxnGoods is running ")

	var txns []TxnGoods

	// Get list of all the keys - TxnGoods
	keysBytes, err := stub.GetState("TxnGoods")
	if err != nil {
		fmt.Println("Error retrieving TxnGoods keys")
		return nil, errors.New("Error retrieving TxnGoods keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling TxnGoods key")
		return nil, errors.New("Error unmarshalling TxnGoods keys")
	}

	// Get each txn from "TxnGoods" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var txn TxnGoods
		err = json.Unmarshal(bytes, &txn)
		if err != nil {
			fmt.Println("Error retrieving txn " + value)
			return nil, errors.New("Error retrieving txn " + value)
		}

		fmt.Println("Appending txn goods details " + value)
		txns = append(txns, txn)
	}

	bytes, err := json.Marshal(txns)
	if err != nil {
		fmt.Println("Error marshaling txns TxnGoods")
		return nil, errors.New("Error marshaling txns TxnGoods")
	}
	return bytes, nil
}

func (t *LoyaltyChaincode) getAllTxnEncash(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("getAllTxnEncash is running ")

	var txns []TxnEncash

	// Get list of all the keys - TxnGoods
	keysBytes, err := stub.GetState("TxnEncash")
	if err != nil {
		fmt.Println("Error retrieving TxnEncash keys")
		return nil, errors.New("Error retrieving TxnEncash keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling TxnEncash key")
		return nil, errors.New("Error unmarshalling TxnEncash keys")
	}

	// Get each txn from "TxnGoods" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var txn TxnEncash
		err = json.Unmarshal(bytes, &txn)
		if err != nil {
			fmt.Println("Error retrieving txn " + value)
			return nil, errors.New("Error retrieving txn " + value)
		}

		fmt.Println("Appending txn encash details " + value)
		txns = append(txns, txn)
	}

	bytes, err := json.Marshal(txns)
	if err != nil {
		fmt.Println("Error marshaling txns TxnEncash")
		return nil, errors.New("Error marshaling txns TxnEncash")
	}
	return bytes, nil
}
func (t *LoyaltyChaincode) appendKey(stub shim.ChaincodeStubInterface, primeKey string, key string) ([]byte, error) {
	fmt.Println("appendKey is running " + primeKey + " " + key)

	bytes, err := stub.GetState(primeKey)
	if err != nil {
		return nil, err
	}
	var keys []string
	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return nil, err
	}
	keys = append(keys, key)
	bytes, err = json.Marshal(keys)
	if err != nil {
		fmt.Println("Error marshaling " + primeKey)
		return nil, errors.New("Error marshaling keys" + primeKey)
	}
	err = stub.PutState(primeKey, bytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
