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
