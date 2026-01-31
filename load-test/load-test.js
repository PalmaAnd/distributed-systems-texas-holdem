/**
 * k6 load test for Texas Hold'em backend API
 * Install: https://k6.io/docs/getting-started/installation/
 * Run: k6 run load-test.js
 * With options: k6 run --vus 50 --duration 30s load-test.js
 */
import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '30s', target: 20 },   // Ramp up to 20 users
    { duration: '1m', target: 20 },    // Stay at 20 users
    { duration: '30s', target: 50 },   // Ramp up to 50 users
    { duration: '1m', target: 50 },    // Stay at 50 users
    { duration: '30s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],    // Less than 1% failure rate
  },
};

export function setup() {
  return { baseUrl: BASE_URL };
}

export default function (data) {
  // Test /api/v1/evaluate
  const evaluateRes = http.post(
    `${data.baseUrl}/api/v1/evaluate`,
    JSON.stringify({
      hole_cards: ['HA', 'HK'],
      community_cards: ['HQ', 'HJ', 'HT', 'S2', 'D3'],
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(evaluateRes, {
    'evaluate status 200': (r) => r.status === 200,
    'evaluate has rank_name': (r) => JSON.parse(r.body).rank_name !== undefined,
  });

  sleep(0.5);

  // Test /api/v1/probability (lighter: fewer sims)
  const probRes = http.post(
    `${data.baseUrl}/api/v1/probability`,
    JSON.stringify({
      hole_cards: ['HA', 'HK'],
      community_cards: [],
      num_players: 2,
      num_sims: 1000,
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(probRes, {
    'probability status 200': (r) => r.status === 200,
    'probability has win_probability': (r) => JSON.parse(r.body).win_probability !== undefined,
  });

  sleep(0.5);
}
