FROM ubuntu-debootstrap:14.04

COPY ./gci-test gci
RUN mv gci /bin

CMD gci
