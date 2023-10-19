# Running flyscrape

Once you've successfully installed flyscrape, you're ready to start using it to scrape data from websites. This section will guide you through the basic commands and steps to run flyscrape effectively.

## Creating Your First Scraping Script

Before you can run flyscrape, you'll need a scraping script to tell it what data to extract from a website. To create a new scraping script, you can use the `new` command followed by the script's filename. Here's an example:

```bash
flyscrape new my_first_script.js
```

This command will generate a sample scraping script named `my_first_script.js` in the current directory. You can edit this script to customize it according to your scraping needs.

## Running a Scraping Script

To execute a scraping script, you can use the `run` command followed by the script's filename. For example:

```bash
flyscrape run my_first_script.js
```

When you run this command, flyscrape will start retrieving and processing data from the website specified in your script.

## Watching for Development

During the development phase, you may want to make changes to your scraping script and see the results quickly. flyscrape provides a convenient way to do this using the `dev` command. It allows you to watch your scraping script for changes and automatically re-run it when you save your changes.

Here's how to use the `dev` command:

```bash
flyscrape dev my_first_script.js
```

With the development mode active, you can iterate on your scraping script, fine-tune your data extraction queries, and see the results in real-time.

## Script Output

After running a scraping script, flyscrape will generate structured data based on your script's logic. You can customize your script to specify what data you want to extract and how it should be formatted.

Congratulations! You've learned how to run flyscrape and create your first scraping script. You can now start gathering data from websites and transforming it into structured information.

Next, explore more advanced topics in the "Using flyscrape" section to refine your scraping skills and learn how to set up more complex scraping scenarios.
