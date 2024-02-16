#!/bin/bash

set -e

error() {
    echo -e "error:" "$@" >&2
    exit 1
}

if [[ ${OS:-} = Windows_NT ]]; then
    error "This installer does not support Windows."
fi

echo "Installing flyscrape"

case $(uname -ms) in
'Darwin x86_64')
    target=macos_amd64
    ;;
'Darwin arm64')
    target=macos_arm64
    ;;
'Linux aarch64' | 'Linux arm64')
    target=linux_arm64
    ;;
'Linux x86_64' | *)
    target=linux_amd64
    ;;
esac

dir="$HOME/.flyscrape"

mkdir -p "$dir" ||
    error "Failed to create directory: $HOME/.flyscrape"


archive="$dir/flyscrape_$target.tar.gz"
url="https://github.com/philippta/flyscrape/releases/latest/download/flyscrape_$target.tar.gz"
curl --fail --location --progress-bar --output "$archive" "$url" ||
    error "Failed to download flyscrape from: $url"

tar -xzf "$archive" -C "$dir" ||
    error "Failed to extract downloaded archive."

chmod +x "$dir/flyscrape" ||
    error "Failed to chmod the flyscrape executable."

rm "$archive" "$dir/README.md" "$dir/LICENSE" ||
    error "Failed to clean up the downloaded archive."

case $(basename "$SHELL") in
zsh)
    # Add paths to zsh
    if [[ ":$PATH:" != *":$HOME/.flyscrape:"* ]]; then
        if [[ -w "$HOME/.zshrc" ]]; then
            echo "# flyscrape" >> "$HOME/.zshrc"
            echo "export PATH=\"$dir:\$PATH\"" >> "$HOME/.zshrc"
        else
            echo ""
            echo "Manually add the directory to ~/.zshrc (or similar):"
            echo "  export PATH=\"$dir:\$PATH\""
        fi
    fi
    ;;
bash)
    # Add paths to bbash
    if [[ ":$PATH:" != *":$HOME/.flyscrape:"* ]]; then
        if [[ -w "$HOME/.bashrc" ]]; then
            echo "# flyscrape" >> "$HOME/.bashrc"
            echo "export PATH=$dir:\$PATH" >> "$HOME/.bashrc"
        else
            echo ""
            echo "Manually add the directory to ~/.bashrc (or similar):"
            echo "  export PATH=$dir:\$PATH"
        fi
    fi
    ;;
*)
    echo ""
    echo "Manually add the directory to ~/.bashrc (or similar):"
    echo "  export PATH=$dir:\$PATH"
    ;;
esac

echo ""
echo "The installation was successfull!"
echo ""
echo "Note:"
echo "Please restart your terminal window. This ensures your system correctly detects flyscrape."
echo ""
