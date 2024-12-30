# **crypto-utxo-gateway**

**crypto-utxo-gateway** is a service that acts as a gateway for UTXO (Unspent Transaction Output) chains, providing a unified RPC interface for interacting with multiple blockchain networks. It currently supports Bitcoin, Bitcoin Cash, Dash, Dogecoin, and Litecoin.

Written in **Go**, this service exposes a **gRPC** interface for seamless integration with upper-layer services.

---

## **Features**

- Supports multiple UTXO-based blockchains:
    - Bitcoin
    - Bitcoin Cash
    - Dash
    - Dogecoin
    - Litecoin
- Exposes a **gRPC** interface for service consumption.
- Written in **Go** for high performance and scalability.
- Easily extendable to support additional UTXO-based chains.

---

## **Installation**

### **Prerequisites**

- Go 1.18+
- Protobuf compiler (for generating gRPC files)

### **Steps to Install**

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/crypto-utxo-gateway.git
   ```

2. Navigate to the project directory:

   ```bash
   cd crypto-utxo-gateway
   ```

3. Install the dependencies:

   ```bash
   go mod tidy
   ```

4. Compile the gRPC proto files:

   ```bash
   make generate
   ```

5. Build the service:

   ```bash
   make build
   ```

6. Run the service:

   ```bash
   ./crypto-utxo-gateway
   ```

---

## **Usage**

### **gRPC Interface**

The service exposes gRPC endpoints to interact with the supported blockchains. The gRPC definition file (`proto/crypto_utxo_gateway.proto`) is included in the repository. This file should be used to generate client stubs.

### **Example Request**

Once the service is running, you can query UTXO data for a specific blockchain. Here's an example of using `grpcurl` to interact with the service:

```bash
grpcurl -d '{"address": "your_bitcoin_address"}' -proto proto/crypto_utxo_gateway.proto \
  localhost:50051 yourservice.UTXOService.GetUTXOs
```

### **Supported Chains**

Specify the chain you want to query in your request. Supported chains include:

- `bitcoin`
- `bitcoincash`
- `dash`
- `dogecoin`
- `litecoin`

Each request will return UTXO information for the specified chain.

---

## **Development**

To contribute to the project:

1. Fork the repository.
2. Clone your fork:

   ```bash
   git clone https://github.com/yourusername/crypto-utxo-gateway.git
   ```

3. Create a new branch for your feature:

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. Make your changes and commit:

   ```bash
   git commit -m "Add feature"
   ```

5. Push your changes:

   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a pull request.

---

## **Testing**

### **Unit Tests**

Run unit tests using the following command:

```bash
go test ./...
```

### **Integration Tests**

Integration tests for interacting with live blockchains are located in the `integration_tests` directory. These tests may require running test blockchain nodes.

---

## **API Documentation**

For complete API documentation, refer to the `proto/crypto_utxo_gateway.proto` file.

---

## **Contributing**

We welcome contributions! To contribute:

- Fork the repository.
- Submit a pull request for new features or bug fixes.

For significant changes, please open an issue first to discuss the intended changes.

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

---

## **Acknowledgements**

- [gRPC](https://grpc.io/) - A high-performance RPC framework.
- [Protobuf](https://developers.google.com/protocol-buffers) - Data serialization format used in gRPC.
- [Go](https://golang.org/) - The Go programming language used to implement the service.

---
