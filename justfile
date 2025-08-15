dev:
    bun run build
    go run . --dev --config test/config.toml --domain example.org --public-dir public/

build:
    bun run build
    go build