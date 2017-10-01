/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
        "strconv"
        "encoding/json"
        "bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Stock struct {
        StockName string `json:"stock"`
        Quantity float64 `json:"quantity"`
}

type User struct {
	Name   string `json:"name"`
        Balance  int `json:"bal"`
        StockUnits []Stock `json:"stock"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var Aval, Bval int // Asset holdings
        var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
        }
        fmt.Println("Values ", args[0])

	// Initialize the chaincode
        Aval, err = strconv.Atoi(args[1])

	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
        }
        var user1 = User{Name: args[0],  Balance: Aval}
        fmt.Println("User 1", user1);

        Bval, err = strconv.Atoi(args[3])

	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
        }

        var user2 = User{Name: args[2],  Balance: Bval}
	fmt.Println("User 2", user2);

        user1AsBytes, _ := json.Marshal(user1)
	// Write the state to the ledger
	err = stub.PutState(user1.Name, user1AsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        user2AsBytes, _ := json.Marshal(user2)
        err = stub.PutState(user2.Name, user2AsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


        return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "buyShares" {
		// Make payment of X units from A to B
		return t.buyShares(stub, args)
	} else if function == "sellShares" {
		// Deletes an entity from its state
		return t.sellShares(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "addUser" {
		// the old "Query" is now implemtned in invoke
		return t.addUser(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"buyShare\" \"sellShare\" \"query\" \"addUser\"")
}

func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Aval int // Asset holdings
        var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        fmt.Println("Values ", args[0])

	// Initialize the chaincode
        Aval, err = strconv.Atoi(args[1])

	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
    }

	var user1 = User{Name: args[0],  Balance: Aval}

	fmt.Println("User 1", user1);

    user1AsBytes, _ := json.Marshal(user1)
	// Write the state to the ledger
	err = stub.PutState(user1.Name, user1AsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
    return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) buyShares(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Buyer, Payee, StockName string    // Entities
	var price int // Asset holdings
	var quantity float64          // Transaction value
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	Buyer = args[0]
        Payee = args[1]
        StockName = args[2]
        price, _ = strconv.Atoi(args[4])
        quantity, _ = strconv.ParseFloat(args[3], 64)

        // Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	BuyerValbytes, err := stub.GetState(Buyer)
	if err != nil {
		return shim.Error("Failed to get state for buyer")
	}
	if BuyerValbytes == nil {
		return shim.Error("Entity not found for buyer")
        }
	buyerState := User{}

        json.Unmarshal(BuyerValbytes, &buyerState)
        fmt.Println("buyerState before", buyerState)



	PayeeValBytes, err := stub.GetState(Payee)
	if err != nil {
		return shim.Error("Failed to get state for payee")
	}
	if PayeeValBytes == nil {
		return shim.Error("Entity not found for payee")
        }
        payeeState := User{}

	json.Unmarshal(PayeeValBytes, &payeeState)
	 fmt.Println("Payee state befoew", payeeState)

	buyerState.Balance= buyerState.Balance - price
        payeeState.Balance= payeeState.Balance + price

        var isStockFound = false
        var stockUnits = Stock{}
        i := 0
		for i < len(buyerState.StockUnits) {

		if buyerState.StockUnits[i].StockName == StockName {
                        stockUnits = buyerState.StockUnits[i]
                        isStockFound = true
                        break
                }
                i= i+1
        }
        if isStockFound {
                stockUnits.Quantity+=quantity
                buyerState.StockUnits[i] = stockUnits
        } else {
                stockUnits = Stock{StockName: StockName, Quantity: quantity }
                buyerState.StockUnits = append(buyerState.StockUnits, stockUnits);
        }

        fmt.Println("buyerState", buyerState)
        fmt.Println("Payee state", payeeState)


        // Write the state back to the ledger
        buyerAsBytes, _ := json.Marshal(buyerState)
	err = stub.PutState(buyerState.Name, buyerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        payeeAsBytes, _ := json.Marshal(payeeState)
	err = stub.PutState(payeeState.Name, payeeAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) sellShares(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Seller, Payer, StockName string    // Entities
	var price int // Asset holdings
	var quantity float64          // Transaction value
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	Seller = args[0]
        Payer = args[1]
        StockName = args[2]
        price, _ = strconv.Atoi(args[4])
        quantity, _ = strconv.ParseFloat(args[3], 64)

        // Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	SellerValbytes, err := stub.GetState(Seller)
	if err != nil {
		return shim.Error("Failed to get state for buyer")
	}
	if SellerValbytes == nil {
		return shim.Error("Entity not found for buyer")
        }
	sellerState := User{}

	json.Unmarshal(SellerValbytes, &sellerState)


	PayerValBytes, err := stub.GetState(Payer)
	if err != nil {
		return shim.Error("Failed to get state for payee")
	}
	if PayerValBytes == nil {
		return shim.Error("Entity not found for payee")
        }
        payerState := User{}

	json.Unmarshal(PayerValBytes, &payerState)



	sellerState.Balance= sellerState.Balance + price
        payerState.Balance= payerState.Balance - price

        i := 0
        var stockUnits = Stock{}
	for i < len(sellerState.StockUnits) {

		if sellerState.StockUnits[i].StockName == StockName {
                        stockUnits = sellerState.StockUnits[i]
                        break
                }
                i= i+1
        }

        stockUnits.Quantity-=quantity
        sellerState.StockUnits[i] = stockUnits

        fmt.Println("seller State", sellerState)
        fmt.Println("Payee state", payerState)


        // Write the state back to the ledger
        sellerAsBytes, _ := json.Marshal(sellerState)
	err = stub.PutState(sellerState.Name, sellerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        payerAsBytes, _ := json.Marshal(payerState)
	err = stub.PutState(payerState.Name, payerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var UserName string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	UserName = args[0]

	// Get the state from the ledger
	uservalbytes, err := stub.GetState(UserName)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + UserName + "\"}"
		return shim.Error(jsonResp)
	}

        var buffer bytes.Buffer
        buffer.WriteString("{\"User\":" +  string(uservalbytes) + "}")
	fmt.Printf("Query Response:%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
