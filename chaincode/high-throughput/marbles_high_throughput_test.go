package general

import (
	"encoding/json"
	"fmt"
	"marbles-meetup/util"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var sampleMarble = &marble{"marble", "RedMarble", "red", 30}
const txInit, txTransfer1, txTransfer2,  txTransfer3, txTransfer4, txPrune = "transfer_init", "transfer1", "transfer2", "transfer3", "transfer4", "transfer_prune"
const alice, bob, carol = "alice", "bob", "carol"
const totalAmount = 100000
const transferAmount1, transferAmount2, transferAmount3, transferAmount4 = 1000, 20, 40, 50

func checkAmount(t *testing.T, stub *shim.MockStub, marbleName, owner string, expectedAmount int)  {
	amount, err := getAmount(stub, marbleName, owner)
	if err != nil {
		fmt.Println("Fail to get Amount")
		t.FailNow()
	}
	if amount != (expectedAmount) {
		fmt.Println("amount for " + owner, amount, "was not", expectedAmount, "as expected")
		t.FailNow()
	}
}

func initMarble(t *testing.T) *shim.MockStub {
	var scc = new(HighThroughputChaincode)
	var stub = shim.NewMockStub("marbles", scc)
	stub.MockInit("1", [][]byte{[]byte("init")})
	arguments := [][]byte{[]byte(FUNCTION_INIT), []byte(sampleMarble.Name),
		[]byte(sampleMarble.Color), []byte(strconv.Itoa(sampleMarble.Size)),
		[]byte(strconv.Itoa(totalAmount)), []byte(alice)}
	util.CheckInvoke(t, stub, arguments, txInit)

	return stub
}

func Test_MARBLES_initMarble_success(t *testing.T) {
	fmt.Println("[TEST] initMarble")

	// invoke initMarbles
	stub := initMarble(t)

	// check marble state
	sampleMarbleBytes, _ := json.Marshal(sampleMarble)
	util.CheckState(t, stub, sampleMarble.Name, string(sampleMarbleBytes))

	// check state
	key, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, "", alice, strconv.Itoa(totalAmount), txInit})
	util.CheckState(t, stub, key, string([]byte{0x00}))

	// check amount
	checkAmount(t, stub, sampleMarble.Name, alice, totalAmount)
}

func Test_MARBLES_transferMarbles_success(t *testing.T) {
	fmt.Println("[TEST] transferMarbles")

	// invoke initMarbles
	stub := initMarble(t)

	// invoke transfer
	arguments := [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(alice), []byte(bob), []byte(strconv.Itoa(transferAmount1))}
	util.CheckInvoke(t, stub, arguments, txTransfer1)

	// check added state
	key, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, alice, bob, strconv.Itoa(transferAmount1), txTransfer1})
	util.CheckState(t, stub, key, string([]byte{0x00}))

	// check sender amount
	checkAmount(t, stub, sampleMarble.Name, alice, totalAmount -transferAmount1)

	// check receiver amount
	checkAmount(t, stub, sampleMarble.Name, bob, transferAmount1)
}

func Test_MARBLES_readMarbles_success(t *testing.T) {
	fmt.Println("[TEST] readMarbles")

	// invoke initMarbles
	stub := initMarble(t)

	// check receiver query (check amount 0)
	receiverResult := &marbleResponse{sampleMarble, bob, 0}
	receiverResultBytes, _ := json.Marshal(receiverResult)
	arguments := [][]byte{[]byte(FUNCTION_READ), []byte(sampleMarble.Name), []byte(bob)}
	util.CheckQuery(t, stub, arguments, string(receiverResultBytes), "1")

	// invoke transfer
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(alice), []byte(bob), []byte(strconv.Itoa(transferAmount1))}
	util.CheckInvoke(t, stub, arguments, "1")

	// check sender query
	senderResult := &marbleResponse{sampleMarble, alice, totalAmount - transferAmount1}
	senderResultBytes, _ := json.Marshal(senderResult)
	arguments = [][]byte{[]byte(FUNCTION_READ), []byte(sampleMarble.Name), []byte(alice)}
	util.CheckQuery(t, stub, arguments, string(senderResultBytes), "1")

	// check receiver query
	receiverResult = &marbleResponse{sampleMarble, bob, transferAmount1}
	receiverResultBytes, _ = json.Marshal(receiverResult)
	arguments = [][]byte{[]byte(FUNCTION_READ), []byte(sampleMarble.Name), []byte(bob)}
	util.CheckQuery(t, stub, arguments, string(receiverResultBytes), "1")
}

