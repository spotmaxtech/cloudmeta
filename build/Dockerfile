FROM alpine
MAINTAINER "SpotMax"
RUN apk add -U --no-cache ca-certificates
COPY awsfactory /spotmax/awsfactory
WORKDIR /spotmax/
CMD ["./awsfactory"]
