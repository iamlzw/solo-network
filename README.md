fabric v1.4.2示例网络,solo共识版本

### 1、clone repository

```bash
$ cd ANY PATH
$ git clone https://github.com/iamlzw/fabric-network-samples.git
```
### 2、copy chaincode to GOPATH
```bash
$ cd fabric-network-samples
$ cp -r chaincode $GOPATH/src/
```

### 3、start fabric network

```bash
$ cd fabric-network-samples/v1.4.2/solo
$ ./start.sh
```

### 4、init 
this step include 
- create channel
- join channel
- install chaincode
- instatinate chaincode
- invoke chaincode

```bash
$ ./init.sh
```

### 5、stop and clean network
```bash
$ ./stop.sh
```
