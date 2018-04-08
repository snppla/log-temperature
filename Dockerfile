from debian as build-rtl_433
RUN apt-get update && apt-get upgrade -y 
RUN apt-get install -y libtool libusb-1.0.0-dev librtlsdr-dev rtl-sdr build-essential autoconf cmake pkg-config git && apt-get clean
RUN git clone https://github.com/merbanan/rtl_433.git
WORKDIR /rtl_433
RUN mkdir -p build
WORKDIR /rtl_433/build
RUN cmake -DCMAKE_INSTALL_PREFIX=/rtl_433/install ../
RUN make -j $(nproc) install

from debian as build-logger

RUN apt-get update && apt-get upgrade -y

#############################

from debian as build-log
RUN apt-get update && apt-get upgrade -y && apt-get clean
RUN apt-get install -y golang git
RUN mkdir -p /go/src/log-temperature
ADD log.go /go/src/log-temperature
WORKDIR /go/src/log-temperature
ENV GOPATH /go
ENV GOBIN /go/bin
RUN go get
RUN go install

##############################
from debian as main
RUN apt-get update && apt-get upgrade -y 
RUN apt-get install -y librtlsdr0 && apt-get clean
COPY --from=build-rtl_433 /rtl_433/install/bin/rtl_433 /usr/bin
COPY --from=build-log /go/bin/log-temperature /usr/bin

ENV DB temperature
ENV USER ""
ENV PASSWORD ""
ENV HOST http://localhost:8086
ENV RTL_ARGS ""
ENV LOCATION room

CMD rtl_433 $RTL_ARGS | log-temperature -host $HOST -db $DB -password $PASSWORD -user $USER -location $LOCATION

##############################


