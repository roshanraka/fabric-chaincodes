package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Account - Structure for an entity account
type Account struct {
	AccountNum           string             `json:"num"` //key for rocksDB
	AssetBalances map[string]float64 `json:"assets"`
}

// User - Structure for an entity
type User struct {
	UserID   string             `json:"userId"` //key for rocksDB
	UserName string             `json:"userName"`
	Accounts []string `json:"accounts"`
}

var NumAccounts int = 0

// AssetMgmt - chaincode struct
type AssetMgmt struct {
}

func main() {
	err := shim.Start(new(AssetMgmt))
	if err != nil {
		fmt.Printf("Error starting AssetMgmt chaincode: %s", err)
	}
}

const transactionKey = "transaction"
const assetNamesKey = "assetNames"

var docMeta = make(map[string]string)

// Init resets all the things
func (t *AssetMgmt) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// initLedger
func (t *AssetMgmt) initLedger(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
		
	accounts := []Account{
		Account{AccountNum: "", AssetBalances: map[string]float64{"T2Parts": 1000}},
		Account{AccountNum: "", AssetBalances: map[string]float64{"T1Parts": 1000, "T2Parts": 100}},
		Account{AccountNum: "", AssetBalances: map[string]{"OEMParts": 1000, "T1Parts": 100}},
		Account{AccountNum: "", AssetBalances: map[string]float64{"OEMParts": 100}},
	}
	i := 0
	for i < len(accounts) {
		fmt.Println("i is ", i)
		NumAccounts++
		accountNum := "Account-"+strconv.Itoa(NumAccounts)
		accounts[i].AccountNum = accountNum
		t.saveInBlockchain(stub, accountNum, accounts[i])
		fmt.Println("Saved to ledger", accounts[i])
		i = i + 1
	}
	// Initialize the collection of  keys for products and various transactions
	fmt.Println("Initializing keys collection")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err = stub.PutState(assetNamesKey, blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize Products key collection")
	}
	
	t.appendKey(stub, assetNamesKey, "T2Parts")
	t.appendKey(stub, assetNamesKey, "T1Parts")
	t.appendKey(stub, assetNamesKey, "OEMParts")
	
	return shim.Success(nil)
}

// Invoke is entry point to invoke a chaincode function
func (t *AssetMgmt) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions/transactions
	if function == "initLedger" {
		return t.initLedger(stub, args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	} else if function == "createUser" {
		return t.createUser(stub, args)
	} else if function == "createAsset" {
		return t.createAsset(stub, args)
	} else if function == "issueMore" {
		return t.issueMore(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "exchange" {
		return t.exchange(stub, args)
	} else if function == "produce" {
		return t.produce(stub, args)
	}
	if function == "getEntity" {
		return t.read(stub, args)
	} else if function == "getAssetTypes" {
		return t.getAssetTypes(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation: " + function)
}

func (t *AssetMgmt) createUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	userID := args[0]
	user := User{
		UserID : userID,
		UserName : args[1],
		Accounts : []string{}
	}

	// Write the state to the ledger
	t.saveInBlockchain(stub, userID, user)
	
}
// write - invoke function to write new key/value pair ex: Account
func (t *AssetMgmt) createAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("running createAccount()")
	userID := args[0]

	// Get state of user from blockchain
	bytes, err := stub.GetState(userID)
	user := User{}
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error("User not found")
	}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		fmt.Println("Error Unmarshaling account ")
		return shim.Error("Error Unmarshaling account")
	}

	// crerateAccount
	assets := make(map[string]float64)
	NumAccounts++
	accountNum := "Account-"+strconv.Itoa(NumAccounts)
 	account := Account{
		AccountNum: accountNum,
		AssetBalances: assets,
	}

	// Write the acct state to the ledger
	bytes, err := t.saveInBlockchain(stub, accountNum, account)
	if err != nil {
		return shim.Error(err.Error())
	}

	user.Accounts =append(user.Accounts,accountNum)
	// Write the user state to the ledger
	bytes, err := t.saveInBlockchain(stub, userID, user)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// read - query function to read key/value pair
func (t *AssetMgmt) read(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("read() is running")

	key := args[0] // name of Account

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return shim.Error("Error retrieving " + key)
	}

	fmt.Println(bytes)
	return shim.Success(bytes)
}

