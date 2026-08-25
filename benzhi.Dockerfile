FROM golang:1.23 AS build
WORKDIR /src
COPY . .
ENV GOPROXY=off
ENV GOSUMDB=off
RUN go build -mod=vendor -o /coldstore ./cmd/coldstore
FROM golang:1.23
WORKDIR /app
COPY . .
COPY --from=build /coldstore /app/coldstore
ENV GOPROXY=off
ENV GOSUMDB=off
CMD ["/app/coldstore"]
