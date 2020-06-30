Sportbuff Assignment
====================

These instructions have been verified on an ubnuntu 10.04.4 machine
and it is probably safe to assume it works on all sane linux distros.

However: There may be some issues on other OS's (e.g. OSX or Windows)
due to different virtualisation layers.

Potential issues are probably related to network issues, and you can probably
diagnose them with the help of the env file in the.

Quickstart:
-----------

Setup the service locally

```
source ./deploy/env.sh
docker-compose -f ./deploy/localdeploy.yaml up --build -d
```

Run some queries

```
curl localhost:8000/v1/video_streams
```

Stop the service:

```
docker-compose -f ./deploy/localdeploy.yaml down --volumes
```

###Run tests

Spin-up the testing env, run the go tests, and tear it down again

```
source ./deploy/env.sh
docker-compose -f ./deploy/testing.yaml up --build
go test ./...
docker-compose -f ./deploy/testing.yaml down --volumes
```

###Example Requests

```
TODO
```

About
-----

###API

The API is implemented as a restful API served over normal HTTP
Only `GET` methods were required by this task, but on further development
the other actions (create, update, delete) can be implemented with the
other HTTP methods (`POST`, `PUT`, `DELETE`).

| route                          | method | paginated? | multi-codec |
|--------------------------------|--------|------------|-------------|
| /v1/video_streams              | GET    | True       | True        |
| /v1/video_streams/{uuid}       | GET    | False      | True        |
| /v1/video_streams/{uuid}/buffs | GET    | False      | True        |
| /v1/buffs                      | GET    | True       | True        |
| /v1/buffs/{uuid}               | GET    | False      | True        |

Pagination:

Paginated endpoints use count and skip parameters (defaulting to `count=10` and `skip=0`)
Count is the number of items to return, and skip is the number of blocks (size of count) to skip.
This is nice as you can scan the dataset by choosing a count size and incrementing the skip value
per request.

These defaults can be over-ridden by setting the URL params on the request.

E.g. `?count=6&skip=3` returning items 12-17 (indexed from 0)

Codec:

The rest API supports multi codec behaviour, returning data in JSON by default.
The codec can be changed by providing the `codec` URL param. Valid values are:

| codec URL param   | Result                 |
|-------------------|------------------------|
| codec=json        | std JSON encoding      |
| codec=json,pretty | indented JSON encoding |
| codec=yaml        | YAML encoding          |

###Database

The database is a simple postgres database. It is maintained using the migration scripts in `deploy/migrations/`
using the migration tool tern.

For convenience, there is a dbinit container (run automatically in the docker-compose) that migrates the database and runs a populate job to fill it with
fake data.

###Observability

There is a basic observability stack using opentracing which is viewable from the Jaeger
service in the docker-compose file.


###Project Structure

The project follows the repo layout suggestions from the [golang standards](https://github.com/golang-standards/project-layout)
project. For ease of digestion, I have provided an annotated tree of the repo here:

```
.
├── api
│   ├── buff
│   │   └── [handlers for the buff subtype]
│   ├── types
│   │   └── [exposed API types (data model the API serves)]
│   ├── videostream
│   │   └── [handlers for the videostream subtype]
│   └── [route definitions for the api]
│
├── build
│   └── [Dockerfiles ans build scripts]
│
├── cmd
│   ├── seed
│   │   └── [entrypoint for the dbinit seed application]
│   └── server
│       └── [entrypoint for the server application]
│
├── deploy
│   ├── migrations
│   │   └── [migration scripts]
│   └── [env and docker-compose files]
│
├── internal
│   ├── config
│   │   └── [internal aplication config]
│   └── model
│       ├── postgres
│       │   └── [postgres backed store]
│       ├── testmodel
│       │   └── [mock store for testing]
│       └── [abstract data-model]
│
├── go.mod
├── go.sum
└── README.md
```
