---
title: 'Cookies'
weight: 9
---

The Cookies configuration in the `flyscrape` script's configuration object allows you to specify the behavior of the cookie store during the scraping process. Cookies are often used for authentication and session management on websites.

## Cookies Configuration

To configure the cookie store behavior, set the `cookies` field in your configuration. The `cookies` option supports three values: `"chrome"`, `"edge"`, and `"firefox"`. Each value corresponds to using the cookie store of the respective local browser.

When the `cookies` option is set to `"chrome"`, `"edge"`, or `"firefox"`, `flyscrape` utilizes the cookie store of the user's installed browser.

```javascript {filename="Configuration"}
export const config = {
    cookies: "chrome",
};
```

In the above example, the `cookies` option is set to `"chrome"`, indicating that `flyscrape` should use the cookie store of the local Chrome browser.

```javascript {filename="Configuration"}
export const config = {
    cookies: "firefox",
};
```

In this example, the `cookies` option is set to `"firefox"`, instructing `flyscrape` to use the cookie store of the local Firefox browser.

```javascript {filename="Configuration"}
export const config = {
    cookies: "edge",
};
```

In this example, the `cookies` option is set to `"edge"`, indicating that `flyscrape` should use the cookie store of the local Edge browser.
