FROM alpine:3.6
MAINTAINER danil@nulana.com

COPY battleship /internal/bin/battleship-server
COPY entrypoint.sh /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/internal/bin/battleship-server"]