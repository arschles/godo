FROM ubuntu-debootstrap:14.04

COPY ./godo godo
COPY ./godo.yaml godo.yaml
RUN mv godo /bin && mv godo.yaml /bin

CMD godo
