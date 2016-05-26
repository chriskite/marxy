FROM phusion/baseimage:0.9.18

# Set correct environment variables.
ENV HOME /root

# Install Haproxy.
RUN \
  sed -i 's/^# \(.*-backports\s\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get install -y haproxy=1.5.14-1ubuntu0.15.10.1~ubuntu14.04.1 wget git && \
  sed -i 's/^ENABLED=.*/ENABLED=1/' /etc/default/haproxy && \
  rm -rf /var/lib/apt/lists/*

# install and setup go
RUN wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
RUN tar -C /usr/local -zxf go1.6.linux-amd64.tar.gz
RUN mkdir /go
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# build marxy
RUN go get github.com/tools/godep
ADD . /go/src/github.com/chriskite/marxy
WORKDIR /go/src/github.com/chriskite/marxy
RUN godep go install

# setup marxy service
RUN mkdir /etc/service/marxy
ADD run /etc/service/marxy/run

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Add files.
ADD haproxy.cfg /etc/haproxy/haproxy.cfg
ADD haproxy.sh /bin/

# Use baseimage-docker's init system.
CMD ["/sbin/my_init", "--", "/bin/haproxy.sh"]

# Expose ports.
EXPOSE 10000 10001 10002 10003 10004 10005 10006 10007 10008 10009 10010 10011 10012 10013 10014 10015 10016 10017 10018 10019 10020 10021 10022 10023 10024 10025 10026 10027 10028 10029 10030 10031 10032 10033 10034 10035 10036 10037 10038 10039 10040 10041 10042 10043 10044 10045 10046 10047 10048 10049 10050 10051 10052 10053 10054 10055 10056 10057 10058 10059 10060 10061 10062 10063 10064 10065 10066 10067 10068 10069 10070 10071 10072 10073 10074 10075 10076 10077 10078 10079 10080 10081 10082 10083 10084 10085 10086 10087 10088 10089 10090 10091 10092 10093 10094 10095 10096 10097 10098 10099
EXPOSE 9090
