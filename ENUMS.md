# Enum Definitions

Dokumentasi tipe data enumerasi yang digunakan secara programatik oleh sistem, **tidak disimpan dalam database**.

## PaymentStatus

| Value        | Description                                |
|--------------|--------------------------------------------|
| `pending`    | Menunggu pembayaran                        |
| `processing` | Menunggu Proses Pembayaran Payment Gateway |
| `success`    | Pembayaran Berhasil                        |
| `failed`     | Pembayaran Gagal                           |
| `cancelled`  | Dibatalkan oleh user                       |
| `expired`    | Waktu Pembayaran telah Habis               |

## PaymentMethod

| Value     | Description                     |
|-----------|---------------------------------|
| `CASH`    | Pembayaran tunai                |
| `QRIS`    | Pembayaran QRIS                 |
| `EWALLET` | E-wallet seperti OVO, DANA, dsb |
| `VA`      | Virtual Account                 |

## FeeType

| Value        | Description                                 |
|--------------|---------------------------------------------|
| `fixed`      | Fee Pembayaran Dengan Fix Pembayaran        |
| `percentage` | Fee Pembayaran Dengan Persentase Pembayaran |
| `none`       | Fee Tanpa Status Apapun                     |

## CodePlatform

| Value      | Description   |
|------------|---------------|
| `web_knet` | platform knet |
| `web_sms`  | platform sms  |

## CodePaymentTypes

| Value          | Description                                 |
|----------------|---------------------------------------------|
| `sales_order`  | code payment type untuk sales order product |
| `topup_token`  | code payment type untuk topup token         |
| `topup_wallet` | code payment type untuk topup wallet        |

## CodePaymentMethods

| Value           | Description                                    |
|-----------------|------------------------------------------------|
| `bank_transfer` | code payment method transfer bank              |
| `e_wallet`      | code payment method transfer elektronik wallet |
| `credit_card`   | code payment method transfer credit card       |
| `va`            | code payment method transfer virtual account   |

## ProviderPaymentMethod

| Value       | Description                       |
|-------------|-----------------------------------|
| `midtrans`  | provider payment method midtrans  |
| `senangpay` | provider payment method senangpay |

## CodeAdminRoles

| Value         | Description                        |
|---------------|------------------------------------|
| `super_admin` | code admin roles untuk super admin |
| `admin`       | code admin roles untuk admin       |

## ActionAdminActivityLogs

| Value            | Description                                          |
|------------------|------------------------------------------------------|
| `login`          | action admin activity log aksi login                 |
| `logout`         | action admin activity log aksi logout                |
| `cancel_payment` | action admin activity log aksi pembatalan pembayaran |

## ResourceTypeAdminActivityLogs

| Value      | Description                                                                                                         |
|------------|---------------------------------------------------------------------------------------------------------------------|
| `payment`  | resource type admin activity log tipe sumber daya yang di lakan admin saat action pembayaran                        |
| `user`     | resource type admin activity log tipe sumber daya yang di lakan admin saat melakukan aksi admin                     |
| `platform` | resource type admin activity log tipe sumber daya yang di lakan admin saat melakukan aksi platform payment gatewayF |
