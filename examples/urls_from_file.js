import urls from "./urls.txt"

export const config = {
  urls: urls.split("\n")
};

export default function({ doc }) {
  return {
    title: doc.find("title").text().trim(),
  };
}
