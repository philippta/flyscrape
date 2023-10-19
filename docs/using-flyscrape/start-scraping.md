# Start Scraping

In this section, we'll dive into the process of initiating web scraping using flyscrape. Now that you have created and fine-tuned your scraping script, it's time to run it and start gathering data from websites.

## The `flyscrape run` Command

The `flyscrape run` command is used to execute your scraping script and retrieve data from the specified website. This command is your gateway to turning your scraping logic into actionable results.

## Running Your Scraping Script

To run your scraping script, simply use the `flyscrape run` command followed by the name of your script file. For example:

```bash
flyscrape run my_scraping_script.js
```

This command will initiate the scraping process as defined in your script. Flyscrape will execute your script and stream the JSON output of the extracted data directly to your terminal.

## Saving Scraped Data to a File

You can easily save the JSON output of the scraped data to a file using standard shell redirection. For example, to save the scraped data to a file named `result.json`, you can use the following command:

```bash
flyscrape run my_scraping_script.js > result.json
```

This command will execute your scraping script and save the extracted data in the `result.json` file in the current directory.

## Example Workflow

Here's a simple workflow for starting web scraping with flyscrape, including saving the scraped data to a file:

1. Create a scraping script using `flyscrape new` and fine-tune it using `flyscrape dev`.

2. Save your script.

3. Run the script using `flyscrape run`.

4. Observe the terminal as flyscrape streams the JSON output of the extracted data in real-time.

5. If you want to save the data to a file, use redirection as shown above (`flyscrape run my_scraping_script.js > result.json`).

6. Customize the script to store, process, or further analyze the extracted data as needed.

7. Continue scraping or iterate on your script for more complex scenarios.

With this workflow, you can efficiently gather and process data from websites using flyscrape, with the option to save the extracted data to a file for later use or analysis.

---

This concludes the "Start Scraping" section, which covers the process of initiating web scraping with the `flyscrape run` command, including an example of how to save the scraped data to a file. Next, you can explore various configuration options and advanced features in the "Options" section to further tailor your scraping experience.
