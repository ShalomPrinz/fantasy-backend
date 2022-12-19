export FIRESTORE_EMULATOR_HOST=127.0.0.1:8090
go test ./... -coverprofile=./coverage/cover.out -coverpkg ../...
go tool cover -o ./coverage/coverage.html -html=./coverage/cover.out