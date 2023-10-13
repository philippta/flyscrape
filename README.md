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
<b>flyscrape</b> is an expressive and elegant web scraper, combining the speed of Go with the <br/> flexibility of JavaScript. — Focus on data extraction rather than request juggling.
</p>

<br />

## Features

- Domains and URL filtering
- Depth control
- Request caching
- Rate limiting
- Development mode
- Single binary executable


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
          "title": "Show HN: flyscrape - An expressive and elegant web scraper",
          "url": "https://flyscrape.com"
        },
        ...
      ]
    }
  }
]
```

## Installation

To install **flyscrape**, follow these simple steps:

1. Install Go: Make sure you have Go installed on your system. If not, you can download it from [https://golang.org/](https://golang.org/).

2. Install **flyscrape**: Open a terminal and run the following command:

   ```bash
   go install github.com/philippta/flyscrape/cmd/flyscrape@latest
   ```

## Usage

```
$ flyscrape
flyscrape is an elegant scraping tool for efficiently extracting data from websites.

Usage:

    flyscrape <command> [arguments]

Commands:

    new    creates a sample scraping script
    run    runs a scraping script
    dev    watches and re-runs a scraping script

```

### Create a new sample scraping script

The `new` command allows you to create a new boilerplate sample script which helps you getting started.

```
flyscrape new example.js
```

### Watch the script for changes during development

The `dev` command allows you to watch your scraping script for changes and quickly iterate during development. In development mode, flyscrape will not follow any links and request caching is enabled.

```
flyscrape dev example.js
```

### Run the scraping script

The `run` command allows you to run your script.

```
flyscrape run example.js
```

## Configuration

Below is an example scraping script that showcases the capabilities of **flyscrape**:

```javascript
export const config = {
    url: "https://example.com/", // Specify the URL to start scraping from.
    depth: 0,                    // Specify how deep links should be followed.  (default = 0, no follow)
    allowedDomains: [],          // Specify the allowed domains. ['*'] for all. (default = domain from url)
    blockedDomains: [],          // Specify the blocked domains.                (default = none)
    allowedURLs: [],             // Specify the allowed URLs as regex.          (default = all allowed)
    blockedURLs: [],             // Specify the blocked URLs as regex.          (default = none)
    rate: 100,                   // Specify the rate in requests per second.    (default = no rate limit)
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

## Contributing

We welcome contributions from the community! If you encounter any issues or have suggestions for improvement, please [submit an issue](https://github.com/philippta/flyscrape/issues).
