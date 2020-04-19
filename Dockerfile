FROM ubuntu:18.04
WORKDIR /app/
ADD ./release/pstree-json-latest-linux-amd64 /app/pstree-json
CMD ["/bin/bash", "-c", "/app/pstree-json", "$PID"]
