FROM metalstack/builder:latest as builder
RUN make platforms

FROM alpine:3.13
LABEL maintainer="metal-stack Authors <info@metal-stack.io>"
COPY --from=builder /work/bin/cloudctl /cloudctl
ENTRYPOINT ["/cloudctl"]
