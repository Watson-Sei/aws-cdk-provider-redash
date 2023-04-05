FROM redash/redash:10.1.0.b50633

# Switch to the root user temporarily
USER root

# Create the directory and set the appropriate permissions
RUN mkdir -p /var/lib/apt/lists/partial && chown _apt /var/lib/apt/lists/partial

# Install nginx
RUN apt-get update && apt-get install -y nginx

COPY nginx.conf /etc/nginx/conf.d/default.conf

CMD service nginx start && /app/bin/docker-entrypoint