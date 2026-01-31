FROM python:3.12-slim-bookworm AS chroot

RUN apt update -y
RUN apt install -y curl
RUN apt install -y build-essential

RUN useradd --create-home -u 1000 user
USER user

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y --profile=minimal --default-toolchain=nightly

FROM disconnect3d/nsjail:3.1-6483728

COPY --from=chroot / /chroot

RUN apt update -y
RUN apt install -y socat python3

RUN useradd --create-home -u 1000 user

COPY . /home/user/
WORKDIR /home/user/

CMD socat TCP-LISTEN:1337,reuseaddr,fork EXEC:"/usr/bin/python3 /home/user/runner.py",pty,stderr,icanon=0,echo=0
