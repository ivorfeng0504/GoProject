FROM dockerhub.emoney.cn/emoney/appbase:alpine
MAINTAINER emoney-panxinming

# 创建基础目录
RUN mkdir -p /emoney/softweb/roboadvisor/web_api_stocktalk/innerlogs
RUN mkdir -p /emoney/softweb/roboadvisor/web_api_stocktalk/logs

# 拷贝运行文件及配置文件
COPY config /emoney/softweb/roboadvisor/web_api_stocktalk/conf
COPY cmd /emoney/softweb/roboadvisor/web_api_stocktalk/roboadvisor.web_api_stocktalk

# 拷贝启动脚本
COPY docker-entrypoint.sh /emoney/docker-entrypoint.sh

# 赋权
RUN chmod +x /emoney/softweb/roboadvisor/web_api_stocktalk/roboadvisor.web_api_stocktalk
RUN chmod +x /emoney/docker-entrypoint.sh

ENTRYPOINT ["/emoney/docker-entrypoint.sh"]