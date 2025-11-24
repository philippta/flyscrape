---
title: 'Full Example Script'
weight: 3
---

This script serves as a reference that show all features of Flyscrape and how to use them. Feel free to copy and paste this as a starter script.

```javascript{filename="Reference"}
import { parse } from "flyscrape";
import { download } from "flyscrape/http";
import http from "flyscrape/http";

export const config = {
    // Specify the URL to start scraping from.
    url: "https://example.com/",

    // Specify the multiple URLs to start scraping from.   (default = [])
    urls: [                          
        "https://anothersite.com/",
        "https://yetanother.com/",
    ],

    // Enable rendering with headless browser.             (default = false)
    browser: true,

    // Specify if browser should be headless or not.       (default = true)
    headless: false,

    // Specify how deep links should be followed.          (default = 0, no follow)
    depth: 5,                        

    // Speficy the css selectors to follow.                (default = ["a[href]"])
    follow: [".next > a", ".related a"],                      
 
    // Specify the allowed domains. ['*'] for all.         (default = domain from url)
    allowedDomains: ["example.com", "anothersite.com"],              
 
    // Specify the blocked domains.                        (default = none)
    blockedDomains: ["somesite.com"],              

    // Specify the allowed URLs as regex.                  (default = all allowed)
    allowedURLs: ["/posts", "/articles/\d+"],                 
 
    // Specify the blocked URLs as regex.                  (default = none)
    blockedURLs: ["/admin"],                 
   
    // Specify the rate in requests per minute.            (default = no rate limit)
    rate: 60,                       

    // Specify the number of concurrent requests.          (default = no limit)
    concurrency: 1,                       

    // Specify a single HTTP(S) proxy URL.                 (default = no proxy)
    // Note: Not compatible with browser mode.
    proxy: "http://someproxy.com:8043",

    // Specify multiple HTTP(S) proxy URLs.                (default = no proxy)
    // Note: Not compatible with browser mode.
    proxies: [
      "http://someproxy.com:8043",
      "http://someotherproxy.com:8043",
    ],                     

    // Enable file-based request caching.                  (default = no cache)
    cache: "file",                   

    // Specify the HTTP request header.                    (default = none)
    headers: {                       
        "Authorization": "Bearer ...",
        "User-Agent": "Mozilla ...",
    },

    // Use the cookie store of your local browser.         (default = off)
    // Options: "chrome" | "edge" | "firefox"
    cookies: "chrome",

    // Specify the output options.
    output: {
        // Specify the output file.                        (default = stdout)
        file: "results.json",
        
        // Specify the output format.                      (default = json)
        // Options: "json" | "ndjson"
        format: "json",
    },
};

export default function ({ doc, url, absoluteURL }) {
  // doc              - Contains the parsed HTML document
  // url              - Contains the scraped URL
  // absoluteURL(...) - Transforms relative URLs into absolute URLs

  // Find all users.
  const userlist = doc.find(".user")

  // Download the profile picture of each user.
  userlist.each(user => {
    const name = user.find(".name").text()
    const pictureURL = absoluteURL(user.find("img").attr("src"));

    download(pictureURL, `profile-pictures/${name}.jpg`)
  })

  // Return users name, address and age.
  return {
    users: userlist.map(user => {
      const name = user.find(".name").text()
      const address = user.find(".address").text()
      const age = user.find(".age").text()

      return { name, address, age };
    })
  };
}
```
