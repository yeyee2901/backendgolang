# Running Instructions

The app is "packaged" using Docker & `docker compose`. I have prepared several make file targets to make things easier:

```bash
# cleans any log & (IMPORTANT) docker database volume
make clean

# build the docker images
make docker-build

# fire up the containers, runs in background
make docker-up

# shutdown
make docker-down 
```

Server can be accessed with host: `127.0.0.1:3000`. The documentation (swagger OpenAPI) can be accessed from `127.0.0.1:3000/swagger/index.html`

As an additional note, we can change the scheduler interval from `docker-compose.yml`:
```yaml
...
...
  scheduler:
    build:
      context: .
      dockerfile: scheduler.Dockerfile
    restart: always

    # NOTE: set the scheduling here
    environment:
      - SCHEDULER_REFRESH_URL=http://api:3000/api/v1/indego-data-fetch-and-store-it-db
      - SCHEDULER_REFRESH_INTERVAL=1h
      # - SCHEDULER_REFRESH_INTERVAL=3m
    networks:
      - epiphyte-corp-dot-net
    depends_on:
      - postgresdb
...
...
```

Application logs can be observed from `docker_log/app.log`. This has a `service` key that can be used to identify which service produce which log.
```bash
tail -f docker_log/app.log
```

# Main Techs Used 
- [Gin](https://github.com/gonic-gin/gin) - Fast & Simple HTTP Web Server Framework for Golang.
- [golang-migrate](https://github.com/golang-migrate/migrate) - manage sequential migration scripts, can be used as standalone CLI, a library, and there is also docker image available for it, which is a good thing for `docker compose` situation.
- For logging, I use a combination of [lumberjack](https://github.com/natefinch/lumberjack) & `slog` (builtin `slog` package).

## Golang  Backend Challenge

[Indego](https://www.rideindego.com) is Philadelphia's bike-sharing program, with many bike stations in the city.

The [Indego GeoJSON station status API](https://www.rideindego.com/stations/json/) provides a realtime snapshot of the number of bikes available, number of open docks available (not currently containing a bike), and total number of docks at every station. This API is free and requires no API key.

The [Open Weather Map API](https://openweathermap.org/current#name) provides a realtime snapshot of the current weather in a given city. Since Philadelphia is a small geographical area it is sufficient to obtain the weather for a geographical location central to Philadelphia. This API has a free plan, you will need to sign up for an API key.

Using Golang, create a new API server which accumulates data over time and provides access to historical data for both weather and Indego bike availability.

## API endpoints

Use a static token and protect all the endpoints. If that static token is not provided or is invalid return error response with relevant HTTP status code.

### Store data from Indego

An endpoints which downloads fresh data from [Indego GeoJSON station status API](https://www.rideindego.com/stations/json/) and stores it inside PostgreSQL.

```bash
# this endpoint will be trigger every hour to fetch the data and insert it in the PostgreSQL database
POST http://localhost:3000/api/v1/indego-data-fetch-and-store-it-db
```

### Snapshot of all stations at a specified time

Data for all stations as of 11am Universal Coordinated Time on September 1st, 2019:

```bash
GET http://localhost:3000/api/v1/stations?at=2019-09-01T10:00:00Z
```

This endpoint should respond as follows, with the actual time of the first snapshot of data on or after the requested time and the data:

```javascript
{
  at: '2019-09-01T10:00:00Z',
  stations: { /* As per the Indego API */ },
  weather: { /* As per the Open Weather Map API response for Philadelphia */ }
}
```

### Snapshot of one station at a specific time

Data for a specific station (by its `kioskId`) at a specific time:

```bash
GET http://localhost:3000/api/v1/stations/{kioskId}?at=2019-09-01T10:00:00Z
```

The response should be the first available on or after the given time, and should look like:

```javascript
{
  at: '2019-09-01T10:00:00',
  station: { /* Data just for this one station as per the Indego API */ },
  weather: { /* As per the Open Weather Map API response for Philadelphia */ }
}
```

Include an `at` property in the same format indicating the actual time of the snapshot.

If no suitable data is available a 404 status code should be given.

## Unit tests

Write functional and API test the tool of your choice.

## Hosting details

You will need to make your API available on a server running in local environment.

## Criteria

Your work will be evaluated primarily on:

* Correct use of proper error handling
* Definitions, types, interfaces.
* Efficient Postgres queries and indexes
* Unit tests for API endpoints and also do functional testing
* API documentation using OpenAPI specification

## How to submit your work

Create a public repo on Github and push your code on it. then share the link back with the team.

## Extra credit

* oAuth2 integration to protect your APIs by registering Auth0 free account.
* A simple front end React application offering a visualization of all or part of the data utilizing the API you have built as a back end.
* Anything else you think is cool, relevant, and consistent with the other requirements.
