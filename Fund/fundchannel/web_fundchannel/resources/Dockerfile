FROM dockerhub.emoney.cn/emoney/appbase:alpine
MAINTAINER emoney-panxinming

# 创建基础目录
RUN mkdir -p /emoney/fundchannel/web_fundchannel/innerlogs
RUN mkdir -p /emoney/fundchannel/web_fundchannel/logs

# 拷贝运行文件及配置文件
COPY config /emoney/fundchannel/web_fundchannel/conf
COPY views /emoney/fundchannel/web_fundchannel/views
COPY static /emoney/fundchannel/web_fundchannel/static
COPY cmd /emoney/fundchannel/web_fundchannel/fundchannel.web_fundchannel
COPY templates /emoney/fundchannel/web_fundchannel/templates

# 拷贝启动脚本
COPY docker-entrypoint.sh /emoney/docker-entrypoint.sh

# 赋权
RUN chmod +x /emoney/fundchannel/web_fundchannel/fundchannel.web_fundchannel
RUN chmod +x /emoney/docker-entrypoint.sh

ENTRYPOINT ["/emoney/docker-entrypoint.sh"]