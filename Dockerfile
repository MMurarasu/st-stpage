FROM golang:1.9-alpine as builder
WORKDIR /build
RUN apk update && apk add git
RUN git clone https://github.com/aso930/st-stpage.git
WORKDIR /build/st-stpage
ENV GOPATH=/build
ENV GOBIN=$GOPATH/bin
RUN go get .
RUN go build .

FROM alpine:latest
COPY --from=builder /build/st-stpage/st-stpage /usr/bin/st-stpage
ENTRYPOINT ["st-stpage", "-agent", "agent:18080"]