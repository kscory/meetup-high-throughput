# Fabric High Throughput

This is Fabric High Throuput Code for Hyperledger Meetup

## Prerequisites and setup:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Node.js** v8.9.x
* * **Go** v1.11.x
* [Download Docker images](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#binaries)

## Test Chaincode without Fabric

move to chaincode directory

```
$ cd chaincode
```

install & remove modules with go mod

```
$ go mod tidy
```

test general chaincode:

```
$ go test ./general
```

test high throughput Chaincode:

```
$ go test ./high-throughput
```

test high throughput removed phantom read Chaincode:

```
$ go test ./high-throughput-phantom
```


## Test Chaincode with SDK

test chaincode with fabric & nodejs sdk

### Running the sample project

move to server root directory

```
$ cd server
```

Installs the fabric-client and fabric-ca-client node modules

starts the node app on PORT 4000

```
$ ./scripts/runApp.sh
```

preInstall to test chaincode

```
$ ./scripts/preInstall.sh
```

### Test general chaincode:

init Marbles & transfer Marbles

```
$ ./scripts/runGeneral.sh
```

read Marbles result

```
$ ./scripts/readGeneralMarbles.sh
```

### Test high throughput Chaincode:

init Marbles & transfer Marbles

```
$ ./scripts/runHighThroughput.sh
```

read Marbles result

```
$ ./scripts/readHighThroughputMarbles.sh
```

### Test high throughput removed phantom read Chaincode:

init Marbles & transfer Marbles

```
$ ./scripts/runHighThroughputPhantom.sh
```

read Marbles result

```
$ ./scripts/readHighThroughputPhantomMarbles.sh
```

### Test high throughput removed phantom read Chaincode with solution2:

init Marbles & transfer Marbles

```
$ ./scripts/runSolution.sh
```

read Marbles result

```
$ ./scripts/readSolution.sh
```