docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t jkbcz/amongus-konfero:0.0.1 \
  --push .