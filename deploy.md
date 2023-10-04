## install go
https://go.dev/doc/install
```
cd ~/tmp
```
```
wget https://go.dev/dl/...linux-amd64.tar.gz
```
```
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
```
Add go PATH variable 
```
vim ~/.profile
```
```
export PATH=$PATH:/usr/local/go/bin
```
gcc compiler (`CGO_ENABLE=1`) is required for sqlite driver
```
apt-get install build-essential
```

## Copying local code to remote
```
rsync -a . root@123.123.123.123:/home/go/src/rss --exclude .env
```

## build executables
```
go build -o bin/web web/*.go ; go build -o bin/cron cron/main.go ; go build -o bin/init db/init.go
```
Create clean database
```
./bin/init
```

## web service
```
cp rss.service /lib/systemd/system/
```
Start service
```
service rss start
service rss status
```

## nginx
```
apt-get install nginx
```
Configure firewall
```
ufw app list
```
```
ufw allow "Nginx Full"
```
```
ufw reload
```
```
vim /etc/nginx/sites-available/rss
```
```
server {
    server_name your_domain www.your_domain;

    location / {
        proxy_pass http://localhost:3000;
    }
}
```
```
server {
   listen 80;
   listen [::]:80;

   location / {
     proxy_pass http://localhost:3000;
   }
}
```
```
ln -s /etc/nginx/sites-available/rss /etc/nginx/sites-enabled/rss
```
```
rm -f /etc/nginx/sites-enabled/default
```
```
nginx -s reload
```

## Cron job
Stored under `/var/spool/cron/crontabs/`.
```
apt update
apt install cron
```
```
systemctl enable cron
```
Synchronizing state of cron.service with SysV service script with /lib/systemd/systemd-sysv-install.
Executing: /lib/systemd/systemd-sysv-install enable cron

```
crontab -e #edit
crontab -l #view
crontab -r #remove
```
Store environment variables in `.profile`
```
vim ~/.profile
```
```
export VAR_NAME=VALUE
...
```
```
0 1 * * * . $HOME/.profile; /home/go/src/rss/bin/cron
```
Debug Cron job
```
grep CRON /var/log/syslog
```

## SSL
[letsencrypt](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-20-04)
