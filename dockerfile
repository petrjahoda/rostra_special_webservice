# FROM alpine:latest
# RUN apk update && apk upgrade && apk add bash && apk add procps && apk add nano
# RUN apk add tzdata
# RUN rm -rf /var/cache/apk/*
# RUN cp /usr/share/zoneinfo/Europe/Prague /etc/localtime
# WORKDIR /bin
# COPY /css /bin/css
# COPY /html /bin/html
# COPY /js /bin/js
# COPY /mif /bin/mif
# COPY /linux /bin
# ENTRYPOINT rostra_special_webservice_linux
FROM alpine:latest as build
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Prague /etc/localtime

FROM scratch as final
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/localtime /etc/localtime
COPY /css /css
COPY /html html
COPY /js js
COPY /mif mif
COPY /linux /
CMD ["/rostra_special_webservice_linux"]