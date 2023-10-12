# flyscrape

flyscrape is an elegant scraping tool for efficiently extracting data from websites. Whether you're a developer, data analyst, or researcher, flyscrape empowers you to effortlessly gather information from web pages and transform it into structured data. 

## Features

- **Simple and Intuitive**: **flyscrape** offers an easy-to-use command-line interface that allows you to interact with scraping scripts effortlessly.

- **Create New Scripts**: The `new` command enables you to generate sample scraping scripts quickly, providing you with a solid starting point for your scraping endeavors.

- **Run Scripts**: Execute your scraping script using the `run` command, and watch as **flyscrape** retrieves and processes data from the specified website.

- **Watch for Development**: The `dev` command allows you to watch your scraping script for changes and quickly iterate during development, helping you find the right data extraction queries.

## Installation

To install **flyscrape**, follow these simple steps:

1. Install Go: Make sure you have Go installed on your system. If not, you can download it from [https://golang.org/](https://golang.org/).

2. Install **flyscrape**: Open a terminal and run the following command:

   ```bash
   go install github.com/philippta/flyscrape/cmd/flyscrape@latest
   ```

## Usage

**flyscrape** offers several commands to assist you in your scraping journey:

### Creating a New Script

Use the `new` command to create a new scraping script:

```bash
flyscrape new example.js
```

### Running a Script

Execute your scraping script using the `run` command:

```bash
flyscrape run example.js
```

### Watching for Development

The `dev` command allows you to watch your scraping script for changes and quickly iterate during development:

```bash
flyscrape dev example.js
```

## Example Script

Below is an example scraping script that showcases the capabilities of **flyscrape**:

```javascript
export const config = {
  url: "https://news.ycombinator.com/", // Specify the URL to start scraping from.
  // depth: 0,                          // Specify how deep links should be followed.  (default = 0, no follow)
  // allowedDomains: [],                // Specify the allowed domains. ['*'] for all. (default = domain from url)
  // blockedDomains: [],                // Specify the blocked domains.                (default = none)
  // allowedURLs: [],                   // Specify the allowed URLs as regex.          (default = all allowed)
  // blockedURLs: [],                   // Specify the blocked URLs as regex.          (default = non blocked)
  // rate: 100,                         // Specify the rate in requests per second.    (default = 100)
  // cache: "file",                     // Enable file-based request caching.          (default = no cache)
};

export default function({ doc, absoluteURL }) {
  const title = doc.find("title");
  const posts = doc.find(".athing");

  return {
    title: title.text(),
    posts: posts.map((post) => {
      const link = post.find(".titleline > a");

      return {
        title: link.text(),
        url: absoluteURL(link.attr("href")),
      };
    }),
  };
}
```

## Contributing

We welcome contributions from the community! If you encounter any issues or have suggestions for improvement, please [submit an issue](https://github.com/philippta/flyscrape/issues).
