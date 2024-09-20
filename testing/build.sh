docker buildx build -t vpa-manager:latest .

kind load docker-image vpa-manager:latest
