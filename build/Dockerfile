FROM alpine
MAINTAINER "SpotMax"
RUN apk add -U --no-cache ca-certificates
COPY awsfactory /spotmax/awsfactory
COPY cloudmeta.json /spotmax/
WORKDIR /spotmax/
CMD ["./awsfactory"]
