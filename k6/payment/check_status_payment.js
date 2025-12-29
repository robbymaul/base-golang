import http from "k6/http";

export function checkStatusPayment(baseUrl, data) {
  const response = http.post(baseUrl, JSON.stringify(data), {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
  });

  return response;
}
