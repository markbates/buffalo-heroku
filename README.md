# buffalo-heroku

This is a plugin for [https://gobuffalo.io](https://gobuffalo.io) that makes working with Heroku easier.

It assumes you are using Docker to deploy to Heroku. It is recommended you read [https://devcenter.heroku.com/articles/container-registry-and-runtime](https://devcenter.heroku.com/articles/container-registry-and-runtime) first.

## Installation

```bash
$ go get -u -v github.com/markbates/buffalo-heroku
```

## Pre-Requisites

* You should absolutely have read [https://devcenter.heroku.com/articles/container-registry-and-runtime](https://devcenter.heroku.com/articles/container-registry-and-runtime) first.
* You should have the Heroku CLI installed [https://devcenter.heroku.com/articles/heroku-cli](https://devcenter.heroku.com/articles/heroku-cli).

## Setup

The `buffalo heroku setup` command will setup and create a new Heroku app for you, with a bunch of defaults that **I** find nice.

### Flags/Options

There are a lot of flags and options you can use to tweak the Heroku app you create. Use the `--help` flag to see a list of them all.

```bash
$ buffalo heroku setup
```

### Interactive Mode

If you are unsure of what you're doing, you can use the `-i` flag to enter interactive mode. This will ask you questions and give you menus of choices to select from.

## Deploying

The initial `setup` command will do a deploy at the end, but after that you'll want to use the `buffalo heroku deploy` command to push a new version of your application, it'll even try to run your migrations for you.

```bash
$ buffalo heroku deploy
```
