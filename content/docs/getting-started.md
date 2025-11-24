---
title: 'Getting Started'
weight: 1
---

In this quick guide we will go over the core functionalities of Flyscrape and how to use it. Make sure you've got `flyscrape` up and running on your system.

The quickest way to install Flyscrape on Mac, Linux or WSL is to run the following command. For more information or how to install it on Windows check out the [installation instructions](/docs/installation).

```bash {filename="Terminal"}
curl -fsSL https://flyscrape.com/install | bash
```

## Overview

Flyscrape is a standalone scraping tool tool that works with so called _scraping scripts_.

Scraping scripts let you define what data you want to extract from a website using familiar JavaScript code you might recognize from jQuery or cherrio. Inside your scraping script, you can also configure how the Flyscrape should behave, e.g. what links to follow, what domains to access, how fast to send out requests, etc.

When your happy with the initial version of your scraping script, you can run Flyscrape and it will go off and start scraping the websites you have defined.

## Your first Scraping Script

A new scraping script can be created using the `new` command. This script is meant as a helpful guide to let you explore the JavaScript API.

Go a head and run the following command:
```sh {filename="Terminal"}
flyscrape new hackernews.js
```

This should have created you a new file called `hackernews.js` in your current directory. You can open it up in your favorite text editor.

## Anatomy of a Scraping Script

Let's look at the previously created `hackernews.js` file and go through it together. Every scraping script consists of two main parts:

### Configuration 

The configuration is used to control the scraping behaviour. Here we can specify what URLs to scrape, how deep it should follow links or what domains should be allowed to acess. Besides these, there are a bunch more to explore.

```javascript {filename="Configuration"}
export const config = {
  url: "https://hackernews.com",
  // depth: 0,
  // allowedDomains: [],
  // ...
}
```

### Data Extraction Logic

The data extracting logic defines what data to extract from a website. In this example it grabs the posts from the website using the `doc` document object and extracts the individual links and their titles. The `absoluteURL` function is used to ensure that every relative link is converted into an absolute one.

```javascript {filename="Data Extraction Logic"}
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

### Starting the Development Mode

Flyscrape has a built in Development Mode that allows you to quickly iterate and see changes to your script immediately.
It does so by watching your script for changes and re-runs the Data Extraction Logic against a cached version of the website.

Let's try and fire that up using the following command:

```sh {filename="Terminal"}
flyscrape dev hackernews.js
```

You should now see the extracted data of your target website. Note that no links are followed in this mode, even when otherwise specified in the configuration.

Now let's try and change our script so we extract some more data like the user, who submitted the post.

```diff {filename="hackernews.js"}
    return {
      title: title.text(),
      posts: posts.map((post) => {
        const link = post.find(".titleline > a");
+       const meta = post.next();
  
        return {
          title: link.text(),
          url: absoluteURL(link.attr("href")),
+         user: meta.find(".hnuser").text(),
        };
      }),
    };
```

When you now save the file and look at your terminal again, the changes should have reflected and the user added to each of the posts.

Once you're happy with the extraction logic, your can exit out by pressing CTRL+C.

### Running the Scraper

Now that your scraping script is configured and the extraction logic is in place, your can use the `run` command to execute the scraper.

```sh {filename="Terminal"}
flyscrape run hackernews.js
```

This should output a JSON array of all scraped pages.

### Learn more

Once you're done experimenting feel fee to check many of Flyscrape's other features. There are plenty to customize it for your specific needs.

- [Full Example Script](/docs/full-example-script)
- [API Reference](/docs/api-reference)
