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
	-kubectl delete deployment my-nginx
	-kubectl delete deployment my-nginx

push:
	-docker build -t $(TAG) -t $(TAG):$(VER) .
	-docker push $(TAG):$(VER)
	
deploy:
	- kubectl apply -f choas_deployment.yaml

auth:
	- kubectl apply -f service-admin-role.yaml
	- kubectl create clusterrolebinding service-admin-pod --clusterrole=cluster-admin --serviceaccount=default:default

proxy:
	minikube service --url yasser-chaos

deploytest:
	- kubectl apply -f nginx_deployment.yaml