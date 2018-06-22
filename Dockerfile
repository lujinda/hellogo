FROM alpine:3.6
EXPOSE 5050
COPY hellogo .
CMD ./hellogo
