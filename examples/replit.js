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

            var result = await fn();

            // Easily access the last expression result with '_'.
            global._ = result;

            if (result !== undefined && result !== null) {
                console.log(result.toString());
            }
        } catch (error) {
            if (error.toString() == "GoError: EOF") {
                break;
            }
            console.error(error.toString());
        }

        input = input.trim();
        if (input.startsWith("let ") || input.startsWith("const ") || input.startsWith("var ")) {
            console.error("Invalid assignment in REPL context.");
            console.info("Hint: In order to set a variable globally, use `foo = 123`.");
        }
    }
}
