# Load Testing

Load tests use [k6](https://k6.io/) to simulate many frontends hitting the backend.

## Install k6

```bash
# Arch Linux
yay -S k6

# Or: https://k6.io/docs/getting-started/installation/
```

## Run Load Test

```bash
# Local backend (default: http://localhost:8080)
k6 run load-test.js

# Against deployed service
BASE_URL=https://your-ingress-url k6 run load-test.js

# Custom load: 100 virtual users for 2 minutes
k6 run --vus 100 --duration 2m load-test.js
```

## Thresholds

- **http_req_duration p(95) < 500ms**: 95% of requests complete in under 500ms
- **http_req_failed rate < 0.01**: Less than 1% of requests fail

Adjust in `load-test.js` if needed.
