# Internet Podcasts Database

## Description

Database of Podcasts from all over the internet

## Prerequisite

1. Install [Golang](https://go.dev/doc/install)
2. Install [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
3. Install [Docker](https://docs.docker.com/engine/install/)
4. Install [Make](https://www.gnu.org/software/make/#download)

## Installation

```
git clone https://github.com/terajari/ipdb.git

cd ipdb

cp env.example .env

make docker/up

make run/api
```
