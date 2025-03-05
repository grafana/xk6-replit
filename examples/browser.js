import { replit } from "k6/x/replit";
import { browser } from "k6/browser";

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
    await replit.with({browser: browser});
}