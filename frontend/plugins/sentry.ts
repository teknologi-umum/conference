import { init } from "@sentry/vue";
import { defaultIntegrations, BrowserTracing, Replay, flush } from "@sentry/browser"
import * as Sentry from "@sentry/browser";

export default defineNuxtPlugin(nuxt => {
    const sentryDSN = nuxt.$config.public["sentryDSN"];

    if (sentryDSN !== undefined && typeof sentryDSN !== "string") {
        return;
    }

    nuxt.hook("app:created", () => {
        init({
            app: nuxt.vueApp,
            dsn: sentryDSN,
            integrations: [
                ...defaultIntegrations,
                new BrowserTracing(),
                new Replay(),
            ],
            sampleRate: 1.0,
            tracesSampleRate: 0.2,
            tracePropagationTargets: ["https://conf.teknologiumum.com"],
            replaysOnErrorSampleRate: 1.0,
            replaysSessionSampleRate: 0.01,
        });
    });

    nuxt.hook(
        "page:transition:finish",
        async () => {
            await flush(5000);
        },
    )

    return {
        provide: {
            sentry: Sentry,
        }
    }
})