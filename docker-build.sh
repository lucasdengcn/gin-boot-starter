#docker buildx build . --output type=tar,dest=./build_cache.tar

docker buildx build . -t localdev/gin-api-service:latest