# Simple Operator
A simple operator to handle your resource.

## Build and run
To build from source using Make tool:
```sh
$ git clone https://github.com/kairen/simple-operator.git $GOPATH/src/github.com/kairen/simple-operator
$ cd $GOPATH/src/github.com/kairen/simple-operator
$ make
```

To run operator:
```sh
$ vagrant up
$ vagrant ssh operator
$ sudo -i
$ cd go/src/github.com/kairen/simple-operator/
$ source /root/tokenrc
$ go run cmd/main.go serve --logtostderr=true --endpoint=https://192.16.35.10:6443 --token=${TOKEN} -v=2
```
