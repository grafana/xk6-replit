// replit.js
// Contains JavaScript functions for the replit module.
// Other functions are defined in Go as well.
// In particular, Immediately Invoked Function Expressions (IIFE) are used here.

var replit; // replit module will be injected by the module itself

(function() {
    const AsyncFunction = Object.getPrototypeOf(async function(){}).constructor;

    function _inspect(obj) {
        // NOTE: Circular references will throw an error unless we handle them,
        // so let's do a naive replacer that short-circuits circular refs.
        // FIXME: JSON has some limitations in representing objects, a bespoke solution
        // would be better in order to properly print undefined, Symbols, Functions, etc.
        const seen = new WeakSet();
        return JSON.stringify(
            obj,
            (key, value) => {
                if (typeof value === "object" && value !== null) {
                    if (seen.has(value)) {
                        return "[Circular]";
                    }
                    seen.add(value);
                } else if (typeof value === "function") {
                    return "[Function]";
                }
                return value;
            },
            2 // indentation
        );
    }

    async function repl(context) {
        let params = [];
        let args = [];
        for (const [k, v] of Object.entries(context || {})) {
            params.push(k);
            args.push(v);
        }

        while (true) {
            try {
                var input = replit.read();
                if (input === "exit") {
                    break;
                }

                // FIXME: See comment in replit.go readMultiLine.

                var fn = undefined;
                try {
                    // Input was an expression
                    fn = AsyncFunction(...params.concat(["return " + input]));
                } catch (error) {
                    // Input was a statement
                    fn = AsyncFunction(...params.concat([input]));
                }

                var result = await fn(...args); // the user's code result

                if (result !== undefined) {
                    // Fall back to .toString() if it's a primitive or can't be stringified
                    // Or do a quick test to see if it's an object
                    if (typeof result === "object") {
                        replit.highlight(_inspect(result), "json");
                    } else {
                        replit.highlight(result.toString(), "javascript");
                    }
                } else {
                    replit.log("undefined")
                }

                // Easily access the last expression result with '_'.
                global._ = result;
            } catch (error) {
                if (error.toString() == "GoError: EOF") {
                    break;
                }
                replit.error(error.toString());
            }

            input = input.trim();
            if (input.startsWith("let ") || input.startsWith("const ") || input.startsWith("var ")) {
                replit.warn("Variable assignment had no effect in REPL context.");
                replit.log("Hint: in order to assign a variable globally, use `foo = 123`.");
            }
        }
    };

    return {repl: repl};
}());

