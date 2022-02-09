FROM alpine:3.13
COPY clonehub-linux-amd64 /bin/clonehub
ENTRYPOINT ["/bin/clonehub"]
