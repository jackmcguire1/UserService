FROM golang:1.16.3-stretch

WORKDIR '/app/'

EXPOSE 8000

CMD [ "/app/entrypoint.sh" ]
