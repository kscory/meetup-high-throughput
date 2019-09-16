package highthroughputphantom

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type marble struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Size       int    `json:"size"`
}

type marbleResponse struct {
	Marble interface{} `json:"marble"`
	Owner  string      `json:"owner"`
	Amount int         `json:"amount"`
}

const (
	KEY_TRANSFER = "Transfer/name/sender/receiver/amount/txid"

	FUNCTION_INIT = "initMarbles"
	FUNCTION_TRANSFER = "transferMarbles"
	FUNCTION_READ = "readMarbles"
	FUNCTION_PRUNE = "pruneMarbles"
)

type HighThroughputChaincode struct {
}

func main() {
	err := shim.Start(new(HighThroughputChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *HighThroughputChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *HighThroughputChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == FUNCTION_INIT {
		return t.initMarbles(stub, args)
	} else if function == FUNCTION_TRANSFER {
		return t.transferMarbles(stub, args)
	} else if function == FUNCTION_READ {
		return t.readMarbles(stub, args)
	} else if function == FUNCTION_PRUNE {
		return t.pruneMarbles(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

/**
 * initMarbles - create a new marble, store into chaincode state
 * to give in the args array are as follows:
 *	- args[0] -> name; name of marble (key)
 *	- args[1] -> color; color of marble
 *	- args[2] -> size; size of marble
 *	- args[3] -> amount; total amount of marble
 *	- args[4] -> owner; owner id for this marble
 *
 * @param stub The chaincode shim
 * @param args The arguments array for the initMarbles invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *HighThroughputChaincode) initMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("[invoke] Call initMarbles")

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// Input sanitation
	fmt.Println("- start init marble")
	marbleName := args[0]
	color := strings.ToLower(args[1])
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	amount := args[3]
	_, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4rd argument must be a numeric string")
	}
	owner := strings.ToLower(args[4])

	// Check if marble already exists
	marbleAsBytes, err := stub.GetState(marbleName)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if marbleAsBytes != nil {
		fmt.Println("This marble already exists: " + marbleName)
		return shim.Error("This marble already exists: " + marbleName)
	}

	// Create marble object and marshal to JSON
	objectType := "marble"
	marble := &marble{objectType, marbleName, color, size}
	marbleJSONasBytes, err := json.Marshal(marble)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Save marble to state
	err = stub.PutState(marbleName, marbleJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Save marble amount to owner
	txID := stub.GetTxID()
	compositeKey, err  := stub.CreateCompositeKey(KEY_TRANSFER, []string{marbleName, "", owner, amount, txID})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(compositeKey, []byte{0x00})

	return shim.Success(nil)
}

/**
 * transferMarbles - transfer a marble by setting a new owner name on the marble
 * to give in the args array are as follows:
 *	- args[0] -> name; name of marble (key)
 *	- args[1] -> sender;
 *	- args[2] -> receiver;
 *	- args[3] -> amount; amount to transfer
 *
 * @param stub The chaincode shim
 * @param args The arguments array for the transferMarbles invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *HighThroughputChaincode) transferMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("[invoke] Call transferMarbles")

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	marbleName := args[0]
	sender := strings.ToLower(args[1])
	receiver := strings.ToLower(args[2])
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4rd argument must be a numeric string")
	} else if amount < 0 {
		return shim.Error("amount cannot be negative")
	}

	// check marble is existed
	marbleAsBytes, err := stub.GetState(marbleName)
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if marbleAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	// check sender amount
	//senderAmount, err := getAmount(stub, marbleName, sender)
	//if err != nil {
	//	return shim.Error("Cannot get sender Amount, err: " + err.Error())
	//}

	// check sender can transfer amount
	//if senderAmount < amount {
	//	fmt.Println(senderAmount)
	//	fmt.Println(amount)
	//	return shim.Error("Cannot transfer amount:")
	//}

	// Save amount
	txID := stub.GetTxID()
	compositeKey, err  := stub.CreateCompositeKey(KEY_TRANSFER, []string{marbleName, sender, receiver, strconv.Itoa(amount), txID})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(compositeKey, []byte{0x00})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

/**
 * readMarbles - read a marble from chaincode state
 * to give in the args array are as follows:
 *	- args[0] -> name; name of marble (required)
 *	- args[1] -> owner; owner of marble (not required)
 *
 * @param stub The chaincode shim
 * @param args The arguments array for the readMarbles query
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *HighThroughputChaincode) readMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("[query] Call readMarbles")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name := args[0]
	marbleAsbytes, err := stub.GetState(name) //get the marble from chaincode state
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + name + "\"}")
	} else if marbleAsbytes == nil {
		return shim.Error("{\"Error\":\"Marble does not exist: " + name + "\"}")
	}
	marble := marble{}
	err = json.Unmarshal(marbleAsbytes, &marble) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	result := &marbleResponse{marble, "", 0}

	// if parameter includes owner, return with amount info
	if len(args) > 1 {
		owner := args[1]
		result.Owner = owner
		ownerAmount, err := getAmount(stub, name, owner)
		if err != nil {
			return shim.Error("Cannot get sender Amount, err: " + err.Error())
		}
		result.Amount = ownerAmount
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultBytes)
}

/**
 * pruneMarbles - pruning for marbles
 * to give in the args array are as follows:
 *	- args[0] -> name; name of marble (key)
 *
 * @param stub The chaincode shim
 * @param args The arguments array for the pruneMarbles invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *HighThroughputChaincode) pruneMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("[invoke] Call pruneMarbles")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to prune")
	}
	name := args[0]

	// check marble is existed
	marbleAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if marbleAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	finalValue := make(map[string]int)
	amountIterator, err := stub.GetStateByPartialCompositeKey(KEY_TRANSFER, []string{name})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer amountIterator.Close()
	if amountIterator.HasNext() {
		for amountIterator.HasNext() {
			responseRange, err := amountIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			// Split Composite Key
			_, keyParts, err := stub.SplitCompositeKey(responseRange.Key)
			if err != nil {
				return shim.Error(err.Error())
			}
			sender, receiver := keyParts[1], keyParts[2]
			amount, err := strconv.Atoi(keyParts[3])
			if err != nil {
				return shim.Error(err.Error())
			}

			if len(sender) != 0 {
				senderAmount, exists := finalValue[sender]
				if !exists {
					senderAmount = 0 - amount
				} else {
					senderAmount = senderAmount - amount
				}
				finalValue[sender] = senderAmount
			}
			if len(receiver) != 0 {
				receiverAmount, exists := finalValue[receiver]
				if !exists {
					receiverAmount = 0 + amount
				} else {
					receiverAmount = receiverAmount + amount
				}
				finalValue[receiver] = receiverAmount
			}

			// Del State
			err = stub.DelState(responseRange.Key)
			if err != nil {
				return shim.Error(err.Error())
			}
		}
	}

	// Save Amount
	txID := stub.GetTxID()
	for owner, value := range finalValue {
		compositeKey, err := stub.CreateCompositeKey(KEY_TRANSFER, []string{name, "", owner, strconv.Itoa(value), txID})
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(compositeKey, []byte{0x00})
		if err != nil {
			return shim.Error(err.Error())
		}
	}


	return shim.Success(nil)
}

func getAmount(stub shim.ChaincodeStubInterface, marbleName, owner string) (int, error) {
	amountResult := 0

	amountIterator, err := stub.GetStateByPartialCompositeKey(KEY_TRANSFER, []string{marbleName})
	if err != nil {
		return 0, err
	}
	defer amountIterator.Close()

	if amountIterator.HasNext() {
		for amountIterator.HasNext() {
			responseRange, err := amountIterator.Next()
			if err != nil {
				return 0, err
			}
			// Split Composite Key
			_, keyParts, err := stub.SplitCompositeKey(responseRange.Key)
			if err != nil {
				return 0, err
			}
			sender, receiver := keyParts[1], keyParts[2]
			if sender == owner {
				amount := keyParts[3]
				amountInt, err := strconv.Atoi(amount)
				if err != nil {
					return 0, err
				}
				amountResult -= amountInt
			} else if receiver == owner {
				amount := keyParts[3]
				amountInt, err := strconv.Atoi(amount)
				if err != nil {
					return 0, err
				}
				amountResult += amountInt
			}
		}
	}
	return amountResult, nil
}
