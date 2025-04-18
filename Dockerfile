FROM debian:stable-slim

# COPY source destination
COPY chirpy /bin/chirpy

# Set Environment Variable PORT
ENV PORT=8080

# Command/Run
CMD ["/bin/chirpy"]