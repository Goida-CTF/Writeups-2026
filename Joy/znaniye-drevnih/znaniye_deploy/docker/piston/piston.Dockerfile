FROM alpine:3.22.2 AS pkgs

RUN apk add --no-cache curl tar

WORKDIR /tmp
RUN mkdir -p /tmp/piston/packages/ \
    && curl -fsSL \
        -o gcc-10.2.0.pkg.tar.gz \
        https://github.com/engineer-man/piston/releases/download/pkgs/gcc-10.2.0.pkg.tar.gz \
    && mkdir -p /tmp/piston/packages/gcc/10.2.0 \
    && tar -xzf gcc-10.2.0.pkg.tar.gz -C /tmp/piston/packages/gcc/10.2.0 \
    && touch /tmp/piston/packages/gcc/10.2.0/.ppman-installed \
    && cat > /tmp/piston/packages/gcc/10.2.0/.env <<'EOF'
LD_LIBRARY_PATH=/piston/packages/gcc/10.2.0/lib:/piston/packages/gcc/10.2.0/lib64
PATH=/piston/packages/gcc/10.2.0/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin:.
EOF

# This is a fix custom libraries include
RUN sed -i "/rename 's\/\\$\/\\\\.cpp\//d" /tmp/piston/packages/gcc/10.2.0/compile

FROM ghcr.io/engineer-man/piston@sha256:2f66b7456189c4d713aa986d98eccd0b6ee16d26c7ec5f21b30e942756fd127a

COPY --from=pkgs --chown=piston:piston /tmp/piston/packages/ /piston/packages/
