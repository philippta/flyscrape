export const config = {
  url: "https://www.airbnb.com/",
  browser: true,
  // headless: false,
};

export default function ({ doc, absoluteURL }) {
  const rooms = doc.find("[itemprop=itemListElement]");

  return {
    listings: rooms.map(room => {
      const link = "https://" + room.find("meta[itemprop=url]").attr("content");
      const image = room.find("img").attr("src");
      const desc = new Set(room.find("[role=group] > div > div > div").map(d => d.text()).filter(Boolean));

      return { link, image, desc }
    }),
  }
}
