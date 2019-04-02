#!/bin/bash

docker run -it --rm --mount type=bind,source="$(pwd)",target=/go/src/dd dd
