export const config = {
  url: "https://news.ycombinator.com/",
  depth: 9,
  cache: "file",
  follow: ["a.morelink[href]"],
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
