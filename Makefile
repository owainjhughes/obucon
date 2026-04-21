.PHONY: dev backend frontend

dev:
	$(MAKE) -j2 backend frontend

backend:
	cd backend && go run cmd/server/main.go

frontend:
	cd frontend && npm start
