FROM golang:1.11.0-stretch

WORKDIR /go/src/github.com/rekkusu/gyotaku
COPY . /go/src/github.com/rekkusu/gyotaku

RUN apt-get update && apt-get install -y supervisor chromium git curl apache2 apache2-utils libapache2-mod-security2 \
  && rm -rf /etc/modsecurity/crs \
  && git clone https://github.com/SpiderLabs/owasp-modsecurity-crs.git /etc/modsecurity/crs \
  && go get -u github.com/golang/dep/... \
  && dep ensure \
  && go build \
  && rm -rf /go/src/github.com/rekkusu/gyotaku/app /go/src/github.com/rekkusu/gyotaku/main.go

ENV APACHE_RUN_USER www-data
ENV APACHE_RUN_GROUP www-data
ENV APACHE_PID_FILE /var/run/apache2/apache2.pid
ENV APACHE_RUN_DIR /var/run/apache2
ENV APACHE_LOCK_DIR /var/lock/apache2
ENV APACHE_LOG_DIR /var/log/apache2
ENV CHROME /usr/bin/chromium
ENV LISTEN 0.0.0.0:9999

COPY ./deploy/supervisord.conf /usr/etc/supervisord.conf
COPY ./deploy/security2.conf /etc/apache2/mods-available/security2.conf
COPY ./deploy/crs-setup.conf /etc/modsecurity/crs/crs-setup.conf
COPY ./deploy/modsecurity.conf /etc/modsecurity/modsecurity.conf
COPY ./deploy/apache2.conf /etc/apache2/apache2.conf

RUN a2enmod security2 proxy proxy_http \
  && a2dissite 000-default \
  && mkdir -p /var/lock/apache2 /var/run/apache2

EXPOSE 80

CMD ["/usr/bin/supervisord"]
