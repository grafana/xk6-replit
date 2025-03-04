import http from "k6/http";
import { browser } from "k6/browser";
import { replit } from "k6/x/replit";
import { sleep } from "k6";

export const options = {
    scenarios: {
        ui: {
            executor: "shared-iterations",
            options: {
                browser: {
                    type: "chromium",
                },
            },
        },
    },
};

export default async function () {
    const AsyncFunction = Object.getPrototypeOf(async function(){}).constructor;

    // console.log("Required modules by this script: ", required_modules);

    // We need to copy these over to the global context otherwise they
    // are not accessible from the newly created AsyncFunction below.
    global.http = http
    global.browser = browser
    global.sleep = sleep

    while (true) {
        try {
            var input = replit.read("> ");
            if (input === "exit") {
                break;
            }

            var fn = undefined;
            try {
                // Input was an expression
                fn = AsyncFunction("return " + input)
            } catch (error) {
                // Input was a statement
                fn = AsyncFunction(input)
            }

            var result = await fn(); // the user's code result
            if (result !== undefined && result !== null) {
                // Fall back to .toString() if it's a primitive or can't be stringified
                // Or do a quick test to see if it's an object
                replit.log(typeof result === "object" ? inspect(result) : result.toString());
            }

            // Easily access the last expression result with '_'.
            global._ = result;
        } catch (error) {
            if (error.toString() == "GoError: EOF") {
                break;
            }
            replit.error(error.toString(), 'red');
        }

        input = input.trim();
        if (input.startsWith("let ") || input.startsWith("const ") || input.startsWith("var ")) {
            replit.error("Invalid assignment in REPL context.", 'red');
            replit.error("Hint: In order to set a variable globally, use `foo = 123`.", 'red');
        }
    }
}

function inspect(obj) {
  // NOTE: Circular references will throw an error unless we handle them,
  // so let's do a naive replacer that short-circuits circular refs.
  const seen = new WeakSet();
  return JSON.stringify(
    obj,
    (key, value) => {
      if (typeof value === "object" && value !== null) {
        if (seen.has(value)) {
          return "[Circular]";
        }
        seen.add(value);
      }
      return value;
    },
    2 // indentation
  );
}

