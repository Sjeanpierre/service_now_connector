FROM scratch
MAINTAINER Stevenson Jean-Pierre <stevenson.jean-pierre@sage.com>
ADD bin/linux_service_now_proxy service_now_proxy
EXPOSE 8080
ENTRYPOINT ["/service_now_proxy"]