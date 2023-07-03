docker build -t go_server:latest .
sudo docker build -f Dockerfile .
docker run --rm -p 1058:1058/udp -ti --name skato go_server:latest
