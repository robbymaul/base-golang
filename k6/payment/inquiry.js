import http from "k6/http";

export function inquiry(baseUrl, data) {
  const response = http.post(baseUrl, JSON.stringify(data), {
    headers: {
      "Content-Type": "application/json",
      "Accept": "application/json",
    },
    timeout: "240s", // <- tambahkan ini
  });

  return response;
}