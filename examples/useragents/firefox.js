import { parse } from "flyscrape";

export const config = {
  url: "https://www.mozilla.org/en-US/firefox/releases/",
};

export default function ({ doc, absoluteURL }) {
  const links = doc.find(".c-release-list a");
  return links
    .map(link => link.text())
    .filter(Boolean)
    .filter(version => parseFloat(version) >= 60);
}
