FROM alpine:3.6
MAINTAINER danil@nulana.com

COPY battleship /internal/bin/battleship_server
COPY entrypoint.sh /entrypoint.sh

RUN ls -l
RUN ls -l /internal/bin/

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/internal/bin/battleship_server"]