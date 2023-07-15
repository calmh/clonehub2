FROM alpine/git:latest
ARG TARGETARCH
COPY clonehub-linux-${TARGETARCH} /bin/clonehub
ENTRYPOINT ["/bin/clonehub"]
