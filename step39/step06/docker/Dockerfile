FROM --platform=linux/amd64 ubuntu:18.04
# 设置工作目录
WORKDIR /app
COPY main  .
COPY ./templates ./templates
EXPOSE 8081
CMD ["./main"]