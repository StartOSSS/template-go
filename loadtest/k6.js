import http from 'k6/http';
import { check, sleep } from 'k6';

const target = __ENV.K6_TARGET || 'http://todo-app.todo.svc.cluster.local:8080';
const vus = Number(__ENV.K6_VUS || '20');
const duration = __ENV.K6_DURATION || '1m';

export const options = {
  stages: [
    { duration: duration, target: vus },
    { duration: '30s', target: 0 },
  ],
};

export default function () {
  const payload = JSON.stringify({ title: 'k6', description: 'load test' });
  const res = http.post(`${target}/todo.v1.TodoService/CreateTodo`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
  sleep(1);
}
