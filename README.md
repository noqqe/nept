# nept

nept is a tiny image manipulation program I write to get used to go

## Install

    brew install noqqe/tap/nept

## Usage

```
> nept -h
NAME:
   nept - A new cli application

USAGE:
   nept [global options] command [command options] [arguments...]

VERSION:
   1.0.1

DESCRIPTION:
   Image manipulation program for commandline on pixel level

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --brightness value, -b value  Brighten the Image (default: 0)
   --darkness value, -d value    Darken the image (default: 0)
   --flatten value, -f value     Flatten the image (default: 0)
   --iso value, -s value         Add iso to the image (default: 0)
   --negative, -n                Convert negative to positive image (default: false)
   --in value, -i value          Image to edit (input)
   --out value, -o value         Image to edit (output)
   --help, -h                    show help (default: false)
   --version, -v                 print the version (default: false)
```

## GitHub Actions

![Go](https://github.com/noqqe/nept/workflows/Go/badge.svg)