// getAssetTypes - query function to return the asset-types for this chain
func (t *AssetMgmt) getAssetTypes(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("getAssetTypes() is running")

	bytes, err := stub.GetState(assetNamesKey)
	if err != nil {
		fmt.Println("Error retrieving ")
		return shim.Error("Error retrieving ")
	}

	fmt.Println(bytes)
	return shim.Success(bytes)
}

// add - invoke funcrion to add/issue assets to an account
func (t *AssetMgmt) issueMore(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("issueMore is running ")

	accountNum := args[0]                       //Account key
	assetName := args[1]                        //Asset_type
	amt, err := strconv.ParseFloat(args[2], 64) //qty

	if t.isAssetCreated(stub, assetName) == false {
		return shim.Error("Asset not created")
	}

	// GET the state of account from the ledger
	account, err := t.getAccount(stub, accountNum)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Perform the addition of assests
	// _, ok := account.AssetBalances[assetName]
	//      if ok == false{
	account.AssetBalances[assetName] = account.AssetBalances[assetName] + amt

	// Write the state back to the ledger
	t.saveInBlockchain(stub, accountNum, account)

	return shim.Success(nil)
}

//transfer - function to transfer an asset between any two entities
func (t *AssetMgmt) transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("transfer is running ")

	fromAccountNum := args[0] // fromAccount
	toAccountNum := args[1]   // toAccount
	assetName := args[2]      //key for Asset map
	remarks := args[4]

	// GET the state of fromAccount from the ledger
	fromAccount, err := t.getAccount(stub, fromAccountNum)
	if err != nil {
		return shim.Error(err.Error())
	}

	// GET the state of toAccount from the ledger
	toAccount, err := t.getAccount(stub, toAccountNum)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Perform transfer of Asset
	amt, err := strconv.ParseFloat(args[3], 64)
	if fromAccount.AssetBalances[assetName] > amt {
		toAccount.AssetBalances[assetName] = toAccount.AssetBalances[assetName] + amt
		fmt.Println("account asset balance = ", toAccount.AssetBalances[assetName])
		fromAccount.AssetBalances[assetName] = fromAccount.AssetBalances[assetName] - amt
	} else {
		return shim.Error("Insufficient assets - transfer")
	}

	// Write the state back to the ledger
	t.saveInBlockchain(stub, fromAccountNum, fromAccount)

	t.saveInBlockchain(stub, toAccountNum, toAccount)

	// Write the remarks(docMeta) to the ledger
	txID := stub.GetTxID()
	bytes, err := json.Marshal(remarks)
	if err != nil {
		fmt.Println("Error marshaling remarks")
		return shim.Error("Error marshaling remarks")
	}
	err = stub.PutState(txID, bytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	t.appendKey(stub, transactionKey, txID)

	return shim.Success(nil)
}

