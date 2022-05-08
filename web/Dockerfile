FROM nginx:stable-alpine

COPY dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

ENV APIADDRESS = "http://127.0.0.1:8080"
RUN echo $APIADDRESS > /usr/share/nginx/html/api_address