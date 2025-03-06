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

        // Let, const, and var are not allowed in the REPL context.
        // However, we automatically remove them when they are used,
        // and warn the user once about it.
        let illegalVarAssignment = {pattern: /^(let|const|var)\s+/, matched: false, warned: false};

        while (true) {
            try {
                var input = replit.read().trim();
                if (input === "exit") {
                    break;
                }

                illegalVarAssignment.matched = false;
                if (input.match(illegalVarAssignment.pattern)) {
                    input = input.replace(illegalVarAssignment.pattern, "");
                    illegalVarAssignment.matched = true;
                }

                // FIXME: See comment in replit.go readMultiLine.

                var fn = undefined;
                try {
                    fn = AsyncFunction(...params.concat(["return " + input])); // Input was an expression
                } catch (error) {
                    fn = AsyncFunction(...params.concat([input]));             // Input was a statement
                }

                var result = await fn(...args); // the user's code result
                global._ = result;              // Easily access the last expression result with '_'.

                // Do a quick check to see if it's an object, if so, pretty print it
                // Otherwise, fall back to stringifying it if it's not undefined.
                if (typeof result === "object") {
                    replit.highlight(_inspect(result), "json");
                } else if (result !== undefined) {
                    replit.highlight(result.toString(), "javascript");
                } else {
                    replit.log("undefined")
                }

                if (illegalVarAssignment.matched && !illegalVarAssignment.warned) {
                    // A long output can easily push the warning out of the screen.
                    // So, we warn the user about illegal variable assignment here.
                    replit.warn("Variable assignment with `let`, `const`, or `var` has no effect in the REPL context.");
                    replit.warn("For convenience, REPLIT automatically removes `let`, `const`, or `var` from assignments.");
                    replit.log("Hint: in order to assign a variable globally, use `foo = 123`.");
                    illegalVarAssignment.warned = true;
                }
            } catch (error) {
                if (error.toString() == "GoError: EOF") {
                    break;
                }
                replit.error(error.toString());
            }
        }
    };

    return {repl: repl};
}());

