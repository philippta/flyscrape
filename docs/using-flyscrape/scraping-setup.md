# Scraping Setup

In this section, we'll delve into the details of setting up your scraping script using the `flyscrape new script.js` command. This command is designed to streamline the process of creating a scraping script, providing you with a structured starting point for your web scraping endeavors.

## The `flyscrape new` Command

The `flyscrape new` command allows you to generate a new scraping script with a predefined structure and sample code. This is incredibly helpful because it provides a quick and easy way to begin your web scraping project.

## Creating a New Scraping Script

To create a new scraping script, use the `flyscrape new` command followed by the desired script filename. For example:

```bash
flyscrape new my_scraping_script.js
```

This command will generate a file named `my_scraping_script.js` in the current directory. You can then open and edit this file with your preferred code editor.

## Script Overview

Let's take a closer look at the structure and components of the generated scraping script:

```javascript
import { parse } from 'flyscrape';

export const options = {
	url: 'https://example.com/', // Specify the URL to start scraping from.
	depth: 1, // Specify how deep links should be followed. (default = 0, no follow)
	allowedDomains: [], // Specify the allowed domains. ['*'] for all. (default = domain from url)
	blockedDomains: [], // Specify the blocked domains. (default = none)
	allowedURLs: [], // Specify the allowed URLs as regex. (default = all allowed)
	blockedURLs: [], // Specify the blocked URLs as regex. (default = non-blocked)
	proxy: '', // Specify the HTTP(S) proxy to use. (default = no proxy)
	rate: 100 // Specify the rate in requests per second. (default = 100)
};

export default function ({ html, url }) {
	const $ = parse(html);

	// Your data extraction logic goes here

	return {
		// Return the structured data you've extracted
	};
}
```

## Implementing the Data Extraction Logic

In the generated scraping script, you'll find the comment "// Your data extraction logic goes here." This is the section where you should implement your custom data extraction logic. You can use tools like [Cheerio](https://cheerio.js.org/) or other libraries to navigate and extract data from the parsed HTML.

Here's an example of how you might replace the comment with data extraction code:

```javascript
// Your data extraction logic goes here
const title = $('h1').text();
const description = $('p').text();

// You can extract more data as needed
```

## Returning the Extracted Data

After implementing your data extraction logic, you should structure the data you've extracted and return it from the scraping function. The comment "// Return the structured data you've extracted" is where you should place this code.

Here's an example of how you might return the extracted data:

```javascript
return {
	title: title,
	description: description
	// Add more fields as needed
};
```

With this setup, you can effectively scrape and structure data from web pages to meet your specific requirements.

---

This concludes the "Scraping Setup" section, which provides insights into creating scraping scripts using the `flyscrape new` command, implementing data extraction logic, and returning extracted data. Next, you can explore more advanced topics in the "Development Mode" section to streamline your web scraping workflow.
