#!/usr/bin/env sh

docker run -it --rm -v "$(pwd):/src" -u "$(id -u):$(id -g)" --network host --workdir /src/webui node:20 /bin/bash

#docker run -it --rm -v "$(pwd):/src" -u "$(id -u):$(id -g)" --workdir /src/webui -p 5173:5173 node:20 yarn run dev --host

