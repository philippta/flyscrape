import { parse } from "flyscrape";

export const options = {
    url: "https://news.ycombinator.com/",     // Specify the URL to start scraping from.
    depth: 1,                                 // Specify how deep links should be followed (0 = no follow).
    allowedDomains: ["news.ycombinator.com"], // Specify the allowed domains to follow.
    rate: 100,                                // Specify the request rate in requests per second.
}

export default function({ html, url }) {
    const $ = parse(html);

    return {
        title: $('title').text(),
        entries: $('.athing').toArray().map(entry => {
            const link = $(entry).find('.titleline > a');
            const rank = $(entry).find('.rank');
            const points = $(entry).next().find('.score');

            return {
                title: link.text(),                                     // Extract the title text.
                url: link.attr('href'),                                 // Extract the link href.
                rank: parseInt(rank.text().slice(0, -1)),               // Extract and cleanup the rank.
                points: parseInt(points.text().replace(' points', '')), // Extract and cleanup the points.
            }
        }),
    };
}
