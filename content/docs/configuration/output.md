---
title: 'Output File and Format'
weight: 10
---

The output file and format are specified in the configuration object of your scraping script. They determine where the scraped data will be saved and in what format.

## Output File

The output file is the file where the scraped data will be saved. If not specified, the data will be printed to the standard output (stdout).

```javascript {filename="Configuration"}
export const config = {
    output: {
        // Specify the output file.
        file: "results.json",
    },
};
```

In the above example, the scraped data will be saved in a file named `results.json`.

## Output Format

The output format is the format in which the scraped data will be saved. The options are `json` and `ndjson`.

```javascript {filename="Configuration"}
export const config = {
    output: {
        // Specify the output format.
        format: "json",
    },
};
```

In the above example, the scraped data will be saved in JSON format.

```javascript {filename="Configuration"}
export const config = {
    output: {
        // Specify the output format.
        format: "ndjson",
    },
};
```

In this example, the scraped data will be saved in newline-delimited JSON (NDJSON) format. Each line in the output file will be a separate JSON object.
