// replit.js
// Contains JavaScript functions for the replit module.
// Other functions are defined in Go as well.
// In particular, Immediately Invoked Function Expressions (IIFE) are used here.

(function() {
    const AsyncFunction = Object.getPrototypeOf(async function(){}).constructor;

    function _inspect(obj) {
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

    async function repl(replit, context) {
        // Set some global properties, from context, for the user to use.
        for (const [k, v] of Object.entries(context || {})) {
            global[k] = v
        }

        while (true) {
            try {
                var input = replit.read();
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
                    replit.log(typeof result === "object" ? _inspect(result) : result.toString());
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
    };

    return {repl: repl};
}());

