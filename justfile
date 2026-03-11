default:
    just --list

[parallel]
watch-and-generate: tailwind-watch templ

tailwind-watch:
    tailwindcss -i ./views/assets/css/input.css -o ./views/assets/css/output.css --watch

templ:
    go tool templ generate --watch --proxy="http://localhost:8090" --cmd="go run ./cmd/server/main.go" --open-browser=false

# Start dev server
dev:
    rm -f ./views/assets/css/output.css
    just watch-and-generate

# Build the server and CLI
build:
    go tool templ generate
    tailwindcss -i ./views/assets/css/input.css -o ./views/assets/css/output.css --minify
    go build -o dist/server ./cmd/server
    go install ./cmd/blkhell

# Run the server
serve: build
    ./dist/server

# Run tests
test:
    go test ./...
