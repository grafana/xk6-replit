# REPLIT

Your REPL to k6.

## Run from CLI

```bash
go build -o repl ./cmd/replit
# run
./repl
```

## Run as an extension

```bash
xk6 run -q --no-summary cmd/replit/assets/replit.js
```

## Build the extension

```bash
xk6 build --with github.com/grafana/xk6-replit=.
# run
./k6 run -q --no-summary cmd/replit/assets/replit.js
```

