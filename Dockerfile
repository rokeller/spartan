FROM golang:1-trixie AS build

WORKDIR /src
COPY go.* ./
RUN go mod download

COPY . .
# Create an optimized statically linked binary so we can create an image from scratch below.
RUN CGO_ENABLED=0 go build -ldflags '-s -w'

FROM scratch

WORKDIR /srv
USER 1000
ENTRYPOINT [ "/srv/spartan" ]
EXPOSE 8080/tcp

COPY --from=build /src/spartan /srv
