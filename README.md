# REPLIT

Your [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) to [k6](https://github.com/grafana/k6).

REPLIT allows you to interact with your k6 script in real-time, inspecting variables, and running code snippets. This approach can be useful for debugging, testing, and learning k6 APIs.

- You write and send commands to k6 to execute and receive execution results.
- Instead of running the whole script, you send each script line through the REPL.
- Benefit: Interactive environment to speed up test creation process. It also helps while learning the k6 APIs.

## Install REPLIT

There's two ways of getting REPLIT:

### Via prebuilt `replit` binaries

You can find REPLIT binaries in the [releases](https://github.com/grafana/xk6-replit/releases/tag/v0.1) section. These binaries simply wrap an extended `k6` binary with some additional utilities and pre-set CLI arguments.

### Via xk6

Clone this repository and build the extended `k6` binary manually, using [xk6](https://github.com/grafana/xk6):
```bash
git clone git@github.com:grafana/xk6-replit.git
cd xk6-replit
xk6 build --with github.com/grafana/xk6-replit=.
./k6 version
```

## Use REPLIT

You can use REPLIT in two different ways: either from within an existing k6 script, or as a standalone application.

### REPLIT in your k6 script

You can drop into a REPL from within an existing k6 test by doing the following:
1. Import from `k6/x/replit`.
2. Use the `replit.with` function to block the script execution.
3. Optionally pass an object containing the context you want to interact with (including variables, modules, etc.).
4. Run REPLIT.

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

To run the script, use either `./replit my_script.js` (in case you downloaded the REPLIT binary) or `./k6 run my_script.js` (if you built `k6` locally using xk6). You will then be able to do the following:

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

### REPLIT Standalone

In order to run REPLIT as as a standalone application (for example, to experiment with different k6 JS libraries), download one of the REPLIT binaries and run it without any arguments (`./replit`).

## Features
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

REPLIT also supports auto-completion, so you can press `Tab` to see the relevant commands you previously typed in (from the history).

```bash
>>> result <- pressed the tab key
result                        result.timings                result.blocked
>>> result.timings
{ ... }
```

> [!TIP]
> Press CTRL+D (or CMD+D on macOS) to exit the REPL and continue running the script.
