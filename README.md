# flyscrape - Elegant Website Scraping Tool

flyscrape is a powerful command-line tool designed to streamline the process of efficiently extracting data from websites. Whether you're a developer, data analyst, or researcher, flyscrape empowers you to effortlessly gather information from web pages and transform it into structured data. With its intuitive command-line interface and versatile capabilities, flyscrape simplifies the scraping process while delivering accurate and customizable results.

## Features

- **Simple and Intuitive**: **flyscrape** offers an easy-to-use command-line interface that allows you to interact with scraping scripts effortlessly.

- **Create New Scripts**: The `new` command enables you to generate sample scraping scripts quickly, providing you with a solid starting point for your scraping endeavors.

- **Run Scripts**: Execute your scraping script using the `run` command, and watch as **flyscrape** retrieves and processes data from the specified website.

- **Watch for Development**: The `watch` command allows you to watch your scraping script for changes and quickly iterate during development, helping you find the right data extraction queries.

## Installation

To install **flyscrape**, follow these simple steps:

1. Install Go: Make sure you have Go installed on your system. If not, you can download it from [https://golang.org/](https://golang.org/).

2. Install **flyscrape**: Open a terminal and run the following command:

   ```bash
   go install github.com/philippta/flyscrape@latest
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

The `watch` command allows you to watch your scraping script for changes and quickly iterate during development:

```bash
flyscrape watch example.js
```

## Example Script

Below is an example scraping script that showcases the capabilities of **flyscrape**:

```javascript
import { parse } from 'flyscrape';

export const options = {
    url: 'https://news.ycombinator.com/',
    depth: 1,
    allowedDomains: ['news.ycombinator.com'],
    blockedDomains: [],
    rate: 100,
};

export default function({ html, url }) {
    const $ = parse(html);
    const title = $('title');
    const entries = $('.athing').toArray();

    if (!entries.length) {
        return null;
    }

    return {
        title: title.text(),
        entries: entries.map(entry => {
            const link = $(entry).find('.titleline > a');
            const rank = $(entry).find('.rank');
            const points = $(entry).next().find('.score');

            return {
                title: link.text(),
                url: link.attr('href'),
                rank: parseInt(rank.text().slice(0, -1)),
                points: parseInt(points.text().replace(' points', '')),
            };
        }),
    };
}
```

## Contributing

We welcome contributions from the community! If you encounter any issues or have suggestions for improvement, please [submit an issue](https://github.com/philippta/flyscrape/issues).

