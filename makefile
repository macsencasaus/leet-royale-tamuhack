run: front
	cd backend && go run .

build: front
	cd backend && go build

front:
	cd frontend && npm install --dev && npm run build

fmt:
	cd backend && gofmt -l -s -w .

clean:
	rm backend/leet-guys

.PHONY: run front build fmt clean
