# Stating URL

The `url` config option allows you to specify the initial URL at which the scraper should start its scraping process.

When no value is provided, the scraper will not start and exit immediately.

Example:

```javascript
export const config = {
    url: "http://example.com/",
    // ...
};
```