//exchange - function - exchange of any two assets between two entities
func (t *AssetMgmt) exchange(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("exchange is running ")

	assetName1 := args[1]
	amt1, err := strconv.ParseFloat(args[2], 64)
	assetName2 := args[4]
	amt2, err := strconv.ParseFloat(args[5], 64)

	// GET the state of Account1 from the ledger
	account1, err := t.getAccount(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	// GET the state of Account2 from the ledger
	account2, err := t.getAccount(stub, args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	// Perform exchange
	if account1.AssetBalances[assetName1] > amt1 && account2.AssetBalances[assetName2] > amt2 {
		account1.AssetBalances[assetName1] = account1.AssetBalances[assetName1] - amt1
		account2.AssetBalances[assetName1] = account2.AssetBalances[assetName1] + amt1
		account2.AssetBalances[assetName2] = account2.AssetBalances[assetName2] - amt2
		account1.AssetBalances[assetName2] = account1.AssetBalances[assetName2] + amt2
	} else {
		return shim.Error("Insufficient assets")
	}

	// Write the account1 state back to the ledger
	t.saveInBlockchain(stub, args[0], account1)

	// Write the account2 state back to the ledger
	t.saveInBlockchain(stub, args[3], account2)

	return shim.Success(nil)
}

func (t *AssetMgmt) isAssetCreated(stub shim.ChaincodeStubInterface, key string) bool {

	bytes, _ := stub.GetState(assetNamesKey)
	var keys []string
	_ = json.Unmarshal(bytes, &keys)
	for _, keyValue := range keys {
		if keyValue == key {
			return true
		}
	}
	return false
}

func (t *AssetMgmt) createAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("createAsset is running.")

	/*
	   args[] - {accountNum, assettypeName for new assetType, qty}
	*/

	accountNum := args[0]
	assetName := args[1]
	qty, err := strconv.ParseFloat(args[2], 64)

	// GET the state of Account from the ledger
	account, err := t.getAccount(stub, accountNum)
	if err != nil {
		return shim.Error(err.Error())
	}

	if t.isAssetCreated(stub, assetName) == true {
		return shim.Error("Asset already created")
	}

	account.AssetBalances[assetName] = qty

	t.saveInBlockchain(stub, accountNum, account)
	// Add assetTypeName to list of assets.
	return t.appendKey(stub, assetNamesKey, assetName)
}

//produce - function - use one asset to produce another asset
func (t *AssetMgmt) produce(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("produce is running ")

	/*
	   args[] - {account, assetName1, amt1, assetName2, amt2,}
	*/

	accountNum := args[0]
	assetName1 := args[1]
	amt1, err := strconv.ParseFloat(args[2], 64)
	assetName2 := args[3]
	amt2, err := strconv.ParseFloat(args[4], 64)

	// GET the state of Account from the ledger
	account, err := t.getAccount(stub, accountNum)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Perform production :D
	if account.AssetBalances[assetName1] > amt1 {
		account.AssetBalances[assetName1] = account.AssetBalances[assetName1] - amt1
		account.AssetBalances[assetName2] = account.AssetBalances[assetName2] + amt2
	} else {
		return shim.Error("Insufficient assets - produce")
	}

	// Write the account1 state back to the ledger
	t.saveInBlockchain(stub, accountNum, account)

	return shim.Success(nil)
}

func (t *AssetMgmt) appendKey(stub shim.ChaincodeStubInterface, primeKey string, key string) peer.Response {
	fmt.Println("appendKey is running " + primeKey + " " + key)

	bytes, err := stub.GetState(primeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	var keys []string
	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return shim.Error(err.Error())
	}
	keys = append(keys, key)

	t.saveInBlockchain(stub, primeKey, keys)

	return shim.Success(nil)
}

func (t *AssetMgmt) saveInBlockchain(stub shim.ChaincodeStubInterface, key string, value interface{}) ([]byte, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		fmt.Println("saveInBlockchain: Error marsalling")
		return nil, errors.New("saveInBlockchain: Error marsalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, errors.New("saveInBlockchain: Error writing")
	}
	return bytes, err
}

func (t *AssetMgmt) getAccount(stub shim.ChaincodeStubInterface, accountNum string) (Account, error) {
	bytes, err := stub.GetState(accountNum)
	account := Account{}
	if err != nil {
		return account, err
	}
	if bytes == nil {
		return account, errors.New("Account not found")
	}

	err = json.Unmarshal(bytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshaling account ")
		return account, errors.New("Error Unmarshaling account")
	}
	return account, nil
}
