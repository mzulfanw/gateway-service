import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    high_load: {
      executor: "constant-arrival-rate",
      rate: 1000,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 200,
      maxVUs: 1000,
    },
  },
};

const BASE_URL = "http://localhost:3000";

export default function () {
  const payload = JSON.stringify({
    productId: "06d44075-00ce-4485-b7ff-b907ac17d45e",
    quantity: 1,
  });

  const headers = {
    "Content-Type": "application/json",
  };

  const res = http.post(`${BASE_URL}/orders`, payload, { headers });

  check(res, {
    "status is 201": (r) => r.status === 201,
    "response time < 500ms": (r) => r.timings.duration < 500,
  });

  sleep(0.001);
}
