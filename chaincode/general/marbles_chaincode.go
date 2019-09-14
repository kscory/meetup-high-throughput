package general

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
	marble interface{} `json:"marble"`
	owner  string      `json:"owner"`
	amount int         `json:"amount"`
}

type SimpleChaincode struct {
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initMarbles" {
		return t.initMarbles(stub, args)
	} else if function == "transferMarbles" {
		return t.transferMarbles(stub, args)
	} else if function == "readMarbles" {
		return t.readMarbles(stub, args)
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
func (t *SimpleChaincode) initMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	key := owner + marbleName
	err = stub.PutState(key, []byte(amount))
	if err != nil {
		return shim.Error(err.Error())
	}

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
func (t *SimpleChaincode) transferMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	}

	// check marble is existed
	marbleAsBytes, err := stub.GetState(marbleName)
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if marbleAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	// check sender amount is existed
	senderAmountAsBytes, err := stub.GetState(sender + marbleName)
	if err != nil {
		return shim.Error("Failed to get sender amount of marbles:" + err.Error())
	} else if senderAmountAsBytes == nil {
		return shim.Error("Sender does not have marbles")
	}
	sendAmount, err := strconv.Atoi(string(senderAmountAsBytes))
	if err != nil {
		return shim.Error("Failed to get sender amount of marbles:" + err.Error())
	}

	// check sender can transfer amount
	if sendAmount < amount {
		return shim.Error("Cannot transfer amount:")
	}

	// receiver amount
	receiverAmountAsBytes, err := stub.GetState(receiver + marbleName)
	receiverAmount := 0
	if err != nil {
		return shim.Error("Failed to get amount of marbles:" + err.Error())
	} else if receiverAmountAsBytes != nil {
		receiverAmount, err = strconv.Atoi(string(senderAmountAsBytes))
		if err != nil {
			return shim.Error("Failed to get receiver amount of marbles:" + err.Error())
		}
	}

	// Save amount
	err = stub.PutState(sender+marbleName, []byte(strconv.Itoa(sendAmount-amount)))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(strconv.Itoa(receiverAmount + amount))
	err = stub.PutState(receiver+marbleName, []byte(strconv.Itoa(receiverAmount+amount)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// readMarbles - read a marble from chaincode state
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
func (t *SimpleChaincode) readMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		result.owner = owner
		ownerAmountAsBytes, err := stub.GetState(owner + name)
		if err != nil {
			return shim.Error("{\"Error\":\"Failed to get amount state for name:" + name + ", owner: " + owner + "\"}")
		}
		if ownerAmountAsBytes != nil {
			ownerAmount, err := strconv.Atoi(string(ownerAmountAsBytes))
			if err != nil {
				return shim.Error("Failed to get owner amount of marbles:" + err.Error())
			}
			result.amount = ownerAmount
		}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultBytes)
}
