# Caching

The `cache` config option allows you to enable file-based request caching. When enabled every request cached with its raw response. When the cache is populated and you re-run the scraper, requests will be served directly from cache.

This also allows you to modify your scraping script afterwards and collect new results immediately.

Example:

```javascript
export const config = {
    url: "http://example.com/",
    cache: "file",
    // ...
};
```

### Cache File

When caching is enabled using the `cache: "file"` option, a `.cache` file will be created with the name of your scraping script.

Example:

```bash
$ flyscrape run hackernews.js # Will populate: hackernews.cache
```

### Shared cache

In case you want to share a cache between different scraping scripts, you can specify where to store the cache file.

```javascript
export const config = {
    url: "http://example.com/",
    cache: "file:/some/path/shared.cache",
    // ...
};
```
