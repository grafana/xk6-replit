import { replit } from "k6/x/replit";

export default async function () {
    console.log(replit.greeting);
    replit.run("var hello = 'Hello, world!'");
    replit.run("console.log(hello)");
}