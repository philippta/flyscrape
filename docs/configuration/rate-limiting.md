# Rate Limiting

The `rate` config option allows you to specify at which rate the scraper should send out requests. The rate is measured in _Requests per Second_ (RPS) and can be set as a whole or decimal number to account for shorter and longer request intervals.

When no `rate` is specified, rate limiting is disabled and the scraper will send out requests as fast as it can.

Example:

```javascript
export const options = {
    url: "http://example.com/",
    rate: 50,
};
```
