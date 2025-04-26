mock:
	go generate -v ./...

mock-svc:
	mockgen -source=internal/service/memberships/service.go -destination=internal/service/memberships/service_mock.go -package=memberships

run:
	go run cmd/main.go