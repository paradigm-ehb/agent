grpcurl -plaintext \
  -d '{
    "time": "",
    "numFromTail": 0,
    "cursor": "",
    "field": "GID",
	"value": "1000",
    "path": ""
  }' \
  localhost:50051 journal.v1.JournalService/Action
