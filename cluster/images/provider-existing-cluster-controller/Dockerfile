FROM BASEIMAGE
RUN apk --no-cache add ca-certificates bash

ADD provider /usr/local/bin/crossplane-existing-cluster-provider

EXPOSE 8080
USER 1001
ENTRYPOINT ["crossplane-existing-cluster-provider"]
