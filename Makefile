deploy:
	kubectl create ns sub
	kubectl get secret redis --namespace=default -o yaml | sed 's/namespace: .*/namespace: sub/' | kubectl apply -f -
	kubectl apply -f sub.yaml

build:
	# docker build -f Dockerfile -t docker.io/jaydipgabani/fake-subscriber:latest .
	# kind load docker-image docker.io/jaydipgabani/fake-subscriber:latest
	docker build -f Dockerfile -t fake-subscriber:latest .
	kind load docker-image fake-subscriber:latest