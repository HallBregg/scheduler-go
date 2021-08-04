FROM golang:latest AS build
WORKDIR /src
COPY . .
RUN go build -o /out/bapp

FROM scratch as main
COPY --from=build /out/bapp /
CMD ["/bapp"]

