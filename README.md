# Hearchco agent repository built using Go

To self-host, you can use the official docker images:

```yaml
version: "3.9"
services:
  traefik:
    image: "traefik:latest"
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entryPoints.web.address=:80"
    ports:
      - "80:80"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
    restart: unless-stopped
  frontend:
    image: ghcr.io/hearchco/frontend:latest
    environment:
      - PUBLIC_URI=https://search.example.org
      - API_URI=http://agent:3030 # server reachable, used for SSR
      - PUBLIC_API_URI=https://api.search.example.org # client reachable, used for CSR
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.frontend.rule=Host(`search.example.org`)"
      - "traefik.http.routers.frontend.entrypoints=web"
      - "traefik.http.services.frontend.loadbalancer.server.port=3000"
  agent:
    image: ghcr.io/hearchco/agent:latest
    environment:
      - HEARCHCO_SERVER_FRONTENDURLS=http://localhost:5173,https://*search.example.org
      - HEARCHCO_SERVER_CACHE_TYPE=redis # set to "none" to disable caching, NOT RECOMMENDED
      - HEARCHCO_SERVER_CACHE_REDIS_HOST=redis
      # - HEARCHCO_SERVER_CACHE_REDIS_PASSWORD=redispassword # empty by default
      - HEARCHCO_SERVER_IMAGEPROXY_SECRETKEY=supersecretkey # used by image proxy for hashing image urls
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.agent.rule=Host(`api.search.example.org`)"
      - "traefik.http.routers.agent.entrypoints=web"
      - "traefik.http.services.agent.loadbalancer.server.port=3030"
  redis:
    image: redis:latest
    volumes:
      - "redis-vol-0:/data"
    restart: unless-stopped
volumes:
  redis-vol-0:
```
