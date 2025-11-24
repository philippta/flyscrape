---
title: 'Depth'
weight: 2
---

The `depth` config option allows you to specify how deep the scraping process should follow links from the initial URL.

When no value is provided or `depth` is set to `0` link following is disabled and it will only scrape the initial URL.

```javascript {filename="Configuration"}
export const config = {
  url: "http://example.com/",
  depth: 2,
  // ...
};
```

With the config provided in the example the scraper would follow links like this:

```
http://example.com/                    (depth = 0, initial URL)
↳ http://example.com/deeply            (depth = 1)
  ↳ http://example.com/deeply/nested   (depth = 2)
```
