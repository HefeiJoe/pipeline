FROM golang:latest as build-env
WORKDIR /go/src/EnSaaS_Pipeline_Backend
ADD . /go/src/EnSaaS_Pipeline_Backend
RUN go build -o /go/main

FROM harbor.wise-paas.io/distroless/base:latest as prod-env
WORKDIR /go/
COPY --from=build-env /go/src/EnSaaS_Pipeline_Backend/resources/ resources
COPY --from=build-env /go/main .
EXPOSE 8080
CMD ["./main"]