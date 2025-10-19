# About

A project to simulate a high number of Devices that are creating Cumulocity MEAs (Measurements, Events, Alarms) in a certain frequency. 

# How to use this project

* Clone the project

* Configure number of Devices, device names/serials and the sending interval in `config.toml` file

* You can define the data sent in each interval in `collectFunctions()` function in `/pkg/app/device.go`

* To build the project, have a look in the `justfile`. A build for macOs on ARM CPUs would be `CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o c8y-device-simulator main.go`, it will place an executable in your current directory

* Now create below `.env` file in your current directory (it tells the script your Cumulocity Tenant and User):

```sh
C8Y_HOST=example.cumulocity.com
C8Y_TENANT=t1234
C8Y_USER=john.doe@cumulocity.com
C8Y_PASSWORD=super-secret-password
```

* Now you're all set, do `./c8y-device-simulator` to start the project. Make sure to have `config.toml` and `.env` file in the same directory where you start this executable from.

# Deploy as Cumulocity Microservices

The project can be deployed as Cumulocity via `just build-amd64-image && just deploy-microservice`. This requires [just](https://github.com/casey/just) and [go-c8y-cli](https://goc8ycli.netlify.app/docs/introduction/) being installed. 

If you do not have these tools installed, no problem, just take a look in the project roots `justfile` to see the necessary steps to build and upload the image.