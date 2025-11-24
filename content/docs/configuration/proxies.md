---
title: 'Proxies'
weight: 8
---

The proxy feature allows you to route your scraping requests through a specified HTTP(S) proxy. This can be useful for bypassing IP-based rate limits or accessing region-restricted content.

```javascript
export const config = {
    // Specify a single HTTP(S) proxy URL.
    proxy: "http://someproxy.com:8043",
};
```

In the above example, all scraping requests will be routed through the proxy at `http://someproxy.com:8043`.

## Multiple Proxies

You can also specify multiple proxy URLs. The scraper will rotate between these proxies for each request.

```javascript
export const config = {
    // Specify multiple HTTP(S) proxy URLs.
    proxies: [
      "http://someproxy.com:8043",
      "http://someotherproxy.com:8043",
    ],                     
};
```

In this example, the scraper will randomly pick between the proxies at `http://someproxy.com:8043` and `http://someotherproxy.com:8043`.

Note: If both `proxy` and `proxies` are specified, all proxies will be respected.
