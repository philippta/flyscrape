export const config = {
  url: "https://news.ycombinator.com/",
  depth: 2,
  follow: [],
};

export default function({ url, doc, follow }) {
  const next = doc.find(".morelink").attr("href");

  follow(next);

  return { url, next };
}
