# Run the gateway-service
run:
	go run ./cmd/main.go

# Run k6 load test
k6:
	k6 run k6-orders.js



