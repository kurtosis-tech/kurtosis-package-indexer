FROM --platform=linux/amd64 alpine:3.17

# We need protobut-dev to run protobuf compiler against startosis .proto files
RUN apk update && apk add --no-cache bash protobuf-dev

WORKDIR /run

COPY ./build/kurtosis-package-indexer ./

EXPOSE 9770

CMD ./kurtosis-package-indexer
