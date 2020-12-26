FROM registry.cn-hangzhou.aliyuncs.com/dreamlu/alpine
MAINTAINER dreamlu <862362681@qq.com>

WORKDIR /
COPY conf /conf
COPY main /
#COPY static /app/static

EXPOSE 9000

#CMD ["./main"]
#CMD ["./main --api_address=0.0.0.0:9000 --registry=consul --api_namespace demo.web api --handler web"]