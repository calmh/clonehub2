FROM alpine/git:latest
COPY clonehub-linux-amd64 /bin/clonehub
ENTRYPOINT ["/bin/clonehub"]
