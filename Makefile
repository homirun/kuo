install: ./cmd/kuo/main.go
	go build -o kubectl-kuo ./cmd/kuo
	mv ./kubectl-kuo /usr/local/bin/
