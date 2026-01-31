.PHONY: backend frontend docker-backend docker-frontend test load-test help

help:
	@echo "Texas Hold'em - Build targets"
	@echo "  make backend        - Run Go backend locally"
	@echo "  make test           - Run backend tests"
	@echo "  make frontend       - Run Flutter frontend (requires: cd frontend && flutter run -d chrome)"
	@echo "  make docker-backend - Build backend Docker image (linux/amd64)"
	@echo "  make docker-frontend - Build frontend Docker image (requires flutter build web first)"
	@echo "  make load-test      - Run k6 load test (start backend first)"

backend:
	cd backend && go run ./cmd/server

test:
	cd backend && go test ./...

docker-backend:
	docker build --platform linux/amd64 -t texas-holdem-backend:latest ./backend

docker-frontend:
	cd frontend && flutter build web && cd .. && docker build --platform linux/amd64 -t texas-holdem-frontend:latest ./frontend

load-test:
	cd load-test && k6 run load-test.js
