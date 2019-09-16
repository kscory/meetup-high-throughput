#!/usr/bin/env bash

CC_NAME_THROUGHPUT="marblehighthroughput"

echo "GET query readMarbles chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/$CC_NAME_THROUGHPUT?peer=peer0.org1.example.com&fcn=readMarbles&args=%5B%22redMarbles%22%2C%22alice%22%5D" \
  -H "authorization: Bearer pass" \
  -H "content-type: application/json"
echo
echo

echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/$CC_NAME_THROUGHPUT?peer=peer0.org1.example.com&fcn=readMarbles&args=%5B%22redMarbles%22%2C%22bob%22%5D" \
  -H "authorization: Bearer pass" \
  -H "content-type: application/json"
echo
echo

echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/$CC_NAME_THROUGHPUT?peer=peer0.org1.example.com&fcn=readMarbles&args=%5B%22redMarbles%22%2C%22carol%22%5D" \
  -H "authorization: Bearer pass" \
  -H "content-type: application/json"
echo
echo
