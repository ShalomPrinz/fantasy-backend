export FIRESTORE_EMULATOR_HOST=127.0.0.1:8090
export FIREBASE_AUTH_EMULATOR_HOST=127.0.0.1:8110

go test ./... -coverprofile=./coverage/cover.out -coverpkg ../...
go tool cover -o ./coverage/coverage.html -html=./coverage/cover.out