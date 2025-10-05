package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Fabric network details
const (
	mspID         = "Org1MSP"
	cryptoPath    = "/test-network/organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath       = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "172.25.1.38:7051" // <-- FIX #1: Use your direct IP address
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "mychannel"
	chaincodeName = "basic"
)

// Asset describes the financial account details.
type Asset struct {
	DEALERID    string  `json:"DEALERID"`
	MSISDN      string  `json:"MSISDN"`
	MPIN        string  `json:"MPIN"`
	BALANCE     float64 `json:"BALANCE"`
	STATUS      string  `json:"STATUS"`
	TRANSAMOUNT float64 `json:"TRANSAMOUNT"`
	TRANSTYPE   string  `json:"TRANSTYPE"`
	REMARKS     string  `json:"REMARKS"`
}

var contract *client.Contract

func main() {
	log.Println("============ application-golang starts ============")

	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to gateway: %w", err))
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract = network.GetContract(chaincodeName)

	http.HandleFunc("/assets", assetsHandler)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAsset(w, r)
	case "POST":
		createAsset(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	dealerId := r.URL.Query().Get("dealerid")
	if dealerId == "" {
		http.Error(w, "dealerid query parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("--> Evaluate Transaction: ReadAsset, for dealer %s", dealerId)
	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", dealerId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to evaluate transaction: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(evaluateResult)
}

func createAsset(w http.ResponseWriter, r *http.Request) {
	var asset Asset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("--> Submit Transaction: CreateAsset, for dealer %s", asset.DEALERID)
	_, err := contract.SubmitTransaction("CreateAsset", asset.DEALERID, asset.MSISDN, asset.MPIN, fmt.Sprintf("%f", asset.BALANCE), asset.STATUS, fmt.Sprintf("%f", asset.TRANSAMOUNT), asset.TRANSTYPE, asset.REMARKS)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit transaction: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Asset %s created successfully", asset.DEALERID)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials), grpc.WithAuthority(gatewayPeer)) // <-- FIX #2: Add WithAuthority
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read keystore directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func loadCertificate(path string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}