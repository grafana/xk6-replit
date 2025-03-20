
<img width="822" alt="Screenshot 2025-03-07 at 2 54 54 PM" src="https://github.com/user-attachments/assets/9f0ee123-e369-440a-ba25-17af102229d4" />

# REPLIT

A [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) (Read-Evaluate-Print Loop) for Grafana [k6](https://github.com/grafana/k6)!

REPLIT provides an interactive shell where you can explore and learn k6 APIs quickly, write new tests iteratively, interact with your web browser, and troubleshoot running tests in real time—all with just one function call.

It makes test authoring for reliability testing more enjoyable and effective.

> [!WARNING]  
> This extension is experimental and not officially supported by Grafana k6. As such, no official support will be available. Please use it with that understanding, and feel free to share any feedback or suggestions!

## Short demo

This short demo shows interacting with the official k6/http module in real time. It talks to a running k6 program and executes commands from the command line.

![demo](https://github.com/user-attachments/assets/9b3ba9c9-97e2-4dce-b40e-e2bd56dd51db)

## A typical workflow of developing a k6 script today can be time-consuming

- Often editing scripts and rerunning them all over until we make them work.
- Wait for the results. Repeat.
- 🥱 This cycle can be time-consuming, disrupt focus, and be unproductive.

## So, we thought, why doesn’t k6 have a REPL?

- A REPL shortens the loop from coding to seeing results in real time.
- Saving time and increasing productivity, creativity, and efficiency.
- Allows us to come up with test scripts much quicker.
- Allows us to experiment and learn k6 APIs more conveniently.

## Features

- We built REPLIT as a k6 extension.
- It also runs as a standalone cross-platform executable without any dependencies.
- REPLIT supports bash-like navigation keys (e.g., up, down, Ctrl-E).
- Automatically stores previously entered input in history.
- Syntax-highlighting different forms of output.
- Runs modern async/await expressions with no problems.
- It can be dropped into any existing script with a single replit.run() call.

### Explore any k6 module APIs

You can explore any k6 module just by including it in the context object passed to `replit.run`. See the [examples](./examples) folder for module exploration examples.

### Colorized output

REPLIT colorizes the output, making it easier to read and understand. It outputs differently depending on the output type (e.g., objects, strings, etc.). For objects, it uses JSON syntax highlighting, and for others, it uses JavaScript syntax highlighting.

### Multi-line input

REPLIT supports multi-line input. It waits to execute the code until your input is compilable. This allows you to copy and paste code snippets directly into the REPL.

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

----

## Install REPLIT

There's two ways of getting REPLIT:

### Via prebuilt `replit` binaries

You can find REPLIT binaries in the [releases](https://github.com/grafana/xk6-replit/releases) section. These binaries simply wrap an extended `k6` binary with some additional utilities and pre-set CLI arguments.

### Via xk6

Clone this repository and build the extended `k6` binary manually, using [xk6](https://github.com/grafana/xk6):
```bash
git clone git@github.com:grafana/xk6-replit.git
cd xk6-replit
xk6 build --with github.com/grafana/xk6-replit=.
./k6 version
```

---

## Use REPLIT

You can use REPLIT in two different ways: either from within an existing k6 script, or as a standalone application.

### REPLIT in your k6 script

You can drop into a REPL from within an existing k6 test by doing the following:
1. Import from `k6/x/replit`.
2. Use the `replit.run` function to block the script execution.
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
    await replit.run({result: result, http: http});

    console.log("All done.")
}
```

> [!NOTE]
> Your exported default function will need to be `async` for `replit` to work.

> [!TIP]
> You can add as many `replit.run` calls as you want in your script, and you can pass different contexts to each one.

To run the script, use either:

* `./replit my_script.js` (in case you downloaded the REPLIT binary) or:
* `./k6 run -q --no-summary my_script.js` (if you built `k6` locally using xk6).

You will then be able to do the following:

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
