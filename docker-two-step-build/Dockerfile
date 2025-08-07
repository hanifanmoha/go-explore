##
## Build
##

FROM public.ecr.aws/docker/library/golang:1.23-bullseye AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN go build -o cmd/main main.go

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/cmd/main /timed-exam-api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/timed-exam-api"]
