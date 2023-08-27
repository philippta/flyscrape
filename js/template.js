import { parse } from 'flyscrape';

export const options = {
    url: 'https://news.ycombinator.com/',     // Specify the URL to start scraping from.
    depth: 1,                                 // Specify how deep links should be followed.  (default = 0, no follow)
    allowedDomains: [],                       // Specify the allowed domains. ['*'] for all. (default = domain from url)
    blockedDomains: [],                       // Specify the blocked domains.                (default = none)
    allowedURLs: [],                          // Specify the allowed URLs as regex.          (default = all allowed)
    blockedURLs: [],                          // Specify the blocked URLs as regex.          (default = non blocked)
    rate: 100,                                // Specify the rate in requests per second.    (default = 100)
}

export default function({ html, url }) {
    const $ = parse(html);
    const title = $('title');
    const entries = $('.athing').toArray();

    if (!entries.length) {
        return null; // Omits scraped pages without entries.
    }

    return {
        title: title.text(),                                            // Extract the page title.
        entries: entries.map(entry => {                                 // Extract all news entries.
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