func Test_MARBLES_multi_transferMarbles_success(t *testing.T)  {
	// invoke initMarbles
	stub := initMarble(t)

	// invoke transfer1 alice -> bob
	arguments := [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(alice), []byte(bob), []byte(strconv.Itoa(transferAmount1))}
	util.CheckInvoke(t, stub, arguments, txTransfer1)

	// invoke transfer2 bob -> carol
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(bob), []byte(carol), []byte(strconv.Itoa(transferAmount2))}
	util.CheckInvoke(t, stub, arguments, txTransfer2)

	// invoke transfer3 alice -> carol
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(alice), []byte(carol), []byte(strconv.Itoa(transferAmount3))}
	util.CheckInvoke(t, stub, arguments, txTransfer3)

	// invoke transfer4 bob -> alice
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(bob), []byte(alice), []byte(strconv.Itoa(transferAmount4))}
	util.CheckInvoke(t, stub, arguments, txTransfer4)

	// check alice amount
	checkAmount(t, stub, sampleMarble.Name, alice, totalAmount - transferAmount1 - transferAmount3 + transferAmount4)

	// check bob amount
	checkAmount(t, stub, sampleMarble.Name, bob, transferAmount1 - transferAmount2 - transferAmount4)

	// check carol amount
	checkAmount(t, stub, sampleMarble.Name, carol, transferAmount2 + transferAmount3)
}

func Test_MARBLES_pruneMarbles_success(t *testing.T) {
	// invoke initMarbles
	stub := initMarble(t)

	// invoke transfer1 alice -> bob
	arguments := [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(alice), []byte(bob), []byte(strconv.Itoa(transferAmount1))}
	util.CheckInvoke(t, stub, arguments, txTransfer1)

	// invoke transfer2 bob -> carol
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(bob), []byte(carol), []byte(strconv.Itoa(transferAmount2))}
	util.CheckInvoke(t, stub, arguments, txTransfer2)

	// invoke transfer3 alice -> carol
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(alice), []byte(carol), []byte(strconv.Itoa(transferAmount3))}
	util.CheckInvoke(t, stub, arguments, txTransfer3)

	// invoke transfer4 bob -> alice
	arguments = [][]byte{[]byte(FUNCTION_TRANSFER), []byte(sampleMarble.Name),
		[]byte(bob), []byte(alice), []byte(strconv.Itoa(transferAmount4))}
	util.CheckInvoke(t, stub, arguments, txTransfer4)

	// result amount
	aliceAmount := totalAmount - transferAmount1 - transferAmount3 + transferAmount4
	bobAmount := transferAmount1 - transferAmount2 - transferAmount4
	carolAmount := transferAmount2 + transferAmount3

	// check State
	keyInit, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, "", alice, strconv.Itoa(totalAmount), txInit})
	util.CheckState(t, stub, keyInit, string([]byte{0x00}))

	keyTransfer1, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, alice, bob, strconv.Itoa(transferAmount1), txTransfer1})
	util.CheckState(t, stub, keyTransfer1, string([]byte{0x00}))

	keyTransfer2, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, bob, carol, strconv.Itoa(transferAmount2), txTransfer2})
	util.CheckState(t, stub, keyTransfer2, string([]byte{0x00}))

	keyTransfer3, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, bob, carol, strconv.Itoa(transferAmount2), txTransfer2})
	util.CheckState(t, stub, keyTransfer3, string([]byte{0x00}))

	keyTransfer4, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, bob, carol, strconv.Itoa(transferAmount2), txTransfer2})
	util.CheckState(t, stub, keyTransfer4, string([]byte{0x00}))

	// invoke pruneMarbles
	arguments = [][]byte{[]byte(FUNCTION_PRUNE), []byte(sampleMarble.Name)}
	util.CheckInvoke(t, stub, arguments, txPrune)

	// check Del State
	util.CheckStateNotExisted(t, stub, keyInit)
	util.CheckStateNotExisted(t, stub, keyTransfer1)
	util.CheckStateNotExisted(t, stub, keyTransfer2)
	util.CheckStateNotExisted(t, stub, keyTransfer3)
	util.CheckStateNotExisted(t, stub, keyTransfer4)

	// check new State
	keyAlice, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, "", alice, strconv.Itoa(aliceAmount), txPrune})
	util.CheckState(t, stub, keyAlice, string([]byte{0x00}))

	keyBob, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, "", bob, strconv.Itoa(bobAmount), txPrune})
	util.CheckState(t, stub, keyBob, string([]byte{0x00}))

	keyCarol, _ := stub.CreateCompositeKey(KEY_TRANSFER, []string{sampleMarble.Name, "", carol, strconv.Itoa(carolAmount), txPrune})
	util.CheckState(t, stub, keyCarol, string([]byte{0x00}))

	// check amount is equal
	checkAmount(t, stub, sampleMarble.Name, alice, aliceAmount)
	checkAmount(t, stub, sampleMarble.Name, bob,bobAmount)
	checkAmount(t, stub, sampleMarble.Name, carol, carolAmount)
}