# Domain Filter

The `allowedDomains` and `blockedDomains` config options allow you to specify a list of domains which are accessible or blocked during scraping.

```javascript
export const options = {
    url: "http://example.com/",
    allowedDomains: ["subdomain.example.com"],
    // ...
};
```

### `allowedDomains`

This config option controls which additional domains are allowed to be visted during scraping. The domain of the initial URL is always allowed.

You can also allow all domains to be accessible by setting `allowedDomains` to `["*"]`. To then further restrict access, you can specify `blockedDomains`.

Example:

```javascript
export const options = {
    url: "http://example.com/",
    allowedDomains: ["*"],
    // ...
};
```

### `blockedDomains`

This config option controls which additional domains are blocked from being accessed. By default all domains other than the domain of the initial URL or those specified in `allowedDomains` are blocked.

You can best use `blockedDomains` in conjunction with `allowedDomains: ["*"]`, allowing the scraping process to access all domains except what's specified in `blockedDomains`.

Example:

```javascript
export const options = {
    url: "http://example.com/",
    allowedDomains: ["*"],
    blockedDomains: ["google.com", "bing.com"],
    // ...
};
```
