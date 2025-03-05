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
    // We need to copy these over to the global context otherwise they
    // are not accessible from the newly created AsyncFunction below.
    global.http = http
    global.browser = browser
    global.sleep = sleep

    await replit.repl(replit);
}

