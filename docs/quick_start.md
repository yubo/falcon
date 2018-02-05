## install && start

```
## falcon
git clone git@git.n.xiaomi.com:falcon/falcon-lite.git
cd falcon-lite
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


