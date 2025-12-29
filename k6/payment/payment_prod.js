import { check, fail } from "k6";
import { sleep } from "k6";
import {inquiry} from "./inquiry.js";
import {checkStatusPayment} from "./check_status_payment.js";

const paymentSchenarioJson = JSON.parse(open("payment_scenario_prod.json"));

const paymentScenario = paymentSchenarioJson.paymentScenario.scenario;

export const options = {
  thresholds: paymentSchenarioJson.thresholds,
  scenarios: {
    paymentScenario: paymentScenario,
  },
};

export function setup() {
  const data = paymentSchenarioJson.paymentScenario;

  return data;
}

export default function (data) {
  const orderId = `${new Date().getTime()}`;
  data.body.payments[0].orderId = orderId;
  data.body.apiKey = data.apiKey;
  data.body.secretKey = data.secretKey;

  const responsePayment = inquiry(`${data.baseUrl}/payments`, data.body);

  const payment = check(responsePayment, {
    "is status payment 200": (r) => r.status === 200,
  });

  if (!payment) {
    fail("payment failed", responsePayment.body);
  }

  const checkPaymentResponse = checkStatusPayment(
    `${data.baseUrl}/payments/status`,
    {
      apiKey: data.apiKey,
      secretKey: data.secretKey,
      orderId: orderId,
    }
  );
  check(checkPaymentResponse, {
    "is status check payment 200": (r) => r.status === 200,
  });

  sleep(1);
}
