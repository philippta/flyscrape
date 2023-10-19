# Proxies

The `proxies` config option allows you to specify a list of HTTP(S) proxies that should used during scraping. When multiple proxies are provided, the scraper will prick a proxy at random for each request.

Example:

```javascript
export const config = {
    url: "http://example.com/",
    proxies: ["https://my-proxy.com:3128", "https://my-other-proxy.com:8080"],
    // ...
};
```
