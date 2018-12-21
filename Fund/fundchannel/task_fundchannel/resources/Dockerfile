FROM dockerhub.emoney.cn/emoney/appbase:alpine
MAINTAINER emoney-panxinming

# 创建基础目录
RUN mkdir -p /emoney/fundchannel/task_fundchannel/innerlogs
RUN mkdir -p /emoney/fundchannel/task_fundchannel/logs

# 拷贝运行文件及配置文件
COPY config /emoney/fundchannel/task_fundchannel/conf
COPY views /emoney/fundchannel/task_fundchannel/views
COPY static /emoney/fundchannel/task_fundchannel/static
COPY cmd /emoney/fundchannel/task_fundchannel/fundchannel.task_fundchannel
COPY templates /emoney/fundchannel/task_fundchannel/templates

# 拷贝启动脚本
COPY docker-entrypoint.sh /emoney/docker-entrypoint.sh

# 赋权
RUN chmod +x /emoney/fundchannel/task_fundchannel/fundchannel.task_fundchannel
RUN chmod +x /emoney/docker-entrypoint.sh

ENTRYPOINT ["/emoney/docker-entrypoint.sh"]