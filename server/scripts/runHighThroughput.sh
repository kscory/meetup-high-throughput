#!/usr/bin/env bash

CC_NAME_THROUGHPUT="marblehighthroughput"

function transfer(){
    echo "transfer marbles to $1"
    echo
    VALUES=$(curl -s -X POST \
      http://localhost:4000/channels/mychannel/chaincodes/${CC_NAME_THROUGHPUT} \
      -H "content-type: application/json" \
      -H "authorization: Bearer pass" \
      -d "{
      \"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
      \"fcn\":\"transferMarbles\",
      \"args\":[\"redMarbles\",\"alice\",\"$1\",\"20000\"]
    }")
    echo $VALUES
    echo "this process transfer to is $1"
    # Assign previous invoke transaction id  to TRX_ID
    MESSAGE=$(echo $VALUES | jq -r ".message")
    TRX_ID=${MESSAGE#*ID:}
    echo
}

echo "POST invoke initMarbles chaincode on peers of Org1 and Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/${CC_NAME_THROUGHPUT} \
  -H "content-type: application/json" \
  -H "authorization: Bearer pass" \
  -d "{
	\"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
	\"fcn\":\"initMarbles\",
	\"args\":[\"redMarbles\",\"red\",\"30\",\"100000\",\"alice\"]
}"
echo
echo

echo "POST invoke transferMarbles chaincode 10 times on peers of Org1 and Org2"
for ((i = 1; i < 6; ++i))
do
	transfer bob &
	transfer carol &
done
