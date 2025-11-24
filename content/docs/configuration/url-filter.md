---
title: 'URL Filter'
weight: 4
prev: /docs/getting-started
---

The `allowedURLs` and `blockedURLs` config options allow you to specify a list of URL patterns (in form of regular expressions) which are accessible or blocked during scraping.

```javascript {filename="Configuration"}
export const options = {
  url: "http://example.com/",
  allowedURLs: ["/articles/.*", "/authors/.*"],
  blockedURLs: ["/authors/admin"],
  // ...
};
```

## Allowed URLs

This config option controls which URLs are allowed to be visted during scraping. When no value is provided all URLs are allowed to be visited if not otherwise blocked.

When a list of URL patterns is provided, only URLs matching one or more of these patterns are allowed to be visted.

```javascript {filename="Configuration"}
export const options = {
  url: "http://example.com/",
  allowedURLs: ["/products/"],
};
```

## Blocked URLs

This config option controls which URLs are blocked from being visted during scraping.

When a list of URL patterns is provided, URLs matching one or more of these patterns are blocked from to be visted.

```javascript {filename="Configuration"}
export const options = {
  url: "http://example.com/",
  blockedURLs: ["/restricted"],
};
```
