Fabric Financial Asset Management System
A complete, end-to-end blockchain application built with Hyperledger Fabric and Go. This system provides a secure, transparent, and immutable ledger for managing financial accounts, featuring a custom smart contract and a containerized REST API.

► Core Features
Immutable Asset Creation: Add new financial accounts to the blockchain.

State Management: Update asset details and query the current state of any account.

Transaction History: Retrieve a full, verifiable history of all changes for a specific asset.

RESTful API: A simple, powerful interface for interacting with the blockchain network.

Containerized: The entire application is containerized with Docker for consistency and ease of deployment.

► Tech Stack
Backend: Go

Blockchain: Hyperledger Fabric v2.5

Containerization: Docker, Docker Compose

► Quick Start Guide
This guide assumes all prerequisites (Go, Docker, Git) are installed.

1. Launch the Network & Deploy Chaincode
Open a terminal, navigate to the fabric-samples/test-network directory, and run the following commands.

# Clean up any previous instances
./network.sh down

# Start the Fabric network and create a channel
./network.sh up createChannel -ca

# Deploy the smart contract
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go -ccv 1.0 -ccs 1

Leave this terminal open. It is now running your live blockchain network.

2. Configure and Run the API Server
Open a new terminal. Find and note your WSL IP address:

ip addr show eth0 | grep "inet\s" | awk '{print $2}' | cut -d/ -f1

Navigate to the API directory and update the main.go file with your IP address:

cd ../financial-api
# Open main.go and set `peerEndpoint` to your IP: "YOUR_IP:7051"
nano main.go 

Build and run the Docker container:

docker build -t financial-api .
docker run --rm -p 8080:8080 --name financial-api \
-v ${PWD}/../test-network:/test-network \
financial-api

Leave this second terminal running. It is your live API server.

3. Test the API
Open a third terminal to interact with your application.

Create an Asset (POST /assets)
curl -X POST http://localhost:8080/assets -d '{
    "DEALERID": "DEALER001", "MSISDN": "9876543210", "MPIN": "1234",
    "BALANCE": 5000, "STATUS": "ACTIVE", "TRANSAMOUNT": 5000,
    "TRANSTYPE": "CREDIT", "REMARKS": "Initial deposit via API"
}'

Success: Asset DEALER001 created successfully

Read an Asset (GET /assets)
curl http://localhost:8080/assets?dealerid=DEALER001

Success: {"DEALERID":"DEALER001","MSISDN":"9876543210",...}

► Shutdown
Stop the API server (Terminal 2) with Ctrl + C.

Stop the Fabric network (Terminal 1) with the command:

./network.sh down
