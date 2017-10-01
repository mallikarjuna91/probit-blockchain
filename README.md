## Probit blockchain app

This code is taken from the fabric-sample module of hyperledger to run the smart contract locally.

The smart contract that is written in go code is present in the **/artifacts/src/github.com/probit/** folder.

It needs to be present in the fabric-samples folder to generate and access the certificates and to start the network.

## Running the sample program

##### Terminal Window 1

* Launch the network using docker-compose

```
docker-compose -f artifacts/docker-compose.yaml up
```
##### Terminal Window 2

* Install the fabric-client and fabric-ca-client node modules

```
npm install
```

* Start the node app on PORT 4000

```
PORT=4000 node app
```



## Sample REST APIs Requests

### Login Request

* Register and enroll new users in Organization - **Org1**:

`curl -s -X POST http://localhost:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=Jim&orgName=org1'`

**OUTPUT:**

```
{
  "success": true,
  "secret": "RaxhMgevgJcm",
  "message": "Jim enrolled Successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI"
}
```

The response contains the success/failure status, an **enrollment Secret** and a **JSON Web Token (JWT)** that is a required string in the Request Headers for subsequent requests.

### Create Channel request

```
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../artifacts/channel/mychannel.tx"
}'
```

Please note that the Header **authorization** must contain the JWT returned from the `POST /users` call

### Join Channel request

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"]
}'
```
### Install chaincode

```
curl -s -X POST  \
    http://localhost:4000/chaincodes \
    -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY4ODg2MjYsInVzZXJuYW1lIjoibWFsbGlrIiwib3JnTmFtZSI6Im9yZzIiLCJpYXQiOjE1MDY4NTI2MjZ9.rTFB-e1h30K73lGbBOkQuCK5HtJFGZf6A4YWRoFS_R8"  \
    -H "content-type: application/json"   \
    -d '{
    "peers": ["peer1","peer2"],
    "chaincodeName":"hackathon",
    "chaincodePath":"github.com/probit",
    "chaincodeVersion":"v0"
}'
```

### Instantiate chaincode

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"chaincodeName":"hackathon",
	"chaincodeVersion":"v0",
	"args":["arjun","1000","probit","2000"]
}'
```

### Invoke request for buy shares
```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/hackathon \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY4ODg2MjYsInVzZXJuYW1lIjoibWFsbGlrIiwib3JnTmFtZSI6Im9yZzIiLCJpYXQiOjE1MDY4NTI2MjZ9.rTFB-e1h30K73lGbBOkQuCK5HtJFGZf6A4YWRoFS_R8" \
  -H "content-type: application/json" \
  -d '{
	"fcn":"buyShares",
	"args":["nishamallesh","probit","apple","0.7","700"]
}'
```

### Invoke request for sell shares

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/hackathon \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY4ODg2MjYsInVzZXJuYW1lIjoibWFsbGlrIiwib3JnTmFtZSI6Im9yZzIiLCJpYXQiOjE1MDY4NTI2MjZ9.rTFB-e1h30K73lGbBOkQuCK5HtJFGZf6A4YWRoFS_R8" \
  -H "content-type: application/json" \
  -d '{
	"fcn":"sellShares",
	"args":["nishamallesh","probit","apple","0.7","700"]
}'
```

### Invoke request for adding user
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/hackathon \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY4ODg2MjYsInVzZXJuYW1lIjoibWFsbGlrIiwib3JnTmFtZSI6Im9yZzIiLCJpYXQiOjE1MDY4NTI2MjZ9.rTFB-e1h30K73lGbBOkQuCK5HtJFGZf6A4YWRoFS_R8" \
  -H "content-type: application/json" \
  -d '{
	"fcn":"addUser",
	"args":["rob","600"]
}'

### Chaincode Query

```
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=query&args=%5B%22arjun%22%5D" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```
