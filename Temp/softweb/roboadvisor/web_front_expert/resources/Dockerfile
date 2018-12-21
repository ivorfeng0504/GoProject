FROM dockerhub.emoney.cn/emoney/appbase:alpine
MAINTAINER emoney-panxinming

# 创建基础目录
RUN mkdir -p /emoney/softweb/roboadvisor/web_front_expert/innerlogs
RUN mkdir -p /emoney/softweb/roboadvisor/web_front_expert/logs

# 拷贝运行文件及配置文件
COPY config /emoney/softweb/roboadvisor/web_front_expert/conf
COPY views /emoney/softweb/roboadvisor/web_front_expert/views
COPY static /emoney/softweb/roboadvisor/web_front_expert/static
COPY templates /emoney/softweb/roboadvisor/web_front_expert/templates
COPY cmd /emoney/softweb/roboadvisor/web_front_expert/roboadvisor.web_front_expert

# 拷贝启动脚本
COPY docker-entrypoint.sh /emoney/docker-entrypoint.sh

# 赋权
RUN chmod +x /emoney/softweb/roboadvisor/web_front_expert/roboadvisor.web_front_expert
RUN chmod +x /emoney/docker-entrypoint.sh

ENTRYPOINT ["/emoney/docker-entrypoint.sh"]