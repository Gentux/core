#!/bin/bash

JSON='{
    "jsonrpc" : "2.0",
    "id"      : "1",
    "method"  : "ServiceUsers.GetList",
    "params"  : [ { } ]
}'

ENDPOINT='http://localhost:8081/rpc'

curl --data-binary "$JSON" -H 'content-type:application/json;' $ENDPOINT
