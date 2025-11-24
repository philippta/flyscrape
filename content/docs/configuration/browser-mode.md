---
title: 'Browser Mode'
weight: 10
---

The Browser Mode controls the interaction with a headless Chromium browser. Enabling the browser mode allows `flyscrape` to download a Chromium browser once and use it to render JavaScript-heavy pages.

## Browser Mode

To enable Browser Mode, set the `browser` option to `true` in your configuration. This allows `flyscrape` to use a headless Chromium browser for rendering JavaScript during the scraping process.

```javascript {filename="Configuration"}
export const config = {
    browser: true,
};
```

In the above example, Browser Mode is enabled, allowing `flyscrape` to render pages that rely on JavaScript execution.

## Headless Option

The `headless` option, when combined with Browser Mode, controls whether the Chromium browser should run in headless mode or not. Headless mode means the browser operates without a graphical user interface, which can be useful for background processes.

```javascript {filename="Configuration"}
export const config = {
    browser: true,
    headless: false,
};
```

In this example, the Chromium browser will run in non-headless mode. If you set `headless` to `true`, the browser will run without a visible GUI.

```javascript {filename="Configuration"}
export const config = {
    browser: true,
    headless: true,
};
```

In this example, the Chromium browser will run in headless mode, suitable for scenarios where graphical rendering is unnecessary.
