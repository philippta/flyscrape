export const config = {
  url: "https://coinmarketcap.com/",
  follow: ["a[href]"],
  depth: 1,
  output: {
    mongodb: {
      uri: "mongodb://localhost:27017",
      database: "test",
      collection: "coinmarketcap",
      maxPoolSize: 100,
    },
  },
};

export default function ({ doc }) {
  const title = doc.find("title");

  return {
    title: title.text(),
  };
}
