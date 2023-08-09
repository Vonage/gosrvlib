FROM gosrvlibexampleowner/dev_gosrvlibexample:dev

ARG HOST_USER="root"
ARG HOST_GROUP="root"

ENV HOST_OWNER=${HOST_USER}:${HOST_GROUP}

RUN mkdir /workspace

# Add only the required project resources
COPY resources /workspace/resources
COPY openapi*.yaml /workspace/
COPY Makefile /workspace
COPY RELEASE /workspace
COPY VERSION /workspace

WORKDIR /workspace
ENTRYPOINT ["/workspace/resources/test/integration/entrypoint.sh"]
HEALTHCHECK CMD go version || exit 1
