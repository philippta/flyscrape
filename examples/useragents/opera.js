
export const config = {
  urls: range("https://blogs.opera.com/desktop/changelog-for-{}/", 60, 110),
};

export default function ({ doc, absoluteURL }) {
  const versions = doc.find(".content h4");
  return versions.map(versions => {
    return versions.text().split(" ")[0].trim();
  }).filter(Boolean);
}

function range(url, from, to) {
  return Array.from({length: to - from + 1}).map((_, i) => url.replace("{}", i + from));
}
