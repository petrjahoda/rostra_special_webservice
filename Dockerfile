FROM alpine:latest as build
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Prague /etc/localtime

FROM scratch as final
COPY --from=build /etc/localtime /etc/localtime
COPY /css /css
COPY /html html
COPY /js js
COPY /mif mif
COPY /linux /
CMD ["/rostra_special_webservice"]