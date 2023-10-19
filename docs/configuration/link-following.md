# Link Following

The `follow` config option allows you to specify a list of CSS selectors that determine which links the scraper should follow.

When no value is provided the scraper will follow all links found with the `a[href]` selector.

Example:

```javascript
export const config = {
    url: "http://example.com/",
    follow: [".pagination > a[href]", ".nav a[href]"],
    // ...
};
```

### Following non `href` attributes

For special cases where the link is not to be found in the `href`, you specify a selector with a different ending attribute.

Example:

```javascript
export const config = {
    url: "http://example.com/",
    follow: [".articles > div[data-url]"],
    // ...
};
```
