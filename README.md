# flora
[![Build Status](https://travis-ci.org/ketchoop/flora.svg?branch=master)](https://travis-ci.org/ketchoop/flora)

Version/upgrade manager for terraform

## What is it

*Flora* is another and missing **upgrade and version manager** for Terraform. Inspired by *tfenv* and written in go.
You can use it to upgrade your Terraform version *by one command*. Another use case: to switch between Terraform versions to use 
compatible with your **.tf manifests** Terraform binary.

## Features

* Upgrade your Terraform by one command
* Switch between Terraform versions easy, fast and without pain
* Bash/Zsh autocompletion. Even for versions.

## Install


1. By *go get*
```
go get -u github.com/ketchoop/flora/cmd/flora

mkdir -p ~/flora/.bin
```

2. By install.sh
```
curl https://raw.githubusercontent.com/ketchoop/flora/master/install.sh | bash
```

All of installation ways require existing of `~/.flora/bin` directory and path to it  in your `PATH`. So...


Add to `PATH`:

```
echo 'export PATH=$PATH:$HOME/.flora/bin' >> .path_to_your_rc_file # .bashrc, .zshrc and so on

source .path_to_your_rc_file # To update state of PATH env
```

## Usage(short description)

```
NAME:
   flora - Simple app to upgrade your terraform

USAGE:
   flora [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     upgrade   Upgrade terraform
     download  Download specific Terraform version
     use       Download(when it's needed) and use specific terraform version
     versions  List all available terraform versions
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## How it works

Like rbenv like tools it downloads binary of Terraform and links it to special folder (`~/.flora/bin`), which path have to be in you `PATH` env.
When you switch between versions flora links another version to bin folder.
