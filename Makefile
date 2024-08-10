# Run the application
run:
	@clear
	@go run cmd/api/main.go

start:
	./bin/web

build: ## build to binary files
	go build -o ./bin/web ./cmd/api/main.go

clean: ## cleaning project
	@echo "cleaning..."
	# remove binary files
	rm -rf ./bin
	# doc: https://pkg.go.dev/cmd/go/internal/clean
	go clean

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: run build clean watch

