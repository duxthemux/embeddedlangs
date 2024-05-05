FROM alpine
ARG TARGETARCH
RUN apk add git
ENV ARCH $TARGETARCH
RUN echo "Building for $ARCH"
WORKDIR /data
COPY ./bin/main_amd64 /srv/main_amd64
COPY ./bin/main_arm64 /srv/main_arm64
CMD [ "sh" ,"-c", "/srv/main_${ARCH}" ]