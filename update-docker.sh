#echo "Fetching latest changes from git..."
#git pull
echo "Trying to stop old container"
docker stop noti-share
echo "Deleting old container"
docker container rm noti-share --force
echo "Deleting old image"
docker rmi dnieln7/noti-share
echo "Creating new image"
docker build . -t dnieln7/noti-share
echo "Creating new container"
docker run -d -p 4300:4444 --restart always --name noti-share dnieln7/noti-share
echo "Deleting dangling images"
docker image prune --force
