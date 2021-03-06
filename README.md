[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/tsaridas/salt-golang) 

# Saltstack event-listener, salt cli client, salt-minion, salt-master, salt-api Proof of Concept in golang.

## Disclaimer
This is me trying to learn golang and had the idea to write a salt client and ended up writing the rest.. This is not something that is intended for production usage but rather as a Proof of Concept. 

## Runnables
#### Event Listener
A Saltstack event listener listening for event from the unix socket. Needs to be ran a on salt master server.

eventlistener/event-listener.go:
```
[root@salt-master salt-golang]#go run eventlistener/event-listener.go
Tag is 15511887253592021111 and ret is map[_stamp:2019-02-26T13:45:25.361952 minions:[salt-minion-01]]

Tag is salt/job/15511887253592021111/new and ret is map[jid:15511887253592021111 user:root tgt:[salt-minion-01] arg:[] fun:test.ping tgt_type:list missing:[] _stamp:2019-02-26T13:45:25.363199 minions:[salt-minion-01]]

Tag is salt/job/15511887253592021111/ret/salt-minion-01 and ret is map[retcode:%!s(int8=0) success:%!s(bool=true) id:salt-minion-01 fun_args:[] jid:15511887253592021111 return:%!s(bool=true) cmd:_return _stamp:2019-02-26T13:45:25.495549 fun:test.ping]
```
```
[root@salt-master salt-golang]# salt -L salt-minion-01 test.ping
salt-minion-01:
    True
```

#### Salt Master
A Saltstack Master listens to ports 4505,4506 and pushes events to a unix socket. It will accept auth requests and publish commands from salt-cli (Official or not). It does not support all the fuctions that a normal Salt Master does like saltfs etc.

salt-master/salt-master.go:
```
[root@salt-master salt-golang]#go run salt-master/salt-master.go
2020/05/08 13:01:51 Generated AES key.
2020/05/08 13:01:51 Started Router on port 4506.
2020/05/08 13:01:51 Started Publisher on port 4505.
2020/05/08 13:01:51 Loaded master public key.
2020/05/08 13:01:51 Loaded master public key.
2020/05/08 13:01:51 Loaded master root key.
2020/05/08 13:01:51 Starting IPC server...
2020/05/08 13:01:51 Started 1 workers.
2020/05/08 13:01:51 Starting proxy
2020/05/08 15:40:19 Received an authentication event from: salt-minion-02
2020/05/08 15:40:19 Accepted connection from minion salt-minion-02.
2020/05/08 15:40:28 Received new ipc connection. map[0xc00000e280:true]
2020/05/08 15:40:28 Received a publish event:{Enc:clear Load:map[arg:[] cmd:publish fun:test.ping jid:15889524281431282091 key:8f46wPP20/f1Ojel3KgwBsJfReDrUXuxcc0/zYlypVkOK5JSHgOr+LTBW+1BiZD55aWhRrUy51A=1
 kwargs:map[delimiter:: show_jid:false show_timeout:true] ret: tgt:[salt-minion-02] tgt_type:list user:root] Token:<nil>}
2020/05/08 15:40:28 Received ret is map[cmd:_return fun:test.ping fun_args:[] id:salt-minion-02 jid:15889524281431282091 retcode:0 return:True success:true]
```

#### Client
A Saltstack client to sent commands from the master directly to the zmq port 4506 written in go. Will run a salt module to a minion. Needs to be ran from a salt master server.

salt-cli/salt.go: 
```
[root@salt-master salt-golang]# go run salt-cli/salt.go -L minion test.ping
salt-minion-01:
   True
[root@salt-master salt-golang]# go run salt-cli/salt.go -L salt-minion-01 cmd.run 'ls'
salt-minion-01:
anaconda-ks.cfg
original-ks.cfg
```
As you can see the results are not formated the same way that the normal salt client does ( adds some spaces infront of the result on each line ).

#### Client-v2
A new version of Saltstack client which uses the libraries for client and listener to sent commands from the master directly to the zmq port 4506 and listen to the event bus from the unix socket written in go. Will run a salt module to a minion. Needs to be ran from a salt master server.

salt-cli-v2/salt.go: 
```
[root@salt-master salt-golang]# go run salt-cli-v2/salt.go -L salt-minion-01,salt-minion-02,salt-minion-03 test.ping
salt-minion-02:
    True
salt-minion-01:
    True
salt-minion-03:
    False
```

