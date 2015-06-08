FROM debian:jessie

RUN mkdir -p /etc/ssl/certs /app

COPY run.sh /app/chat.sh
RUN chmod +x /app/chat.sh
COPY components /app/components
COPY bower_components /app/bower_components
COPY templates /app/templates
COPY chat /app/chat
RUN chmod +x /app/chat


ENV DOMAIN "localhost"
ENV PORT "8080"

EXPOSE 8080

WORKDIR /app/
CMD ["/app/chat.sh"]