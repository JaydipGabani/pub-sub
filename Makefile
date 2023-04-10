deploy:
	kubectl get secret redis --namespace=default -o yaml | sed 's/namespace: .*/namespace: kube-system/' | kubectl apply -f -
	kubectl apply -f sub.yaml

build:
	docker build -f Dockerfile -t docker.io/jaydipgabani/dummy-subscriber:test .