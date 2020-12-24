FROM golang:1.15.6 AS builder
WORKDIR /go/src/github.com/joberly/aws-ssm-env
COPY ./ ./
RUN sh -eux; \
    go install;

FROM debian:buster
COPY --from=builder /go/bin/aws-ssm-env ./
CMD ["./aws-ssm-env"]
