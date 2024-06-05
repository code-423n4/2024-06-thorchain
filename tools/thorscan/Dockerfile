FROM python:3.12-alpine@sha256:ef097620baf1272e38264207003b0982285da3236a20ed829bf6bbf1e85fe3cb

RUN pip3 install --no-cache-dir pdm==2.9.3

COPY . /thorscan

WORKDIR /build
RUN pip3 install --no-cache-dir /thorscan

ENTRYPOINT ["thorscan"]
CMD []
