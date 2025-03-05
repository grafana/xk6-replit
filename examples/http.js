import { replit } from "k6/x/replit";
import http from "k6/http";

export default async function () {
    let result = http.get("https://quickpizza.grafana.com");

    // As context, we pass 'result', and the http module in case
    // we want to make additional requests.
    await replit.with({result: result, http: http});

    console.log("All done.")
}