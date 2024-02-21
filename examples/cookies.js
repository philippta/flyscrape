export const config = {
  url: "https://news.ycombinator.com/",

  // This will use cookies from your Chrome browser.
  // Options: "chrome" | "firefox" | "edge"
  cookies: "chrome",
};

export default function({ doc }) {
  return {
    user: doc.find("#me").text(),
    karma: doc.find("#karma").text(),
  }
}
