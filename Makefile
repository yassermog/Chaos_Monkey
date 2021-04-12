NAME=yasser-chaos
TAG=yassermog/$(NAME)
VER=latest

all: clean build push deploy

install: 
	helm install chaos-monkey ./chaos-monkey/ --set service.type=NodePort

build:
	docker build -t $(TAG) -t $(TAG):$(VER) .

run:
	docker run -d -p 7070:7070 -e PORT=7070 --name=$(NAME) $(TAG)

clean:
	-docker stop $(NAME)
	-docker rm $(NAME)
	-kubectl delete deployment chaos-monkey
	-kubectl delete deployment my-nginx
	-kubectl delete clusterrolebinding service-admin-pod
	-helm del chaos-monkey
	-helm del prometheus
	
push:
	-docker build -t $(TAG) -t $(TAG):$(VER) .
	-docker push $(TAG):$(VER)
	
deploy:
	- kubectl apply -f choas_deployment.yaml

auth:
	- kubectl apply -f service-admin-role.yaml
	- kubectl create clusterrolebinding service-admin-pod-chaos --clusterrole=cluster-admin --serviceaccount=default:chaos-monkey

proxy:
	- kubectl port-forward service/prometheus-kube-prometheus-prometheus 9090
	- kubectl port-forward deployment/prometheus-grafana 3000 
	- minikube service --url chaos-monkey

deploytest:
	- kubectl apply -f nginx_deployment.yaml

monitoring: 
	-helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	-helm repo update
	-helm install prometheus prometheus-community/kube-prometheus-stack