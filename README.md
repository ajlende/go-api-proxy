# go-api-proxy

A Go proxy app, which can easily be deployed to Heroku.

## Dependencies

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

## Environment

Create a `.env` file with the following variables:

```sh
GITHUB_ACCESS_TOKEN=<your-github-access-token>
GITHUB_API_URL=https://api.github.com
ALLOWED_ORIGINS=https://one.example.com,https://two.example.com
```

## Running Locally

```sh
git clone https://github.com/ajlende/go-api-proxy.git
cd go-api-proxy
make local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
heroku create
make deploy
heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
