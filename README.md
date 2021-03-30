# Chaos_Monkey
Chaos_Monkey

## Build 
```
make build

```
## Deployment
to deploy to kubernetes 

```
make deploy
```

in minikube you can get the service url by running 
```
 minikube service --url yasser-chaos
``` 

## Testing  

```
http://localhost:7070/hello?name=yasser
```