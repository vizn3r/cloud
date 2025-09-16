FROM nginx:alpine

# Install bash for scripting
RUN apk add --no-cache curl

# Create nginx user and group
RUN addgroup -g 1001 -S nginx && \
    adduser -u 1001 -S nginx -G nginx

# Copy static files
COPY webclient/public/ /usr/share/nginx/html/

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Set proper ownership
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]