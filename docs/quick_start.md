## Falcon Qucik Start

## overview
![](img/falcon-overview.svg?raw=true)

## install && start

```
#install gcc make automak libtool golang ...

# install protoc
wget https://github.com/google/protobuf/archive/v3.4.1.tar.gz
tar xzvf v3.4.1.tar.gz
cd protobuf-3.4.1
./autogen.sh
./configure
make
sudo make install

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/beego/bee

# build falcon
git clone https://github.com/yubo/falcon
cd falcon
make
make install

## nginx
sudo apt-get install nginx
mv /etc/nginx/conf.d/falcon.conf.example /etc/nginx/conf.d/falcon.conf
#edit /etc/nginx/conf.d/falcon.conf
#edit /etc/nginx/mime.types
  -    text/html                             html htm shtml;
  +    text/html                             html htm shtml md;

## mysql
sudp apt-get isntall mysql
mysql -u root -p < ./scripts/db_schema/01_database.sql
mysql -u root -p < ./scripts/db_schema/02_database_user.sql
mysql -u falcon -p1234 falcon < ./scripts/db_schema/03_falcon.sql
mysql -u falcon -p1234 alarm < ./scripts/db_schema/04_alarm.sql
mysql -u falcon -p1234 idx < ./scripts/db_schema/05_index.sql

## start falcon
cp /etc/falcon/falcon.example.conf /etc/falcon/falcon.conf
sudo service falcon start
```

## API flow
![](img/falcon-api.svg?raw=true)

## trigger

#### event
![](img/falcon-event.svg?raw=true)

#### action
![](img/falcon-action.svg?raw=true)



## benchmark

```
cd backend
go test -bench=Add -benchtime=20s
go test -bench=.
```


## file/dir list

dir 			| desc
--  			| --
/etc/falcon/falcon.conf	| config
/etc/init.d/falcon	| init.d script
/opt/falcon/log		| log
/opt/falcon/tsdb	| tsdb storage directry
/opt/falcon/emu_tpl	| emulator template file directry
/opt/falcon/html	| nginx document root
/sbin/falcon		| falcon binary executable file(all module)
/sbin/agent		| falcon-agent binary executable file(just single module)


