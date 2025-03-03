import { replit } from "k6/x/replit";

export default function () {
    console.log(replit.greeting);
    replit.block(); // Execution stops here until 'exit' is sent in the REPL
    console.log("Script continues after REPL exit");
}
