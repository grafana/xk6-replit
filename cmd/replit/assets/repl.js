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

// TODO:
// We want to import more modules by default (maybe even all of them?)
// or allow the user to somehow specify this via the CLI helper.

export default async function () {
    await replit.repl(replit, {http: http, browser: browser, sleep: sleep});
}

