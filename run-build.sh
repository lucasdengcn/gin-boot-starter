go env -w GOPATH=/Users/yamingdeng/goprojects
go build -o gin-app main.go
./gin-app -w /Users/yamingdeng/goprojects/src/gin-boot-starter 2>&1 | tee log.txt