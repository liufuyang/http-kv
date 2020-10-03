
## Start Prometheus

```
docker run \
    --rm -d --name prometheus\
    -p 9090:9090 \
    -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus

docker run --rm -d --name grafana -p 3000:3000 \
    -v $(pwd)/grafana:/etc/grafana \
    grafana/grafana
```
Add prometheus data source on `http://host.docker.internal:9090` (should already be provisioned with config files here). To know more https://grafana.com/tutorials/provision-dashboards-and-data-sources/#3

(Incase `host.docker.internal` not working, try `192.168.65.2` as the internal IP for process in container to reach host)