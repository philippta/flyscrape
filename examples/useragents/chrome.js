import { parse } from "flyscrape";

export const config = {
  url: "https://chromereleases.googleblog.com/search/label/Stable%20updates",
  follow: [".blog-pager-older-link"],
  depth: 30,
  cache: "file",
};

export default function ({ doc, absoluteURL }) {
  const posts = doc.find(".post");
  return posts.map(post => {
    const title = post.find("h2").text().trim();
    const body = parse(post.find(".post-content").text()).find("p:nth-child(1)").text().trim();

    const regexes = [
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (Mac)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (Windows)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (Linux)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (iOS)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (Android)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (ChromeOS)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (Mac,Linux)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)? for (Mac and Linux)/,
       /(\d+\.\d+\.\d+\.\d+(\/\.\d+)?)\)?\s\(Platform version:\s[\d\.]+\)\sfor\smost\s(ChromeOS)/,
    ];

    const versions = new Set();
    for (const regex of regexes) {
      const matches = body.match(regex);
      if (!matches) {
        continue;
      }

      let versionStr = matches[1];

      let vv = versionStr.split("/");
      if (vv.length == 2) {
        vv[1] = vv[0].substring(0, vv[0].lastIndexOf(".")) + vv[1];
      }

      for (const version of vv) {
        versions.add(version)
      }
    }


    return versions
  }).filter(Boolean).flat();
}
