import { parse } from "flyscrape";

export const options = {
    url: "https://esbuild.github.io/plugins/",
    depth: 1,
    allowedDomains: [
        "esbuild.github.io", 
        "nodejs.org",
    ],
}

export default function({ html }) {
    const doc = parse(html);

    return {
        headline: doc('h1').text().trim(),
        body: doc('main > p:nth-of-type(1)').text().trim(),
    };
}
