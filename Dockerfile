FROM ubuntu-debootstrap:14.04

COPY ./godo godo
COPY ./gci.yaml gci.yaml
RUN mv godo /bin && mv gci.yaml /bin

CMD godo
