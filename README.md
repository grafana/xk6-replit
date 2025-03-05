# REPLIT

Your REPL to k6.

## Using REPLIT inside an existing k6 test

You can drop into a REPL from within an existing k6 test via the following methods.
You will need to modify your script so that it includes the `replit` import, and await the `replit.with` function.
The `replit.with` function takes an (optional) argument an object containing a context to make available for use in the REPL. This object can contain variables and modules. Here's an example:

```js
// You can find this example in examples/http.js
import { replit } from "k6/x/replit";
import http from "k6/http";

export default async function () {
    let result = http.get("https://quickpizza.grafana.com");

    // As context, we pass 'result', and the http module in case
    // we want to make additional requests.
    await replit.with({result: result, http: http});

    console.log("All done.")
}
```

> [!NOTE]
> Your exported default function will need to be `async` for `replit` to work.

## Running

You can run a standalone REPL (i.e. without any context) via the following methods:

### Run using CLI helper

```bash
go run ./cmd/replit examples/http.js
```

### Build the CLI tool and run it separately

```bash
go build ./cmd/replit
./replit examples/http.js
```

> [!NOTE]
> The CLI tool requires [Go](https://go.dev/doc/install) and [xk6](https://github.com/grafana/xk6) to be installed.


### Run as an extension

```bash
xk6 run -q --no-summary examples/http.js
```

### Build the extension and run the k6 binary separately

```bash
xk6 build --with github.com/grafana/xk6-replit=.
./k6 run -q --no-summary examples/http.js
```
