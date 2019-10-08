#!/usr/bin/env bash

set -ex

function get_latest_ver() {
    local latest_release_info="$(curl -s https://api.github.com/repos/ketchoop/flora/releases/latest)"

    local latest_ver=$(echo "$latest_release_info" | grep 'tag_name' | cut -d: -f2 | cut -d\" -f2)

    echo "$latest_ver"
}

function install_autocompletion() {
    case $(basename $SHELL) in
    zsh)
        local zsh_autocomplete_path="/usr/local/share/zsh/site-functions"

        pushd "$zsh_autocomplete_path" > /dev/null
            curl -sL "https://raw.githubusercontent.com/ketchoop/flora/master/configs/autocomplete/flora_zsh_autocomplete" -o _flora
        popd
        ;;
    bash)
        local bash_autocomplete_path="/etc/bash_completion.d"

        pushd "$bash_autocomplete_path" > /dev/null
            curl -sLO "https://raw.githubusercontent.com/ketchoop/flora/master/configs/autocomplete/flora_bash_autocomplete"
        popd
        ;;
    esac
}

function main() {
    local FLORA_DIR="$HOME/.flora"
    local DOWNLOAD_URL="https://github.com/ketchoop/flora/releases/download"
    local INSTALL_FILES_DIR="$(mktemp -d /tmp/flora_install.XXXXXX)"

    local kernel="$(uname -s | tr '[:upper:]' '[:lower:]')"
    local arch="amd64"
    echo "Installing flora"

    echo "Getting latest version tag"

    local latest_ver="$(get_latest_ver)"
    echo "$latest_ver"

    local latest_release_download_url="${DOWNLOAD_URL}/${latest_ver}/flora_${latest_ver}_${kernel}_${arch}.tar.gz"

    pushd $INSTALL_FILES_DIR > /dev/null
        echo "Downloading flora $latest_ver"
        
        curl -sLO "$latest_release_download_url"

        echo "Unpacking "

        tar -xzf "flora_${latest_ver}_${kernel}_${arch}.tar.gz"

        mv flora /usr/local/bin
    popd > /dev/null

    echo "Cleaning up"

    rm -rf $INSTALL_FILES_DIR

    echo "Installing autcompletion script"

    install_autocompletion 

    echo "Making flora dir: $FLORA_DIR"

    mkdir -p $FLORA_DIR/bin

    echo "Flora was succesfully installed."
    echo "Please update your PATH env var in rc(e.g. .bashrc) script to point to $FLORA_DIR"
    echo "export PATH=\$PATH:$FLORA_DIR/bin"
}

main
