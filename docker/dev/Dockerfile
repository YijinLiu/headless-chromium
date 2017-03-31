FROM ubuntu:16.04
ARG version

RUN echo America/Los_Angeles > /etc/timezone
RUN dpkg-reconfigure --frontend noninteractive tzdata

RUN apt update
RUN apt install -y --no-install-recommends bash build-essential dpkg git golang-go python \
        software-properties-common ssh sudo vim wget

RUN useradd hcdev --shell /bin/bash --groups sudo
RUN echo 'hcdev ALL = NOPASSWD: ALL' > /etc/sudoers.d/hcdev
RUN addgroup --gid 999 docker
RUN usermod -G docker hcdev

COPY docker_tmp_data /home/hcdev/external
RUN chown -R hcdev:hcdev /home/hcdev
USER hcdev
WORKDIR /home/hcdev/
RUN ./external/scripts/install_chromium.sh --version=$version
RUN ./external/scripts/setup.sh
RUN rm -rf external
