# REPLIT

Your REPL to k6.

## Run from CLI

```bash
go build ./cmd/replit
# run
./replit
```

## Run as an extension

```bash
xk6 run -q --no-summary examples/replit.js
```

## Build the extension

```bash
xk6 build --with github.com/grafana/xk6-replit=.
# run
./k6 run -q --no-summary examples/replit.js
```

