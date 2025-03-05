# REPLIT

Your REPL to k6.

## Running a standalone REPL

You can run a standalone REPL (i.e. without any context) via the following methods:

### Run using CLI helper

```bash
go run ./cmd/replit
```

### Run as an extension

```bash
xk6 run -q --no-summary cmd/replit/assets/repl.js
```

### Build the extension and run the k6 binary separately

```bash
xk6 build --with github.com/grafana/xk6-replit=.
./k6 run -q --no-summary cmd/replit/assets/repl.js
```

## Running a REPL inside an existing k6 test

You can drop into a REPL from within an existing k6 test via the following methods.
You will need to modify your script so that it includes the `replit` import, and await the `replit.repl` function.
The `replit.repl` function takes as its first argument the `replit` module itself, and as a second (optional) argument an object containing a context to make available for use in the REPL. This object can contain variables and modules. Here's an example:

```js
import { replit } from "k6/x/replit";
import http from "k6/http";

export default async function () {
    let result = http.get("https://quickpizza.grafana.com");

    // The first argument to 'repl' is the replit module itself.
    // As context, we pass 'result', and the http module in case
    // we want to make additional requests.
    await replit.repl(replit, {result: result, http: http});

    console.log("All done.")
}
```

> [!NOTE]
> Your exported default function will need to be `async` for `replit` to work.

### Run as an extension

```bash
xk6 run -q --no-summary my_script.js
```

### Build the extension and run the k6 binary separately

```bash
xk6 build --with github.com/grafana/xk6-replit=.
./k6 run -q --no-summary my_script.js
```
