echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

grpcurl -plaintext -d '{ "time": "1m", "numFromTail": 10, "cursor": "", "field": "Systemd", "value": "ssh.service", "path": "" }' localhost:50051  journal.v1.JournalService/Action
