import http from 'k6/http';
import { check, sleep } from "k6";


const healthOptions = JSON.parse(open("health.json"))

export const options = healthOptions.options

export function setup() {
  const data = healthOptions.data
  return data
}

export default function (data) {
  const responseHealth = http.get(`${data.baseUrl}`);

  check(responseHealth, {
    "is status 200": (r) => r.status === 200,
    "is body not null": (r) => r.body !== null
  })

  sleep(1)
}
