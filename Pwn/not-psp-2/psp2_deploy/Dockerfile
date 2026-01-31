FROM debian:bookworm-slim AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    ninja-build \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
RUN git clone https://github.com/maximxlss/melonDS-headless . && git reset --hard 9a7ec82

RUN cmake -S . -B build -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DBUILD_QT_SDL=OFF \
    -DBUILD_HEADLESS=ON \
    -DENABLE_OGLRENDERER=OFF \
    && cmake --build build --target melonDS-headless

FROM python:3.12-slim-bookworm AS runtime
COPY --from=builder /src/build/melonDS-headless /usr/local/bin/melonDS-headless
COPY requirements.txt ./
RUN pip install -r requirements.txt

COPY task.py ./

EXPOSE 8000

CMD ["python3", "task.py"]