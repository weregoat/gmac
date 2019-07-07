# Standard golang image
FROM golang

RUN mkdir -p /go/src/gmac
ADD . /go/src/gmac
WORKDIR /go/src/gmac

RUN go build -o /go/bin/gmac

# Alternative (smaller) builds
# I find kind of insane that the resulting image would be around 700MB.
# You could use another smaller image like the golang:alpine (ca. 350MB)
# Here is an alternative that is good enough for me:
#
#FROM golang:alpine AS builder
#WORKDIR $GOPATH/src/gmac/
#COPY . .
#RUN go build -o /go/bin/gmac
#FROM scratch
#COPY --from=builder /go/bin/gmac /go/bin/gmac
#
# From https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324


ENTRYPOINT ["/go/bin/gmac"]
