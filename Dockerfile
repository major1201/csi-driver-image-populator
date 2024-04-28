FROM major1201/debian:bookworm-runit

RUN \
  apt update && \
  apt install -y buildah && \
  apt clean

COPY ./bin/imagepopulatorplugin /imagepopulatorplugin
ENTRYPOINT ["/imagepopulatorplugin"]
