package general

import (
	"encoding/json"
	"marbles-meetup/util"
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var sampleMarble = &marble{"marble", "RedMarble", "red", 30}
const sender = "alice"
const receiver = "bob"
const totalAmount = 100000
const transferAmount = 30

func initMarble(t *testing.T) *shim.MockStub {
	var scc = new(SimpleChaincode)
	var stub = shim.NewMockStub("marbles", scc)
	stub.MockInit("1", [][]byte{[]byte("init")})
	arguments := [][]byte{[]byte(FUNCTION_INIT), []byte(sampleMarble.Name),
		[]byte(sampleMarble.Color), []byte(strconv.Itoa(sampleMarble.Size)),
		[]byte(strconv.Itoa(totalAmount)), []byte(sender)}
	util.CheckInvoke(t, stub, arguments)

	return stub
}

func Test_MARBLES_initMarble_success(t *testing.T) {
	fmt.Println("[TEST] initMarble")

	// invoke initMarbles
	stub := initMarble(t)

	// check marble state
	sampleMarbleBytes, _ := json.Marshal(sampleMarble)
	util.CheckState(t, stub, sampleMarble.Name, string(sampleMarbleBytes))

	// check owner amount
	util.CheckState(t, stub, sender + sampleMarble.Name, strconv.Itoa(totalAmount))
}

func Test_MARBLES_transferMarbles_readMarbles_success(t *testing.T) {
	fmt.Println("[TEST] transferMarbles")

	// invoke initMarbles
	stub := initMarble(t)

	// invoke transfer
	arguments := [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(sender), []byte(receiver), []byte(strconv.Itoa(transferAmount))}
	util.CheckInvoke(t, stub, arguments)

	// check sender amount state
	util.CheckState(t, stub, sender + sampleMarble.Name, strconv.Itoa(totalAmount - transferAmount))

	// check receiver amount state
	util.CheckState(t, stub, receiver + sampleMarble.Name, strconv.Itoa(transferAmount))
}

func Test_MARBLES_readMarbles_success(t *testing.T) {
	fmt.Println("[TEST] readMarbles")

	// invoke initMarbles
	stub := initMarble(t)

	// check receiver query (check amount 0)
	receiverResult := &marbleResponse{sampleMarble, receiver, 0}
	receiverResultBytes, _ := json.Marshal(receiverResult)
	arguments := [][]byte{[]byte(FUNCTION_READ), []byte(sampleMarble.Name), []byte(receiver)}
	util.CheckQuery(t, stub, arguments, string(receiverResultBytes))

	// invoke transfer
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(sender), []byte(receiver), []byte(strconv.Itoa(transferAmount))}
	util.CheckInvoke(t, stub, arguments)

	// check sender query
	senderResult := &marbleResponse{sampleMarble, sender, totalAmount - transferAmount}
	senderResultBytes, _ := json.Marshal(senderResult)
	arguments = [][]byte{[]byte(FUNCTION_READ), []byte(sampleMarble.Name), []byte(sender)}
	util.CheckQuery(t, stub, arguments, string(senderResultBytes))

	// check receiver query
	receiverResult = &marbleResponse{sampleMarble, receiver, transferAmount}
	receiverResultBytes, _ = json.Marshal(receiverResult)
	arguments = [][]byte{[]byte(FUNCTION_READ), []byte(sampleMarble.Name), []byte(receiver)}
	util.CheckQuery(t, stub, arguments, string(receiverResultBytes))
}