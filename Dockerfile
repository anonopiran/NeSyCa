FROM golang:1.19.3-alpine as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w"  -o /NeSyCa

FROM scratch as app
COPY --from=build /NeSyCa /NeSyCa
ENTRYPOINT [ "/NeSyCa" ]