#### API
An HTTP api that recevies a request sends them to the minion and gives the response. Will initialize a http api and ping the minion-id that you define. If minion does not return it will timeout and return the JID. Needs to be ran on a Salt Master server. The API will wait 5 seconds to receive results and it will return false.

api/salt-api.go: 
```
[root@salt-master salt-golang]# go run api/salt-api.go
[root@salt-master salt-golang]# go run salt-api.go
2019/02/27 10:44:06 Sending command to: salt-minion-01 .
2019/02/27 10:44:06 Got result from: salt-minion-01 .
2019/02/27 10:44:08 Sending command to: salt-minion-02 .
2019/02/27 10:44:13 Timeout 15512642482554418571
```
```
[root@salt-master salt-golang]# time curl 127.0.0.1:8080/minion-id
true

real	0m0.125s
user	0m0.001s
sys	0m0.007s
[root@salt-master salt-golang]# time curl 127.0.0.1:8080/salt-minion-02
false

real	0m5.008s
user	0m0.001s
sys	0m0.006s
```

#### Salt-minion
A Salstack minion that connects to the saltstack master and receives events. You need to use target type list or glob in order to target the minion and make sure the normal salt-minion is installed and has the keys generated. 
salt-minion/salt-minion.go : Will start a subscriber that connects to a salt master server. 

- If you do not define --id it will try to load minion id from the configuration or generate it from hostname.
- If you do not define --masterip it will try to load it from the configuration.

salt-minion/salt-minion.go:
```
[root@salt-master salt-golang]# go run salt-minion/salt-minion.go -id salt-minion-01 --masterip 192.168.1.1
2019/02/28 01:42:25 Loading config file: /etc/salt/minion
2019/02/28 01:42:25 Loading config file: /etc/salt/minion.d/minion.conf
2019/02/28 01:42:25 Using configured master ip : 192.168.1.1
2019/02/28 01:42:25 Using configured minion id : salt-minion-01
2019/02/28 01:42:25 Authenticated with Master.
2019/02/28 01:42:25 Subscribed to Master.
2019/02/28 01:42:29 Got function : test.ping with event map[tgt_type:list jid:20190228014229414793 tgt:[salt-minion-01] ret: user:sudo_vagrant arg:[] fun:test.ping]
2019/02/28 01:42:29 Replied to event : map[user:sudo_vagrant arg:[] fun:test.ping tgt_type:list jid:20190228014229414793 tgt:[salt-minion-01] ret:]
```
The minion will only answer to test.ping requests.

I added some functionality to be able to load Plugins which its really useless but it was nice to play with. You can build the existing plugins with :

```
[root@salt-master salt-minion]# go build -ldflags="-s -w" -buildmode=plugin -o modules/saltutil.so modules/saltutil.go 
[root@salt-master salt-minion]# go build -ldflags="-s -w" -buildmode=plugin -o modules/cmd.so modules/cmd.go 
[root@salt-master salt-minion]# go build -ldflags="-s -w" -buildmode=plugin -o modules/test.so modules/test.go 
```
and you will be able to run test.ping, cmd.run with a parameter a command and saltutil.find_job is just dummy.
 

## Requirements
```
[root@salt-master]# yum install zeromq-devel -y
```

## Installation
```
[root@salt-master]# go get github.com/tsaridas/salt-golang
[root@salt-master]# cd $GOPATH/src/github.com/tsaridas/salt-golang
[root@salt-master]# go get -d ./...
```

## Build
You can cd to any runnable directory and run the build
```
[root@salt-master salt-golang]# go build salt-cli/salt.go
[root@salt-master salt-golang]# ls -ltah salt
-rwxr-xr-x. 1 root root 3.3M Feb 27 09:53 salt
[root@salt-master salt-golang]# ./salt
Application Flags:
  -L string
    	Minion comma separated list of minions.
```

## Tested
This was tested on a Centos7.5 and salt-2018.3.3. It should work on lower and higher versions of Saltstack.

## Performance
With the go client
```
[root@salt-master salt-golang]# time ./salt -L salt-minion-01 test.ping
salt-minion-01:
   True

real	0m0.030s
user	0m0.004s
sys	0m0.005s
```
With the python client
```
[root@salt-master salt-golang]# time salt -L salt-minion-01 test.ping
salt-minion-01:
    True

real	0m0.909s
user	0m0.544s
sys	0m0.184s
```
## ToDo's
- Cleanup code
- Improve salt client to support multiple target types and be as close to the original client
