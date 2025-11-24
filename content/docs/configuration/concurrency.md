---
title: 'Concurrency'
weight: 6
---

The concurrency setting controls the number of simultaneous requests that the scraper can make. This is specified in the configuration object of your scraping script.

```javascript
export const config = {
    // Specify the number of concurrent requests.
    concurrency: 5,
};
```

In the above example, the scraper will make up to 5 requests at the same time.

If the concurrency setting is not specified, there is no limit to the number of concurrent requests.

