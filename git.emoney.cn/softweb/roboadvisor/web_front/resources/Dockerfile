FROM dockerhub.emoney.cn/emoney/appbase:alpine
MAINTAINER emoney-panxinming

# 创建基础目录
RUN mkdir -p /emoney/softweb/roboadvisor/web_front/innerlogs
RUN mkdir -p /emoney/softweb/roboadvisor/web_front/logs

# 拷贝运行文件及配置文件
COPY config /emoney/softweb/roboadvisor/web_front/conf
COPY views /emoney/softweb/roboadvisor/web_front/views
COPY static /emoney/softweb/roboadvisor/web_front/static
COPY cmd /emoney/softweb/roboadvisor/web_front/roboadvisor.web_front
COPY templates /emoney/softweb/roboadvisor/web_front/templates

# 拷贝启动脚本
COPY docker-entrypoint.sh /emoney/docker-entrypoint.sh

# 赋权
RUN chmod +x /emoney/softweb/roboadvisor/web_front/roboadvisor.web_front
RUN chmod +x /emoney/docker-entrypoint.sh

ENTRYPOINT ["/emoney/docker-entrypoint.sh"]