---
title: 'Installation'
weight: 2
---

## Recommended

The easiest way to install Flyscrape is to use the following command. Note: This only works on macOS, Linux and WSL (Windows Subsystem for Linux).

```bash {filename="Terminal"}
curl -fsSL https://flyscrape.com/install | bash
```

## Alternative 1: Homebrew (macOS)

If you are on macOS, you can install Flyscrape via [Homebrew](https://formulae.brew.sh/formula/flyscrape).

```bash {filename="Terminal"}
brew install flyscrape
```

Otherwise you can download and install Flyscrape by using one of the pre-compiled binaries.

## Alternative 2: Manual installation (all systems)

Whether you are on macOS, Linux or Windows you can download one of the following archives 
to your local machine or visit the [releases page on Github](https://github.com/philippta/flyscrape/releases).

#### macOS
{{< cards >}}
{{< card link="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_0.9.0_darwin_arm64.tar.gz" title="macOS (Apple Silicon)" icon="cloud-download" >}}
{{< card link="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_0.9.0_darwin_amd64.tar.gz" title="macOS (Intel)" icon="cloud-download" >}}
{{< /cards >}}

#### Linux
{{< cards >}}
{{< card link="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_0.9.0_linux_amd64.tar.gz" title="Linux" icon="cloud-download" >}}
{{< card link="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_0.9.0_linux_arm64.tar.gz" title="Linux (arm64)" icon="cloud-download" >}}
{{< /cards >}}

#### Windows
{{< cards >}}
{{< card link="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_0.9.0_windows_amd64.zip" title="Windows" icon="cloud-download" >}}
{{< card link="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_0.9.0_windows_arm64.zip" title="Windows" icon="cloud-download" >}}
{{< /cards >}}

### Unpack

Unpack the downloaded archive by double-clicking on it or using the command line:

```bash {filename="Terminal"}
tar xf flyscrape_<os>_<arch>.tar.gz
```

After unpacking you should find a folder with the same name as the archive, which contains the `flyscrape` executable.
Change directory into it using:

```bash {filename="Terminal"}
cd flyscrape_<os>_<arch>/
```

### Install

In order to make the `flyscrape` executable globally available, you can move it to either location in your `$PATH` variable.
A good default location for that is `/usr/local/bin`. So move it using the following command:

```bash {filename="Terminal"}
mv flyscrape /usr/local/bin/flyscrape
```

### Verify

From here on you should be able to run `flyscrape` from any directory on your machine. To verify you can run the following command.
If everything went to plan you should see Flyscrapes help text:

```text {filename="Terminal"}
flyscrape --help
```
```text {filename="Terminal"}
flyscrape is a standalone and scriptable web scraper for efficiently extracting data from websites.

Usage:

    flyscrape <command> [arguments]

Commands:

    new    creates a sample scraping script
    run    runs a scraping script
    dev    watches and re-runs a scraping script
```
