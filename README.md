## watch command

> watch [options] command [command options]


build
```
go build -o watch

GOOS=windows GOARCH=386 go build -o watch.exe
```

use demo
```
watch.exe -n 3 curl 'sftp://root:123@192.168.149.128/root/test/[0-10].txt' -sO
```