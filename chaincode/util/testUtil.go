package util

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CheckState(t *testing.T, stub *shim.MockStub, name string, value string) {
	byteState := stub.State[name]
	if byteState == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(byteState) != value {
		fmt.Println("State value", string(byteState), "was not", value, "as expected")
		t.FailNow()
	}
}

func CheckStateWithByte(t *testing.T, stub *shim.MockStub, name string, value []byte) {
	byteState := stub.State[name]
	if byteState == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if bytes.Compare(byteState, value) != 0 {
		fmt.Println("State value", byteState, "was not", value, "as expected")
		t.FailNow()
	}
}

func CheckQuery(t *testing.T, stub *shim.MockStub, args [][]byte, expect string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Query", args, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", args, "failed to get result")
		t.FailNow()
	}
	if string(res.Payload) != expect {
		fmt.Println("Query result ", string(res.Payload), "was not", expect, "as expected")
		t.FailNow()
	}
}

func CheckInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}