FROM golang:1-trixie AS build

WORKDIR /src
COPY go.* ./
RUN go mod download

COPY . .
ARG SPARTAN_TAGS=""
# Create an optimized statically linked binary so we can create an image from scratch below.
RUN CGO_ENABLED=0 go build -tags "$SPARTAN_TAGS" -ldflags '-s -w'

FROM scratch

LABEL org.opencontainers.image.source=https://github.com/rokeller/spartan
LABEL org.opencontainers.image.licenses=Apache-2.0

WORKDIR /srv
USER 1000
ENTRYPOINT [ "/srv/spartan" ]
EXPOSE 8080/tcp

COPY --link --from=build /src/spartan /srv/
