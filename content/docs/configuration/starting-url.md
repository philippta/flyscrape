---
title: 'Starting URL'
weight: 1
prev: '/docs/configuration'
---

The `url` config option allows you to specify the initial URL at which the scraper should start its scraping process.

```javascript {filename="Configuration"}
export const config = {
  url: "http://example.com/",
  // ...
};
```

## Multiple starting URLs

In case you have more than one URL you want to scrape (or to start from) you can specify them with the `urls` config option.

```javascript {filename="Configuration"}
export const config = {
  urls: [
    "http://example.com/",
    "http://anothersite.com/",
    "http://yetanothersite.com/",
  ],
  // ...
};
```
