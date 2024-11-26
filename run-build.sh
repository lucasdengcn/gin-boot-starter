go env -w GOPATH=/Users/yamingdeng/goprojects
go build
./gin001 -cfg /Users/yamingdeng/goprojects/src/gin001 2>&1 | tee log.txt