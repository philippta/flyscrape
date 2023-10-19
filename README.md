<br />

<p align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/logo-alt.png">
  <source media="(prefers-color-scheme: light)" srcset="docs/logo.png">
  <img width="200" src="docs/logo.png">
</picture>

</p>

<br />

<p align="center">
<b>flyscrape</b> is a standalone and scriptable web scraper, combining the speed of Go with the flexibility of JavaScript. — Focus on data extraction rather than request juggling.
</p>

<br />

<p align="center">
<a href="#installation">Installation</a> · <a href="/docs">Documentation</a> · <a href="https://github.com/philippta/flyscrape/releases">Releases</a>
</p>

## Features

- **Highly Configurable:** 10 options to fine-tune your scraper.
- **Standalone:** flyscrape comes as a single binary executable.
- **Scriptable:** Use JavaScript to write your data extraction logic.
- **Simple API:** Extract data from HTML pages with a familiar API.
- **Fast Iteration:** Use the development mode to get quick feedback.
- **Request Caching:** Re-run scripts on websites you already scraped.
- **Zero Dependencies:** No need to fill up your disk with npm packages.

## Example script

```javascript
export const config = {
    url: "https://news.ycombinator.com/",
}

export default function ({ doc, absoluteURL }) {
    const title = doc.find("title");
    const posts = doc.find(".athing");

    return {
        title: title.text(),
        posts: posts.map((post) => {
            const link = post.find(".titleline > a");

            return {
                title: link.text(),
                url: link.attr("href"),
            };
        }),
    }
}
```

```bash
$ flyscrape run hackernews.js
[
  {
    "url": "https://news.ycombinator.com/",
    "data": {
      "title": "Hacker News",
      "posts": [
        {
          "title": "Show HN: flyscrape - An standalone and scriptable web scraper",
          "url": "https://flyscrape.com/"
        },
        ...
      ]
    }
  }
]
```

## Installation

### Pre-compiled binary

`flyscrape` is available for MacOS, Linux and Windows as a downloadable binary from the [releases page](https://github.com/philippta/flyscrape/releases).

### Compile from source

To compile flyscrape from source, follow these steps:

1. Install Go: Make sure you have Go installed on your system. If not, you can download it from [https://golang.org/](https://golang.org/).

2. Install flyscrape: Open a terminal and run the following command:

   ```bash
   go install github.com/philippta/flyscrape/cmd/flyscrape@latest
   ```

## Usage

```
flyscrape is a standalone and scriptable web scraper for efficiently extracting data from websites.

Usage:

    flyscrape <command> [arguments]

Commands:

    new    creates a sample scraping script
    run    runs a scraping script
    dev    watches and re-runs a scraping script
```

## Configuration

Below is an example scraping script that showcases the capabilities of flyscrape. For a full documentation of all configuration options, visit the [documentation page](docs/readme.md#configuration).

```javascript
export const config = {
    url: "https://example.com/", // Specify the URL to start scraping from.
    depth: 0,                    // Specify how deep links should be followed.  (default = 0, no follow)
    follow: [],                  // Speficy the css selectors to follow         (default = ["a[href]"])
    allowedDomains: [],          // Specify the allowed domains. ['*'] for all. (default = domain from url)
    blockedDomains: [],          // Specify the blocked domains.                (default = none)
    allowedURLs: [],             // Specify the allowed URLs as regex.          (default = all allowed)
    blockedURLs: [],             // Specify the blocked URLs as regex.          (default = none)
    rate: 100,                   // Specify the rate in requests per second.    (default = no rate limit)
    proxies: [],                 // Specify the HTTP(S) proxy URLs.             (default = no proxy)
    cache: "file",               // Enable file-based request caching.          (default = no cache)
};

export default function ({ doc, url, absoluteURL }) {
    // doc              - Contains the parsed HTML document
    // url              - Contains the scraped URL
    // absoluteURL(...) - Transforms relative URLs into absolute URLs
}
```

## Query API

```javascript
// <div class="element" foo="bar">Hey</div>
const el = doc.find(".element")
el.text()                                 // "Hey"
el.html()                                 // `<div class="element">Hey</div>`
el.attr("foo")                            // "bar"
el.hasAttr("foo")                         // true
el.hasClass("element")                    // true

// <ul>
//   <li class="a">Item 1</li>
//   <li>Item 2</li>
//   <li>Item 3</li>
// </ul>
const list = doc.find("ul")
list.children()                           // [<li class="a">Item 1</li>, <li>Item 2</li>, <li>Item 3</li>]

const items = list.find("li")
items.length()                            // 3
items.first()                             // <li>Item 1</li>
items.last()                              // <li>Item 3</li>
items.get(1)                              // <li>Item 2</li>
items.get(1).prev()                       // <li>Item 1</li>
items.get(1).next()                       // <li>Item 3</li>
items.get(1).parent()                     // <ul>...</ul>
items.get(1).siblings()                   // [<li class="a">Item 1</li>, <li>Item 2</li>, <li>Item 3</li>]
items.map(item => item.text())            // ["Item 1", "Item 2", "Item 3"]
items.filter(item => item.hasClass("a"))  // [<li class="a">Item 1</li>]
```

## Issues and Suggestions

If you encounter any issues or have suggestions for improvement, please [submit an issue](https://github.com/philippta/flyscrape/issues).
