# grpcurl -plaintext  -d '{ "time": "", "numFromTail": 0, "cursor": "", "field": "GID", "value": "1000", "path": "" }'  localhost:50051 journal.v1.JournalService/Action

# grpcurl -plaintext  -d '{"numFromTail":10,"field":"Systemd","value":"systemd-journald.service"}'  localhost:50051  journal.v1.JournalService/Action
#
# grpcurl -plaintext \
#   -d '{"numFromTail":10000,"field":"Systemd","value":"systemd-journald.service"}' \
#   localhost:50051 \
#   journal.v1.JournalService/Action
#
grpcurl -plaintext  -d '{"numFromTail":5,"field":"Systemd","value":"systemd-journald.service"}'  localhost:5000  journal.v1.JournalService/Action
