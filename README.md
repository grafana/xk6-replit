# REPLIT

Your [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) to [k6](https://github.com/grafana/k6).

REPLIT allows you to interact with your k6 script in real-time, inspecting variables, and running code snippets. This approach can be useful for debugging, testing, and learning k6 APIs.

- You write and send commands to k6 to execute and receive execution results.
- Instead of running the whole script, you send each script line through the REPL.
- Benefit: Interactive environment to speed up test creation process. It also helps while learning the k6 APIs.

## Using REPLIT inside an existing k6 test

You can drop into a REPL from within an existing k6 test via the following methods.

1. Import from `k6/x/replit`.
1. Use the `replit.with` function to block the script execution.
1. Optionally pass an object containing the context you want to interact with (including variables, modules, etc.).
1. Run the script as explained in the [running](#running-replit) section.
1. Interact with the REPL as explained in the [interacting](#interacting-with-the-repl) section.

Here's an example:

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

> [!TIP]
> You can add as many `replit.with` calls as you want in your script, and you can pass different contexts to each one.

### Interacting with the REPL

Once you [run](#running) the above script, you will be dropped into a REPL where you can interact with the context you passed in.

```bash
>>> result
{ status: 200, statusText: 'OK', headers: {...}, body: ... }
>>> result.status
200
>>> http
{ get: [Function: get], post: [Function: post], ... }
>>> other = http.get("https://quickpizza.grafana.com")
>>> other
{ status: 200, statusText: 'OK', headers: {...}, body: ... }
```

### Multi-line input

REPLIT supports basic multi-line input. Once it sees a semicolon, it will execute the code. Otherwise, it will wait for more input. This allows you to copy and paste code snippets into the REPL directly and execute them.

```bash
>>> await new Promise(
...     resolve
... ) => setTimeout(resolve, 1000))
```

> [!TIP]
> You can use `await` to wait for promises to resolve.

### Auto-completion

REPLIT also supports auto-completion, so you can press `Tab` to see available properties and methods.

```bash
>>> result <- pressed the tab key
result                        result.timings                result.blocked
>>> result.timings
{ ... }
```

> [!TIP]
> Press CTRL+D (or CMD+D on macOS) to exit the REPL and continue running the script.

## Running REPLIT

You can run a standalone REPL (i.e. without any context) via the following methods:

### Run using CLI helper

```bash
go run ./cmd/replit examples/http.js
```

### Build the CLI tool and run it separately

This will allow you to re-run the CLI tool much faster.

```bash
go build -o bin/replit ./cmd/replit
./bin/replit examples/http.js
```

> [!NOTE]
> The CLI tool requires [Go](https://go.dev/doc/install) and [xk6](https://github.com/grafana/xk6) to be installed.


### Run as an extension

```bash
xk6 run -q --no-summary examples/http.js
```

### Build the extension and run the k6 binary separately

This will allow you to re-run the k6 binary much faster.

```bash
xk6 build --with github.com/grafana/xk6-replit=.
./k6 run -q --no-summary examples/http.js
```
