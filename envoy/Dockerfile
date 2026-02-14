# Use official Envoy image
FROM envoyproxy/envoy:v1.30-latest

COPY envoy.yaml /etc/envoy/envoy.yaml
COPY api.pb /etc/envoy/protos/api.pb

# Expose port that clients hit
EXPOSE 10000 9901

# Run Envoy
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml", "--log-level", "info"]


# 1. Remove the old container
# docker rm -f envoy

# # 2. Run it again with the --add-host flag
# docker run -d --name envoy \
#     --add-host=host.docker.internal:host-gateway \
#     -p 10000:10000 \
#     -p 9901:9901 \
#     -v $(pwd)/envoy.yaml:/etc/envoy/envoy.yaml \
#     -v $(pwd)/api.pb:/etc/envoy/api.pb \
#     envoyproxy/envoy:v1.30-latest