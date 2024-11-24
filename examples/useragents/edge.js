import { parse } from "flyscrape";

export const config = {
  url: "https://learn.microsoft.com/en-us/deployedge/microsoft-edge-release-schedule",
};

export default function ({ doc, absoluteURL }) {
  const links = doc.find("table a");
  return links
    .map(link => link.text())
    .filter(Boolean)
}
