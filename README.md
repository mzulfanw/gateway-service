# Gateway Service

This is the Gateway Service for our application. It acts as a single entry point for all client requests and routes them to the appropriate backend services.

## Architecture Overview
The gateway uses HTTP to receive client requests and forwards them to internal services (e.g., order-service, product-service) based on the request path. It can handle  logging, and request validation centrally.

## Prerequisites
- Go 1.20+
- Docker (optional, for containerization)
- Make sure you have the backend services ([order-service](https://github.com/mzulfanw/order-service), [product-service](https://github.com/mzulfanw/products-service)) running before starting the gateway.

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone git@github.com:mzulfanw/gateway-service.git
   cd gateway-service
    ```
2. Install dependencies:
    ```bash
    go mod tidy 
    ```
3. Configure environment variables:
    ```bash
   cp .env.example .env
   ```

4. Run the gateway service:
   ```bash
   go run cmd/main.go
   ```
   
The gateway will start on `http://localhost:3000`.

## API Endpoints
The gateway exposes the following endpoints:
- `POST /products`: Forwards requests to the product-service to create a new product.
- `GET /products/${id}`: Forwards requests to the product-service to retrieve product details by ID.
- `POST /orders`: Forwards requests to the order-service to create a new order.
- `GET /orders/product/${id}`: Forwards requests to the order-service to retrieve orders by product ID.

## Test with k6
You can use k6 to perform load testing on the gateway service. An example test script is provided in the `root` directory.
1. Install k6 if you haven't already. Follow the instructions at [k6 Installation](https://k6.io/docs/getting-started/installation/).
2. Run the k6 test script:
    ```bash
   k6 run k6-orders.js
   ```
   
## Postman Collection
A Postman collection is available to test the API endpoints. You can download the collection postman_collection.json and postman_environment.json files from the repository and import them into Postman.

