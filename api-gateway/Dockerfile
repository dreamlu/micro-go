FROM alpine
MAINTAINER dreamlu <862362681@qq.com>

WORKDIR /app
#COPY conf \
#     static \
#     yyd \
#     /root/yyd/
#COPY conf /app/conf
COPY main /app/
#COPY static /app/static

EXPOSE 9000

#CMD ["./main"]
#CMD ["./main --api_address=0.0.0.0:9000 --registry=consul --api_namespace micro-go.web api --handler web"]