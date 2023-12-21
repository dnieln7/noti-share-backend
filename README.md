# NotiShare Backend

Back end for NotiShare mobile app.

# Deploy with commands

Create image

```shell
docker build . -t dnieln7/noti-share
```

Create and run container

```shell
docker run -d -p 4300:4444 --restart always --name noti-share dnieln7/noti-share
```

# Deploy with script

Give executable permissions

```shell
chmod +x update-docker.sh 
```

Run script

```shell
./update-docker.sh
```
