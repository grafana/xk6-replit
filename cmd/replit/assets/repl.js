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
    await replit.repl(replit, {http: http, browser: browser, sleep: sleep});
}

