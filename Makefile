NAME=yasser-chaos
TAG=yassermog/$(NAME)
VER=latest

all: clean build push deploy

build:
	docker build -t $(TAG) -t $(TAG):$(VER) .

run:
	docker run -d -p 7070:7070 -e PORT=7070 --name=$(NAME) $(TAG)

clean:
	-docker stop $(NAME)
	-docker rm $(NAME)
	-kubectl delete deployment yasser-chaos

push:
	docker push $(TAG):$(VER)
	
deploy:
	- kubectl apply -f Choas_deployment.yaml
