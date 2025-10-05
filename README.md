# Fabric Financial Asset Management System
A complete, end-to-end blockchain application built with Hyperledger Fabric and Go. This system provides a secure, transparent, and immutable ledger for managing financial accounts, featuring a custom smart contract and a containerized REST API.

# 🏛️ Project Architecture
This project consists of three main components that work together to create a robust asset management system:

Hyperledger Fabric Network: Managed by the test-network script from fabric-samples. It consists of two peer organizations and an ordering service, all running in Docker containers. This forms the decentralized backbone of the application.

Go Smart Contract (Chaincode): The business logic deployed on the Fabric network. It defines the rules and functions for interacting with the ledger, such as creating and querying financial accounts.

Go REST API Server: A client-facing web server that acts as a gateway to the blockchain. It uses the Fabric Gateway SDK to securely connect to the network and invoke smart contract transactions. The API is containerized with Docker for portability.

# ✨ Core Features
🔐 Immutable Ledger: All account records and transactions are permanent and tamper-proof.

✅ Asset Management: Full CRUD (Create, Read, Update) functionality for financial accounts on the blockchain.

📖 Transaction History: Ability to retrieve a full, verifiable audit trail for any asset.

🌐 RESTful API: A simple and powerful interface for client applications to interact with the blockchain.

🐳 Containerized: The entire application is containerized with Docker for consistency and ease of deployment.

🚀 Quick Start Guide
This guide will walk you through setting up and running the entire application.

# Prerequisites
Docker & Docker Compose

Go Programming Language (v1.24+)

Git

The official Hyperledger fabric-samples repository cloned locally.

Step 1: Launch the Blockchain Network (Terminal 1)
First, we start the Hyperledger Fabric network and deploy our smart contract.

Navigate to the test-network directory:

cd ~/fabric-internship-project/fabric-samples/test-network

Clean up any old instances and start a fresh network:

./network.sh down
./network.sh up createChannel -ca

Deploy the chaincode to the network. Note: Since this is a new network, the sequence number must be 1.

./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go -ccv 1.0 -ccs 1

<!-- INSERT SCREENSHOT HERE: Successful chaincode deployment, ending with "Chaincode initialization is not required" -->

<!-- Suggested file: Screenshot 2025-10-05 121752.jpg -->

Leave this terminal running. It is now your live blockchain network.

Step 2: Configure & Run the API Server (Terminal 2)
Next, we configure and launch our API server, which will connect to the blockchain.

Open a new, second terminal. Find your machine's local IP address, as the Docker container will use this to connect to the peer.

ip addr show eth0 | grep "inet\s" | awk '{print $2}' | cut -d/ -f1

(Copy the IP address, e.g., 172.25.1.38)

Navigate to the financial-api directory and update the main.go file with this IP address:

cd ~/fabric-internship-project/fabric-samples/financial-api
nano main.go

In the editor, find the peerEndpoint constant and set its value:

// Example
const peerEndpoint = "172.25.1.38:7051"

Save and exit (Ctrl + X, Y, Enter).

Build and run the Docker container for the API server:

docker build -t financial-api .
docker run --rm -p 8080:8080 --name financial-api \
-v ${PWD}/../test-network:/test-network \
financial-api

You will see the log Starting server on port 8080.

<!-- INSERT SCREENSHOT HERE: The API server logs showing it has started successfully -->

<!-- Suggested file: Screenshot 2025-10-05 121942.png -->

Leave this second terminal running. It is your live API server.

Step 3: Test the Live Application (Terminal 3)
Finally, let's interact with our running application.

Open a third, new terminal.

Use curl to send requests to your API.

Create an Asset (POST /assets)
curl -X POST http://localhost:8080/assets -d '{
    "DEALERID": "DEALER001", "MSISDN": "9876543210", "MPIN": "1234",
    "BALANCE": 5000, "STATUS": "ACTIVE", "TRANSAMOUNT": 5000,
    "TRANSTYPE": "CREDIT", "REMARKS": "Initial deposit via API"
}'

Read the Asset (GET /assets)
curl http://localhost:8080/assets?dealerid=DEALER001

You should see a successful response for both commands, proving the end-to-end flow is working correctly.

<!-- INSERT SCREENSHOT HERE: The third terminal showing the successful curl commands and their output -->

<!-- Suggested file: Screenshot 2025-10-05 121633.png -->

► Shutdown
To stop the application, stop the processes in reverse order:

Stop the API server (Terminal 2) with Ctrl + C.

Stop the Fabric network (Terminal 1) with the command:

./network.sh down
