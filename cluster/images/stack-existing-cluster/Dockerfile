FROM BASEIMAGE
RUN apk --no-cache add ca-certificates bash

ARG ARCH
ARG TINI_VERSION

ADD stack /usr/local/bin/crossplane-existing-cluster-stack

COPY stack-package /

EXPOSE 8080
USER 1001
ENTRYPOINT ["crossplane-existing-cluster-stack", "--debug", "--sync=30s"]
