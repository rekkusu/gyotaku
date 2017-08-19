FROM ubuntu:16.04

RUN apt-get update \
  && apt-get install -y git apache2 apache2-utils libapache2-mod-security2 \
  && git clone https://github.com/SpiderLabs/owasp-modsecurity-crs.git /etc/modsecurity/crs


#COPY ./httpd.conf /usr/local/apache2/conf/httpd.conf

ENV APACHE_RUN_USER www-data
ENV APACHE_RUN_GROUP www-data
ENV APACHE_PID_FILE /var/run/apache2/apache2.pid
ENV APACHE_RUN_DIR /var/run/apache2
ENV APACHE_LOCK_DIR /var/lock/apache2
ENV APACHE_LOG_DIR /var/log/apache2

COPY modsecurity.conf /etc/modsecurity/modsecurity.conf
COPY crs-setup.conf /etc/modsecurity/crs/crs-setup.conf
COPY apache2.conf /etc/apache2/apache2.conf

RUN a2enmod security2 proxy proxy_http
RUN a2dissite 000-default

EXPOSE 80
CMD ["/usr/sbin/apache2", "-D", "FOREGROUND"]
