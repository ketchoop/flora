#!/usr/bin/env bash

set -e

function get_latest_ver() {
    local latest_release_info="$(curl -s https://api.github.com/repos/ketchoop/flora/releases/latest)"

    local latest_ver=$(echo "$latest_release_info" | grep 'tag_name' | cut -d: -f2 | cut -d\" -f2)

    echo "$latest_ver"
}

function main() {
    local FLORA_DIR="$HOME/.flora"
    local DOWNLOAD_URL="https://github.com/ketchoop/flora/releases/download"

    local kernel="$(uname -s | tr '[:upper:]' '[:lower:]')"
    local arch="amd64"
    echo "Installing flora"

    echo "Getting latest version tag"

    local latest_ver="$(get_latest_ver)"

    local latest_release_download_url="$DOWNLOAD_URL/$latest_ver/flora-$latest_ver-$kernel-$arch.tar.gz"

    pushd /tmp > /dev/null
        echo "Downloading flora $latest_ver"
        
        curl -sLO "$latest_release_download_url"

        echo "Unpacking "

        tar -xzf flora-$latest_ver-$kernel-$arch.tar.gz

        mv $kernel-$arch/flora /usr/local/bin
    popd > /dev/null

    echo "Making flora dir: $FLORA_DIR"

    mkdir -p $FLORA_DIR/bin

    echo "Flora was succesfully installed."
    echo "Please update your PATH env var in rc(e.g. .bashrc) script to point to $FLORA_DIR"
    echo "export PATH=\$PATH:$FLORA_DIR/bin"
}

main
