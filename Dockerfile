FROM ubuntu-debootstrap:14.04

COPY ./gci gci
RUN mv gci /bin

CMD gci
