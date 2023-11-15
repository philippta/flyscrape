export const config = {
  urls: [
    "https://news.ycombinator.com/show",
    "https://news.ycombinator.com/ask",
  ],
};

export default function({ doc, absoluteURL }) {
  const posts = doc.find(".athing");

  return {
    posts: posts.map((post) => {
      const link = post.find(".titleline > a");
      const meta = post.next();

      return {
        url: absoluteURL(link.attr("href")),
        user: meta.find(".hnuser").text(),
        title: link.text(),
        points: meta.find(".score").text().replace(" points", ""),
        created: meta.find(".age").attr("title"),
      };
    }),
  };
}
