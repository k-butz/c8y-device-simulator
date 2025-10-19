set dotenv-load

run:
    go run main.go

build:
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o c8y-device-simulator main.go

build-amd64-image TAG="latest":
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o c8y-device-simulator main.go
    docker buildx build --platform linux/amd64 -t c8y-device-simulator:{{TAG}} .
    docker save c8y-device-simulator:{{TAG}} > image.tar
    zip c8y-device-simulator.zip image.tar cumulocity.json

deploy-microservice:
    c8y microservices create --file c8y-device-simulator.zip -f