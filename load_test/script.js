import http from 'k6/http';
import { check, group } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 500 },
    { duration: '10s', target: 1000 },
    { duration: '20s', target: 8000 },
  ],
  thresholds: {
      http_req_failed: [
            { threshold: 'rate<0.15', abortOnFail: true }, // Автоматически остановить тест при превышении
        ],
        http_req_duration: [
          { threshold: 'p(95)<1000', abortOnFail: true }
        ],
    }
};

export default function () {
  const token = 'dCagXAppjfkLggxhgMfvYGVvbOYLvPiT';

    const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
    };

   group('API uptime check', () => {
       const response = http.get('http://localhost:3002/v2/dialog/2/list', {
        headers: headers
    });
       check(response, {
           "status code should be 200": res => res.status === 200,
       });
   });
};
