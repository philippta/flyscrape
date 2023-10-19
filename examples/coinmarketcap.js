export const config = {
  url: "https://coinmarketcap.com/",
};

export default function({ doc }) {
  const rows = doc.find(".cmc-table tbody tr");

  return {
    currencies: rows
      .map((row) => {
        const cols = row.find("td");

        return {
          position: cols.get(1).text(),
          currency: cols.get(2).find("p").get(0).text(),
          symbol: cols.get(2).find("p").get(1).text(),
          price: cols.get(3).text(),
          change: {
            "1h": cols.get(4).text(),
            "24h": cols.get(5).text(),
            "7dh": cols.get(6).text(),
          },
          marketcap: cols.get(7).find("span").get(1).text(),
          volume: cols.get(8).find("p").get(0).text(),
          supply: cols.get(9).text(),
        };
      })
      .slice(0, 10),
  };
}
