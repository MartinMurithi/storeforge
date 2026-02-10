# Use official Envoy image
FROM envoyproxy/envoy:v1.29-latest

COPY envoy.yaml /etc/envoy/envoy.yaml
COPY protos/api.pb /etc/envoy/protos/api.pb

# Expose port that clients hit
EXPOSE 10000 9901

# Run Envoy
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml", "--log-level", "info"]