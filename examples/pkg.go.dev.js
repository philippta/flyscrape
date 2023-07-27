import { parse } from "flyscrape";

export const options = {
    url: "https://pkg.go.dev/github.com/stretchr/testify/require",
}

export default function({ html }) {
    console.log("start")
    const $ = parse(html);
    console.log("done")

    return {
        package: $('h1').text().trim(),
        meta: {
            version: $('[data-test-id=UnitHeader-version] > a').text().replace("Version: ", "").trim(),
            license: $('[data-test-id=UnitHeader-licenses] > a').text().trim(),
            published: $('[data-test-id=UnitHeader-commitTime]').text().replace("Published: ", "").trim(),
            imports: $('[data-test-id=UnitHeader-imports] > a').text().replace("Imports: ", "").trim(),
            importedBy: $('[data-test-id=UnitHeader-importedby] > a').text().replace("Imported by: ", "").replace(/,/g,"").trim(),
        },
        functions: $('.Documentation-indexList .Documentation-indexFunction > a').toArray().map(el => $(el).text()),
        types: $('.Documentation-indexList .Documentation-indexType > a').toArray().map(el => $(el).text()),
    };
}
