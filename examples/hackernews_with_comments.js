export const config = {
  url: "https://news.ycombinator.com/",
};

export default function({ doc, scrape }) {
  const post = doc.find(".athing.submission").first();
  const title = post.find(".titleline > a").text();
  const commentsLink = post.next().find("a").last().attr("href");

  const comments = scrape(commentsLink, function({ doc }) {
    return doc.find(".comtr").map(comment => {
      return {
        author: comment.find(".hnuser").text(),
        text: comment.find(".commtext").text(),
      };
    });
  });

  return {
    title,
    comments,
  };
}
