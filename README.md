### What's watch application

### Example commands
```
# run application from cli
make build
./whatswatchcmd api --config ./config.yaml

# check response
curl 127.0.0.1:8081/random 2>/dev/null|jq .
{
  "title": "Только представь!",
  "genre": "(драма)",
  "poster_link": "https://avatars.mds.yandex.net/get-kinopoisk-image/1600647/f7ae2460-1e86-4949-afa1-defedbc48065/600x900",
  "rating_kp": [
    "7.552 (2 879)"
  ],
  "rating_imdb": [
    "IMDb: 7.402 542"
  ],
  "country": [
    "Польша",
    "Португалия",
    "Франция"
  ],
  "link_to_kp": "https://www.kinopoisk.ru/level/1/film/586873/"
}


Main main commands
# build
make build
# build docker container
make docker
# run go fmt
make fmt
```
