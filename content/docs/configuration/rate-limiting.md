---
title: 'Rate Limiting'
weight: 6
---

The `rate` config option allows you to specify at which rate the scraper should send out requests. The rate is measured in _Requests per Minute_ (RPM).

When no `rate` is specified, rate limiting is disabled and the scraper will send out requests as fast as it can.

```javascript {filename="Configuration"}
export const options = {
  url: "http://example.com/",
  rate: 100,
};
```
