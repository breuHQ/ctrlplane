FROM cgr.dev/chainguard/go:latest as base
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod go mod download

FROM base as src

WORKDIR /src
COPY . .
RUN git status


# migrate
FROM src as build-migrate
LABEL io.quantm.artifacts.app="quantm"
LABEL io.quantm.artifacts.component="migrate"
RUN --mount=type=cache,target=/root/go/pkg/mod go build -o ./build/migrate ./cmd/jobs/migrate

FROM cgr.dev/chainguard/git:latest-glibc as migrate

COPY --from=build-migrate /src/build/migrate /bin/migrate

ENTRYPOINT [ ]
CMD ["/bin/migrate"]


# mothership
FROM src as build-mothership
LABEL io.quantm.artifacts.app="quantm"
LABEL io.quantm.artifacts.component="mothership"
RUN --mount=type=cache,target=/root/go/pkg/mod go build -o ./build/mothership ./cmd/workers/mothership

FROM cgr.dev/chainguard/git:latest-glibc as mothership

COPY --from=build-mothership /src/build/mothership /bin/

ENTRYPOINT [ ]
CMD ["/bin/mothership"]


# api
FROM src as build-api
RUN --mount=type=cache,target=/root/go/pkg/mod go build -o ./build/api ./cmd/api
LABEL io.quantm.artifacts.app="quantm"
LABEL io.quantm.artifacts.component="api"

FROM cgr.dev/chainguard/git:latest-glibc as api

COPY --from=build-api /src/build/api /bin/api

ENTRYPOINT [ ]
CMD ["/bin/api"]
