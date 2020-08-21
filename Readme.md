# Nano NTP
An NTP server implementation in golang base on [btfak](https://github.com/btfak/sntp) and [beevik](https://github.com/beevik/ntp)  
This NTP Server first acts like a ntp relay and ask time from upstream NTP servers, if relay can't reach them, then will pass it's own local time to the clients 

### Installation
Installation for linux, runing only on linux. ([TinyCore Linux 32-bit](http://tinycorelinux.net/11.x/x86/release) Is Recommended)

```
# go get https://github.com/sina-ghaderi/nanontp.git
# cd nanontp
# GOOS=linux go build -o ntp-amd64  # 64-Bit Build
# GOOS=linux GOARCH=386 go build -o ntp-i386  # 32-Bit Build

```

### Usage and Options
```
usage of snix ntp server:
./ntp-server -net [ipv4:port] ntp-domain.com:port ntp-domain.org:port ...

options:
  --net string     udp network to listen on <ipv4:port> (default "0.0.0.0:123")
  --h              print this banner and exit
example: 
  ./nanontp --net 0.0.0.0:123 time.google.com:123 ntp.day.ir:123 10.10.10.10:123 

Copyright (c) 2020 slc.snix.ir, All rights reserved.
Developed BY a.esmaeilpour@irisaco.com And s.ghaderi1999@gmail.com
This work is licensed under the terms of the MIT license.

```

### Runing Nano NTP
```
# ./ntp-amd64 -net 0.0.0.0:123 ntp.day.ir:123 132.163.96.5:123 129.6.15.27:123
2020/08/21 20:41:31 ntp server listening on (UDP) 0.0.0.0:123
------------------ Logs and Errors ------------------
2020/08/21 20:44:20 request ---> asking for time from 127.0.0.1:60924
2020/08/21 20:44:20 access ----> trying to ntp server: ntp.day.ir:123
2020/08/21 20:44:21 success ---> time received from: ntp.day.ir:123
2020/08/21 20:44:21 response --> answering to the client 127.0.0.1:60924

```
