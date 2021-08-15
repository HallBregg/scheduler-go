FROM golang:latest AS build
ENV GO111MODULE=off
WORKDIR /go/src/scheduler
COPY . .
RUN go build -o /out/bapp

FROM scratch as main
COPY --from=build /out/bapp /
CMD ["/bapp"]

