#!/bin/bash

JSON='{ 
    "jsonrpc" : "2.0",
    "id"      : "curltext",
    "method"  : "ServiceUsers.RegisterUser",
    "params"  : [ { "firstname" : "'$1'" } ]
}'

ENDPOINT="http://localhost:8081/rpc"

curl --data-binary "$JSON" -H 'content-type:application/json;' $ENDPOINT
