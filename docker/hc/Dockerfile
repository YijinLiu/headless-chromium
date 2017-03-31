FROM ubuntu:16.04

RUN echo America/Los_Angeles > /etc/timezone
RUN dpkg-reconfigure --frontend noninteractive tzdata

RUN apt update
RUN apt install -y --no-install-recommends bash build-essential dpkg git golang-go \
        libfontconfig1 libnss3 python software-properties-common ssh sudo vim

RUN useradd hcdev --shell /bin/bash --groups sudo
RUN echo 'hcdev ALL = NOPASSWD: ALL' > /etc/sudoers.d/hcdev

COPY docker_tmp_data /home/hcdev/build
RUN cp -av /home/hcdev/build/headless_chromium /usr/local/
RUN chown -R hcdev:hcdev /home/hcdev
USER hcdev
WORKDIR /home/hcdev/build
RUN ./scripts/setup.sh
WORKDIR /home/hcdev
RUN rm -rf build
