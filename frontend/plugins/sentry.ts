import * as Sentry from "@sentry/vue";

export default defineNuxtPlugin(nuxt => {
    const sentryDSN = nuxt.$config.public["sentryDSN"];

    if (sentryDSN !== undefined && typeof sentryDSN !== "string") {
        return;
    }

    nuxt.hook("app:created", () => {
        Sentry.init({
            app: nuxt.vueApp,
            dsn: sentryDSN,
            integrations: [
                new Sentry.Integrations.Dedupe(),
                new Sentry.BrowserTracing(),
                new Sentry.Replay(),
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
            await Sentry.flush(5000);
        },
    )

    return {
        provide: {
            sentry: Sentry,
        }
    }
})