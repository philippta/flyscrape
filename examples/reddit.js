export const config = {
  url: "https://old.reddit.com/",
};

export default function({ doc, absoluteURL }) {
  const posts = doc.find("#siteTable .thing:not(.promoted)");

  return {
    posts: posts.map((post) => {
      const rank = post.find(".rank");
      const user = post.find(".author");
      const created = post.find("time");
      const title = post.find("a.title");
      const comments = post.find(".comments");
      const subreddit = post.find(".subreddit");
      const upvotes = post.find(".score.unvoted");
      const thumbnail = post.find("a.thumbnail img");

      return {
        rank: rank.text(),
        user: user.text(),
        created: created.attr("datetime"),
        title: title.text(),
        link: absoluteURL(title.attr("href")),
        comments: comments.text().replace(" comments", ""),
        comments_link: comments.attr("href"),
        subreddit: subreddit.text(),
        upvotes: upvotes.text(),
        thumbnail: absoluteURL(thumbnail.attr("src")),
      };
    }),
  };
}
