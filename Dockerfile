FROM ubuntu-debootstrap:14.04

COPY ./gci gci
COPY ./gci.yaml gci.yaml
RUN mv gci /bin && mv gci.yaml /bin

CMD gci
