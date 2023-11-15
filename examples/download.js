import { download } from "flyscrape/http";

export const config = {
  url: "https://commons.wikimedia.org/wiki/London",
};

export default function ({ doc }) {
  const symbols = doc.find("#mw-content-text .mw-gallery-traditional:first-of-type li");

  return {
    symbols: symbols.map(symbol => {
      const name = symbol.text().trim();
      const url = symbol.find("img").attr("src");
      const file = `symbols/${basename(url)}`;

      download(url, file);

      return { name, url, file };
    })
  };
}

function basename(path) {
  return path.split("/").slice(-1)[0];
}
