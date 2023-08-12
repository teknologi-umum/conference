FROM gomicro/goose

ADD *.sql /migrations/
ADD entrypoint.sh /migrations/

ENTRYPOINT ["/migrations/entrypoint.sh"]