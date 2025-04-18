package errs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	InternalServerError = "Terjadi kesalahan pada server, silakan coba beberpa saat lagi"
	BadRequest          = "Permintaan tidak sesuai"
	NotFound            = "Data tidak ditemukan"
)

type CodeErr int

const (

	// Maintenance
	MAINTENANCE_MODE             CodeErr = 1000
	PROVIDER_PRODUCT_MAINTENANCE CodeErr = 1100

	// Common Error
	GENERAL_ERROR                        CodeErr = 9999
	CONNECTION_RPC_ERROR                 CodeErr = 9998
	CONNECTION_HTTP_ERROR                CodeErr = 9997
	MARSHALLING_FAILED                   CodeErr = 9996
	UNMARSHALLING_FAILED                 CodeErr = 9995
	PROVIDER_SERVER_DOWN                 CodeErr = 9994
	PROVIDER_PRODUCT_PROBLEM             CodeErr = 9993
	PROVIDER_FAILED_UNKNOWN              CodeErr = 9992
	CONNECTION_PROVIDER_TIMEOUT          CodeErr = 9991
	S3_PUT_FAILED                        CodeErr = 9990
	S3_GET_FAILED                        CodeErr = 9989
	FCM_CANNOT_INIT_APP                  CodeErr = 9988
	FCM_CANNOT_SEND_MESSAGE              CodeErr = 9987
	ALL_PRODUCT_FAIL                     CodeErr = 9986
	KLIKOO_PRODUCT_FAILED_MULTIPLE_TIMES CodeErr = 9985
	ISSUER_OR_SWITCH_IS_INOPERATIVE      CodeErr = 9984
	ADD_FAVOURITE_FAILED                 CodeErr = 9983
	ADD_FAVOURITE_SUCCESS                CodeErr = 9982
	DELETE_FAVOURITE_SUCCESS             CodeErr = 9981
	PIN_MAX_TRY_COUNT_REACHED            CodeErr = 9980
	PARTIIAL_PAYMENT_FAILED              CodeErr = 9979
	PLEASE_PAY_AT_BRANCH_OFFICE          CodeErr = 9978
	BILLER_ERROR                         CodeErr = 9977
	UNIMPLEMENTED                        CodeErr = 9976
	USER_NOT_SET                         CodeErr = 9975
	FAIL_TO_UPSERT_BULK                  CodeErr = 9974
	PRODUCT_IS_IN_MAINTENANCE            CodeErr = 9973
	CONTEXT_DEADLINE_EXCEEDED            CodeErr = 9972
	CONTEXT_CANCELLED                    CodeErr = 9971

	// Not Found Error
	USER_NOT_FOUND          CodeErr = 7001
	CUSTOMER_NOT_FOUND      CodeErr = 7002
	PRODUCT_NOT_FOUND       CodeErr = 7003
	DATA_NOT_FOUND          CodeErr = 7004
	TRANSACTION_NOT_FOUND   CodeErr = 7005
	BILLING_NOT_FOUND       CodeErr = 7006
	PROVIDER_NOT_FOUND      CodeErr = 7007
	REFERRAL_NOT_FOUND      CodeErr = 7008
	CUSTOMER_ACCOUNT_BLOCK  CodeErr = 7009
	BILL_HAS_PAID           CodeErr = 7010
	CUSTOMER_ID_WRONG       CodeErr = 7011
	LAST_INSTALLMENT        CodeErr = 7012
	USER_ADMIN_NOT_FOUND    CodeErr = 7013
	BILL_NOT_FOUND          CodeErr = 7014
	TICKET_NOT_FOUND        CodeErr = 7015
	RESOURCE_NOT_FOUND      CodeErr = 7016
	PROVIDER_PRODUCT_CLOSED CodeErr = 7017

	// Validation Error
	INVALID_TOKEN                       CodeErr = 6500
	UNAUTHORIZED                        CodeErr = 6501
	SESSION_EXPIRED                     CodeErr = 6502
	INVALID_REFRESH_TOKEN               CodeErr = 6503
	REFRESH_TOKEN_EXPIRED               CodeErr = 6504
	INVALID_TIMESTAMP                   CodeErr = 6505
	NOT_MATCH_DEVICE                    CodeErr = 6506
	IP_BLOCKED                          CodeErr = 6507
	INVALID_SIGNATURE                   CodeErr = 6508
	BAD_REQUEST                         CodeErr = 6509
	VALIDATION_ERROR                    CodeErr = 6510
	ACCESS_DENIED                       CodeErr = 6511
	REQUIRED_FIELD                      CodeErr = 6512
	TOO_MANY_REQUEST                    CodeErr = 6513
	REQUEST_TO_EARLY                    CodeErr = 6514
	DATA_DUPLICATED                     CodeErr = 6515
	MIN_REF_ID                          CodeErr = 6516
	MAX_REF_ID                          CodeErr = 6517
	BALANCE_INSUFFICIENT                CodeErr = 6518
	PIN_NOT_MATCH                       CodeErr = 6519
	INVALID_PHONE_NUMBER                CodeErr = 6520
	OTP_TOO_MANY_VERIFICATION_ATTEMPT   CodeErr = 6521
	OTP_NOT_CORRECT                     CodeErr = 6522
	OTP_NUMBER_BLOCKED                  CodeErr = 6523
	NEED_TO_VERIFY                      CodeErr = 6524
	PRODUCT_NOT_VALID                   CodeErr = 6525
	INVOICE_AMOUNT_DOES_NOT_MATCH       CodeErr = 6526
	INQUIRY_GENERAL_ERROR               CodeErr = 6527
	MAX_AMOUNT_VALIDATION_CHECK         CodeErr = 6528
	MIN_AMOUNT_VALIDATION_CHECK         CodeErr = 6529
	MUST_BE_PAID_AT_THE_PROVIDER_LOCKET CodeErr = 6530
	TRANSACTION_REPEATED                CodeErr = 6531
	NOT_SUBSCRIBE                       CodeErr = 6532
	INVALID_AUTH                        CodeErr = 6533
	BAD_REQUEST_HEADER                  CodeErr = 6534
	TOKEN_EXPIRED                       CodeErr = 6535
	USERNAME_TOO_SHORT                  CodeErr = 6536
	NO_USER_PROFILE                     CodeErr = 6537
	PIN_DIDNT_MATCH                     CodeErr = 6538
	PIN_NOT_SIX_CHAR                    CodeErr = 6539
	USER_DIDNT_HAVE_PIN                 CodeErr = 6540
	USER_B2B_NOT_ACTIVE                 CodeErr = 6541
	FAILED_TO_VERIFY                    CodeErr = 6542
	FACE_RECOGNITION_NOT_SAME           CodeErr = 6543
	FAILED_OCR                          CodeErr = 6544
	REFERRAL_ALREADY_APPLY              CodeErr = 6545
	CREATE_PASSWORD_ERROR               CodeErr = 6546
	INPUT_PASSWORD_ERROR                CodeErr = 6547
	TOKEN_USED                          CodeErr = 6548
	PVG_EMAIL_NOT_FOUND                 CodeErr = 6549
	EMAIL_ALREADY_REGISTERED            CodeErr = 6550
	PHONE_ALREADY_REGISTERED            CodeErr = 6551
	GENERATE_VA_REACH_MAX               CodeErr = 6552
	OTP_NEED_TO_WAIT                    CodeErr = 6553
	NEED_TO_INQUIRY                     CodeErr = 6554
	NOT_ELIGIBLE_TRANSACTION            CodeErr = 6555
	WRONG_CUSTOMER_NUMBER               CodeErr = 6557
	WALLET_NOT_ENOUGH                   CodeErr = 6558
	ALREADY_SUBSCRIBE                   CodeErr = 6559
	REQUEST_CANNOT_NIL                  CodeErr = 6560
	PASSWORD_ALREADY_SET                CodeErr = 6561
	INVALID_ROLE_ID                     CodeErr = 6562
	DUPLICATE_REF_ID                    CodeErr = 6563
	DATA_TOO_LARGE                      CodeErr = 6564
	INVALID_IP                          CodeErr = 6565
	INVALID_DEVICE_ID                   CodeErr = 6566
	USER_TYPE_IS_NOT_B2B                CodeErr = 6567
	INVALID_UUID                        CodeErr = 6568
	INVALID_GRANT_TYPE                  CodeErr = 6569
	INVALID_CREDENTIALS                 CodeErr = 6570
	IP_NOT_WHITELISTED                  CodeErr = 6571
	USER_TYPE_IS_NOT_IRS                CodeErr = 6572
	USER_HAVE_ACTIVE_PIN                CodeErr = 6573
	PARSE_TIME_FAILED                   CodeErr = 6574
	DELETE_OWN_USER                     CodeErr = 6575
	DELETE_OWNER_USER                   CodeErr = 6576
	LAST_INVITATION_LESS_THAN_24_HOURS  CodeErr = 6577
	AMOUNT_MORE_THAN_SELLING_PRICE      CodeErr = 6578
	EMPTY_FILE                          CodeErr = 6579
	SAME_PASSWORD_AS_BEFORE             CodeErr = 6580
	INVALID_FORMAT_PUBLIC_KEY           CodeErr = 6581
	INVALID_FORMAT_URL                  CodeErr = 6582
	FAILED_SENT_EMAIL                   CodeErr = 6583
	FRAUD_NUMBER                        CodeErr = 6584
	PHONE_NUMBER_BLOCKED                CodeErr = 6585
	DATE_RANGE_EXCEEDS_LIMIT            CodeErr = 6586
	EMAIL_NOT_FOUND                     CodeErr = 6587
	PHONE_NUMBER_NONACTIVE              CodeErr = 6588
	MAX_BALANCE_EXCEEDED                CodeErr = 6589
	PASSWORD_NOT_MATCH                  CodeErr = 6590
	SEAT_BOOKED                         CodeErr = 6591
	TRANSACTION_FAILED_TRY_AGAIN        CodeErr = 6592
	EDIT_OWNER_USER                     CodeErr = 6593
	EDIT_OWN_USER                       CodeErr = 6594
	EDIT_TO_OWNER_ROLES                 CodeErr = 6595
	EMAIL_DOMAIN_NOT_MATCH_WITH_OWNER   CodeErr = 6596
	BOOKING_TICKET_EXPIRED              CodeErr = 6597
	CHANNEL_TOPUP_NOT_ALOWED            CodeErr = 6598
	TOPUP_EXPIRED                       CodeErr = 6599
	PAYMENT_CHANNEL_NOT_IMPLEMENTED     CodeErr = 6600
	PAYMENT_CHANNEL_IS_DISABLE          CodeErr = 6601
	PAYMENT_SOURCE_NOT_IMPLEMENTED      CodeErr = 6602
	TOTAL_B2B_USERS_EXCEEDS_THE_LIMIT   CodeErr = 6603
	PAYMENT_STATUS_IS_FINAL             CodeErr = 6604
	CANCEL_PAYMENT_FAILED               CodeErr = 6605
	TOPUP_STATUS_IS_FINAL               CodeErr = 6606
	USER_HAS_TOPUP_ACTIVE               CodeErr = 6607
	MAXIMUM_KTP_USE                     CodeErr = 6608
	INVALID_INTEGER_VARIABLE            CodeErr = 6609
	USER_DONT_HAVE_TOPUP_ACTIVE         CodeErr = 6610
	AMOUNT_MISMATCH                     CodeErr = 6611
	SAME_EXCHANGE_CURRENCY              CodeErr = 6612
	INVALID_QRIS                        CodeErr = 6613
	TOKEN_ALREADY_EXISTS                CodeErr = 6614
	INVALID_DATE_FORMAT                 CodeErr = 6615
	GROUP_BUY_FULL                      CodeErr = 6616
	GROUP_BUY_NOT_FOUND                 CodeErr = 6617
	GROUP_BUY_EXPIRED                   CodeErr = 6618
	GROUP_BUY_NOT_AVAILABLE             CodeErr = 6619
	GROUP_BUY_CLOSED                    CodeErr = 6620
	GROUP_BUY_ALREADY_JOINED            CodeErr = 6621
	CANNOT_CANCEL_ORDER                 CodeErr = 6622
	PAYMENT_METHOD_NOT_ALLOWED          CodeErr = 6623
	INVALID_REJECT_REASON               CodeErr = 6624
	INVALID_STATUS                      CodeErr = 6625
	INVALID_CONTENT_TYPE                CodeErr = 6626
	INVALID_AMOUNT                      CodeErr = 6627
	NUMBER_ONLY_AMOUNT                  CodeErr = 6628
	MAX_BEST_SELLER_DATA                CodeErr = 6629
	INVALID_CATEGORY_BEST_SELLER        CodeErr = 6630
	PAYMENT_METHOD_MIN_VERSION          CodeErr = 6631
	PRODUCT_INQUIRY_NOT_ALLOWED         CodeErr = 6632
	MINIMUM_AMOUNT                      CodeErr = 6633
	MAXIMUM_AMOUNT                      CodeErr = 6634
	INVALID_MULTIPLES_AMOUNT            CodeErr = 6635
	OTP_EXPIRED                         CodeErr = 6636

	// SQL Database Error
	SQL_SYNTAX_ERROR     CodeErr = 8500
	SQL_ERROR_NO_ROWS    CodeErr = 8501
	SQL_CONNECTION_ERROR CodeErr = 8502
	SQL_INSERT_ERROR     CodeErr = 8503
	SQL_UPDATE_ERROR     CodeErr = 8504
	SQL_GET_DATA_ERROR   CodeErr = 8505

	// Redis Database Error
	REDIS_NIL              CodeErr = 8600
	REDIS_INSERT_FAILED    CodeErr = 8601
	REDIS_UPDATE_FAILED    CodeErr = 8602
	REDIS_DELETE_FAILED    CodeErr = 8603
	REDIS_CONNECTION_ERROR CodeErr = 8604

	// Brokers Kafka Error
	KAFKA_DIAL_ERROR       CodeErr = 8700
	KAFKA_CONNECTION_ERROR CodeErr = 8701
	KAFKA_PUBLISH_ERROR    CodeErr = 8702

	// Brokers NSQ Error
	NSQ_DIAL_ERROR       CodeErr = 8800
	NSQ_CONNECTION_ERROR CodeErr = 8801
	NSQ_PUBLISH_ERROR    CodeErr = 8802
)

var mapCodeErrStatusCode = map[CodeErr]int{
	GENERAL_ERROR:                        http.StatusInternalServerError,
	CONNECTION_RPC_ERROR:                 http.StatusServiceUnavailable,
	MARSHALLING_FAILED:                   http.StatusInternalServerError,
	UNMARSHALLING_FAILED:                 http.StatusInternalServerError,
	UNAUTHORIZED:                         http.StatusUnauthorized,
	NOT_SUBSCRIBE:                        http.StatusForbidden,
	USER_ADMIN_NOT_FOUND:                 http.StatusNotFound,
	ACCESS_DENIED:                        http.StatusForbidden,
	INVALID_AUTH:                         http.StatusUnauthorized,
	BAD_REQUEST:                          http.StatusBadRequest,
	BAD_REQUEST_HEADER:                   http.StatusBadRequest,
	DATA_NOT_FOUND:                       http.StatusNotFound,
	TICKET_NOT_FOUND:                     http.StatusNotFound,
	REDIS_INSERT_FAILED:                  http.StatusInternalServerError,
	REDIS_UPDATE_FAILED:                  http.StatusInternalServerError,
	REDIS_DELETE_FAILED:                  http.StatusInternalServerError,
	TOKEN_EXPIRED:                        http.StatusUnauthorized,
	REFRESH_TOKEN_EXPIRED:                http.StatusUnauthorized,
	USER_NOT_FOUND:                       http.StatusNotFound,
	USERNAME_TOO_SHORT:                   http.StatusBadRequest,
	NO_USER_PROFILE:                      http.StatusNotFound,
	PIN_DIDNT_MATCH:                      http.StatusBadRequest,
	PIN_NOT_SIX_CHAR:                     http.StatusNotAcceptable,
	USER_DIDNT_HAVE_PIN:                  http.StatusBadRequest,
	USER_B2B_NOT_ACTIVE:                  http.StatusNotAcceptable,
	FAILED_TO_VERIFY:                     http.StatusPreconditionFailed,
	FACE_RECOGNITION_NOT_SAME:            http.StatusNotAcceptable,
	FAILED_OCR:                           http.StatusBadRequest,
	REFERRAL_NOT_FOUND:                   http.StatusNotFound,
	CUSTOMER_ACCOUNT_BLOCK:               http.StatusNotFound,
	BILL_HAS_PAID:                        http.StatusNotFound,
	CUSTOMER_ID_WRONG:                    http.StatusNotFound,
	LAST_INSTALLMENT:                     http.StatusNotFound,
	REFERRAL_ALREADY_APPLY:               http.StatusConflict,
	NEED_TO_VERIFY:                       http.StatusForbidden,
	S3_PUT_FAILED:                        http.StatusInternalServerError,
	S3_GET_FAILED:                        http.StatusInternalServerError,
	FCM_CANNOT_INIT_APP:                  http.StatusInternalServerError,
	FCM_CANNOT_SEND_MESSAGE:              http.StatusInternalServerError,
	OTP_NEED_TO_WAIT:                     http.StatusBadRequest,
	OTP_NUMBER_BLOCKED:                   http.StatusForbidden,
	OTP_NOT_CORRECT:                      http.StatusBadRequest,
	OTP_TOO_MANY_VERIFICATION_ATTEMPT:    http.StatusConflict,
	GENERATE_VA_REACH_MAX:                http.StatusTooManyRequests,
	PROVIDER_PRODUCT_PROBLEM:             http.StatusInternalServerError,
	PROVIDER_SERVER_DOWN:                 http.StatusServiceUnavailable,
	PROVIDER_FAILED_UNKNOWN:              http.StatusInternalServerError,
	PROVIDER_NOT_FOUND:                   http.StatusNotFound,
	CUSTOMER_NOT_FOUND:                   http.StatusNotFound,
	ALL_PRODUCT_FAIL:                     http.StatusInternalServerError,
	PRODUCT_NOT_VALID:                    http.StatusBadRequest,
	NEED_TO_INQUIRY:                      http.StatusPreconditionFailed,
	KLIKOO_PRODUCT_FAILED_MULTIPLE_TIMES: http.StatusNotAcceptable,
	NOT_ELIGIBLE_TRANSACTION:             http.StatusNotAcceptable,
	ISSUER_OR_SWITCH_IS_INOPERATIVE:      http.StatusInternalServerError,
	WRONG_CUSTOMER_NUMBER:                http.StatusBadRequest,
	WALLET_NOT_ENOUGH:                    http.StatusUnprocessableEntity,
	ADD_FAVOURITE_FAILED:                 http.StatusBadRequest,
	ADD_FAVOURITE_SUCCESS:                http.StatusBadRequest,
	DELETE_FAVOURITE_SUCCESS:             http.StatusBadRequest,
	PIN_MAX_TRY_COUNT_REACHED:            http.StatusForbidden,
	PARTIIAL_PAYMENT_FAILED:              http.StatusInternalServerError,
	ALREADY_SUBSCRIBE:                    http.StatusConflict,
	CREATE_PASSWORD_ERROR:                http.StatusBadRequest,
	INPUT_PASSWORD_ERROR:                 http.StatusBadRequest,
	TOKEN_USED:                           http.StatusConflict,
	PVG_EMAIL_NOT_FOUND:                  http.StatusNotFound,
	EMAIL_ALREADY_REGISTERED:             http.StatusConflict,
	PHONE_ALREADY_REGISTERED:             http.StatusConflict,
	REQUEST_CANNOT_NIL:                   http.StatusBadRequest,
	BILL_NOT_FOUND:                       http.StatusNotFound,
	PLEASE_PAY_AT_BRANCH_OFFICE:          http.StatusNotAcceptable,
	BILLER_ERROR:                         http.StatusInternalServerError,
	TOO_MANY_REQUEST:                     http.StatusTooManyRequests,
	MAINTENANCE_MODE:                     http.StatusServiceUnavailable,
	CONNECTION_HTTP_ERROR:                http.StatusServiceUnavailable,
	CONNECTION_PROVIDER_TIMEOUT:          http.StatusServiceUnavailable,
	PRODUCT_NOT_FOUND:                    http.StatusNotFound,
	TRANSACTION_NOT_FOUND:                http.StatusNotFound,
	BILLING_NOT_FOUND:                    http.StatusNotFound,
	INVALID_TOKEN:                        http.StatusUnauthorized,
	SESSION_EXPIRED:                      http.StatusUnauthorized,
	INVALID_REFRESH_TOKEN:                http.StatusUnauthorized,
	INVALID_TIMESTAMP:                    http.StatusBadRequest,
	NOT_MATCH_DEVICE:                     http.StatusForbidden,
	IP_BLOCKED:                           http.StatusForbidden,
	INVALID_SIGNATURE:                    http.StatusBadRequest,
	VALIDATION_ERROR:                     http.StatusBadRequest,
	REQUIRED_FIELD:                       http.StatusBadRequest,
	REQUEST_TO_EARLY:                     http.StatusTooEarly,
	DATA_DUPLICATED:                      http.StatusConflict,
	MIN_REF_ID:                           http.StatusBadRequest,
	MAX_REF_ID:                           http.StatusBadRequest,
	BALANCE_INSUFFICIENT:                 http.StatusBadRequest,
	PIN_NOT_MATCH:                        http.StatusBadRequest,
	USER_HAVE_ACTIVE_PIN:                 http.StatusBadRequest,
	INVALID_PHONE_NUMBER:                 http.StatusBadRequest,
	SQL_SYNTAX_ERROR:                     http.StatusInternalServerError,
	SQL_ERROR_NO_ROWS:                    http.StatusNotFound,
	SQL_CONNECTION_ERROR:                 http.StatusInternalServerError,
	SQL_INSERT_ERROR:                     http.StatusInternalServerError,
	SQL_UPDATE_ERROR:                     http.StatusInternalServerError,
	SQL_GET_DATA_ERROR:                   http.StatusInternalServerError,
	REDIS_NIL:                            http.StatusNotFound,
	REDIS_CONNECTION_ERROR:               http.StatusInternalServerError,
	KAFKA_DIAL_ERROR:                     http.StatusInternalServerError,
	KAFKA_CONNECTION_ERROR:               http.StatusInternalServerError,
	KAFKA_PUBLISH_ERROR:                  http.StatusInternalServerError,
	NSQ_DIAL_ERROR:                       http.StatusInternalServerError,
	NSQ_CONNECTION_ERROR:                 http.StatusInternalServerError,
	NSQ_PUBLISH_ERROR:                    http.StatusInternalServerError,
	INVOICE_AMOUNT_DOES_NOT_MATCH:        http.StatusBadRequest,
	INQUIRY_GENERAL_ERROR:                http.StatusBadRequest,
	MAX_AMOUNT_VALIDATION_CHECK:          http.StatusBadRequest,
	MIN_AMOUNT_VALIDATION_CHECK:          http.StatusBadRequest,
	MUST_BE_PAID_AT_THE_PROVIDER_LOCKET:  http.StatusBadRequest,
	TRANSACTION_REPEATED:                 http.StatusBadRequest,
	PASSWORD_ALREADY_SET:                 http.StatusBadRequest,
	INVALID_ROLE_ID:                      http.StatusBadRequest,
	DUPLICATE_REF_ID:                     http.StatusConflict,
	DATA_TOO_LARGE:                       http.StatusRequestEntityTooLarge,
	UNIMPLEMENTED:                        http.StatusNotImplemented,
	USER_NOT_SET:                         http.StatusInternalServerError,
	INVALID_IP:                           http.StatusBadRequest,
	INVALID_DEVICE_ID:                    http.StatusBadRequest,
	USER_TYPE_IS_NOT_B2B:                 http.StatusForbidden,
	INVALID_UUID:                         http.StatusBadRequest,
	INVALID_GRANT_TYPE:                   http.StatusBadRequest,
	INVALID_CREDENTIALS:                  http.StatusUnauthorized,
	IP_NOT_WHITELISTED:                   http.StatusForbidden,
	USER_TYPE_IS_NOT_IRS:                 http.StatusForbidden,
	PARSE_TIME_FAILED:                    http.StatusBadRequest,
	DELETE_OWN_USER:                      http.StatusBadRequest,
	DELETE_OWNER_USER:                    http.StatusBadRequest,
	LAST_INVITATION_LESS_THAN_24_HOURS:   http.StatusBadRequest,
	AMOUNT_MORE_THAN_SELLING_PRICE:       http.StatusBadRequest,
	EMPTY_FILE:                           http.StatusBadRequest,
	SAME_PASSWORD_AS_BEFORE:              http.StatusBadRequest,
	INVALID_FORMAT_PUBLIC_KEY:            http.StatusBadRequest,
	INVALID_FORMAT_URL:                   http.StatusBadRequest,
	FAILED_SENT_EMAIL:                    http.StatusBadRequest,
	FRAUD_NUMBER:                         http.StatusBadRequest,
	PHONE_NUMBER_BLOCKED:                 http.StatusBadRequest,
	DATE_RANGE_EXCEEDS_LIMIT:             http.StatusBadRequest,
	EMAIL_NOT_FOUND:                      http.StatusBadRequest,
	PHONE_NUMBER_NONACTIVE:               http.StatusBadRequest,
	MAX_BALANCE_EXCEEDED:                 http.StatusBadRequest,
	PASSWORD_NOT_MATCH:                   http.StatusBadRequest,
	SEAT_BOOKED:                          http.StatusConflict,
	TRANSACTION_FAILED_TRY_AGAIN:         http.StatusBadRequest,
	EDIT_OWNER_USER:                      http.StatusBadRequest,
	EDIT_OWN_USER:                        http.StatusBadRequest,
	EDIT_TO_OWNER_ROLES:                  http.StatusBadRequest,
	CHANNEL_TOPUP_NOT_ALOWED:             http.StatusBadRequest,
	EMAIL_DOMAIN_NOT_MATCH_WITH_OWNER:    http.StatusBadRequest,
	BOOKING_TICKET_EXPIRED:               http.StatusBadRequest,
	TOPUP_EXPIRED:                        http.StatusBadRequest,
	FAIL_TO_UPSERT_BULK:                  http.StatusBadRequest,
	PAYMENT_CHANNEL_NOT_IMPLEMENTED:      http.StatusBadRequest,
	PAYMENT_CHANNEL_IS_DISABLE:           http.StatusServiceUnavailable,
	PAYMENT_SOURCE_NOT_IMPLEMENTED:       http.StatusBadRequest,
	TOTAL_B2B_USERS_EXCEEDS_THE_LIMIT:    http.StatusBadRequest,
	PAYMENT_STATUS_IS_FINAL:              http.StatusBadRequest,
	CANCEL_PAYMENT_FAILED:                http.StatusBadRequest,
	TOPUP_STATUS_IS_FINAL:                http.StatusBadRequest,
	USER_HAS_TOPUP_ACTIVE:                http.StatusBadRequest,
	MAXIMUM_KTP_USE:                      http.StatusBadRequest,
	INVALID_INTEGER_VARIABLE:             http.StatusBadRequest,
	USER_DONT_HAVE_TOPUP_ACTIVE:          http.StatusNotFound,
	AMOUNT_MISMATCH:                      http.StatusBadRequest,
	SAME_EXCHANGE_CURRENCY:               http.StatusBadRequest,
	INVALID_QRIS:                         http.StatusBadRequest,
	TOKEN_ALREADY_EXISTS:                 http.StatusBadRequest,
	INVALID_DATE_FORMAT:                  http.StatusBadRequest,
	RESOURCE_NOT_FOUND:                   http.StatusNotFound,
	GROUP_BUY_FULL:                       http.StatusForbidden,
	GROUP_BUY_NOT_FOUND:                  http.StatusNotFound,
	GROUP_BUY_EXPIRED:                    http.StatusGone,
	GROUP_BUY_NOT_AVAILABLE:              http.StatusNotFound,
	GROUP_BUY_CLOSED:                     http.StatusGone,
	GROUP_BUY_ALREADY_JOINED:             http.StatusConflict,
	CANNOT_CANCEL_ORDER:                  http.StatusForbidden,
	PAYMENT_METHOD_NOT_ALLOWED:           http.StatusForbidden,
	INVALID_REJECT_REASON:                http.StatusBadRequest,
	INVALID_STATUS:                       http.StatusBadRequest,
	INVALID_CONTENT_TYPE:                 http.StatusBadRequest,
	PRODUCT_IS_IN_MAINTENANCE:            http.StatusBadRequest,
	INVALID_AMOUNT:                       http.StatusBadRequest,
	NUMBER_ONLY_AMOUNT:                   http.StatusBadRequest,
	MAX_BEST_SELLER_DATA:                 http.StatusBadRequest,
	INVALID_CATEGORY_BEST_SELLER:         http.StatusBadRequest,
	PAYMENT_METHOD_MIN_VERSION:           http.StatusBadRequest,
	PRODUCT_INQUIRY_NOT_ALLOWED:          http.StatusForbidden,
	PROVIDER_PRODUCT_MAINTENANCE:         http.StatusServiceUnavailable,
	PROVIDER_PRODUCT_CLOSED:              http.StatusNotFound,
	MINIMUM_AMOUNT:                       http.StatusBadRequest,
	MAXIMUM_AMOUNT:                       http.StatusBadRequest,
	INVALID_MULTIPLES_AMOUNT:             http.StatusBadRequest,
	CONTEXT_DEADLINE_EXCEEDED:            http.StatusInternalServerError,
	CONTEXT_CANCELLED:                    http.StatusInternalServerError,
	OTP_EXPIRED:                          http.StatusBadRequest,
}

var mapCodeErrMessage = map[CodeErr]string{
	GENERAL_ERROR:                        "Harap hubungi admin",
	CONNECTION_RPC_ERROR:                 "Koneksi RPC gagal",
	MARSHALLING_FAILED:                   "Gagal untuk marshal",
	UNMARSHALLING_FAILED:                 "Gagal untuk unmarshal",
	UNAUTHORIZED:                         "Unauthorized",
	NOT_SUBSCRIBE:                        "Tidak ada langganan",
	USER_ADMIN_NOT_FOUND:                 "Pengguna Admin tidak ditemukan",
	ACCESS_DENIED:                        "Akses ditolak",
	INVALID_AUTH:                         "Invalid Authorization",
	BAD_REQUEST:                          "Permintaan tidak lengkap",
	BAD_REQUEST_HEADER:                   "Permintaan tidak sesuai",
	DATA_NOT_FOUND:                       NotFound,
	TICKET_NOT_FOUND:                     "Tiket tidak ditemukan",
	REDIS_INSERT_FAILED:                  "Gagal masukan data ke redis",
	REDIS_UPDATE_FAILED:                  "Gagal ubah data di redis",
	REDIS_DELETE_FAILED:                  "Gagal hapus data di redis",
	TOKEN_EXPIRED:                        "Token kadaluarsa",
	REFRESH_TOKEN_EXPIRED:                "Refresh token kadaluarsa",
	USER_NOT_FOUND:                       "Gagal ambil data user",
	USERNAME_TOO_SHORT:                   "username terlalu pendek",
	NO_USER_PROFILE:                      "Data pengguna tidak ada, lakukan OCR dan Face Recognition terlebih dahulu",
	PIN_DIDNT_MATCH:                      "Pin tidak sesuai",
	PIN_NOT_SIX_CHAR:                     "Pin harus 6 digit",
	USER_DIDNT_HAVE_PIN:                  "Pengguna tidak mempunyai pin",
	USER_B2B_NOT_ACTIVE:                  "Pengguna belum aktif",
	FAILED_TO_VERIFY:                     "Gagal memverikasi pengguna",
	FACE_RECOGNITION_NOT_SAME:            "Selfie tidak sama dengan foto KTP",
	FAILED_OCR:                           "OCR gagal, pastikan gambar sesuai",
	REFERRAL_NOT_FOUND:                   "Referral tidak ditemukan",
	CUSTOMER_ACCOUNT_BLOCK:               "Akun telah dinonaktifkan atau diblokir, silahkan hubungi customer support terkait.",
	BILL_HAS_PAID:                        "Tagihan sudah dibayar",
	CUSTOMER_ID_WRONG:                    "No Pelanggan salah, pastikan nomor yang diinput sesuai.",
	LAST_INSTALLMENT:                     "Angsuran terakhir. Silakan kunjungi kantor cabang.",
	REFERRAL_ALREADY_APPLY:               "Sudah pernah mengambil referral",
	NEED_TO_VERIFY:                       "Butuh verifikasi terlebih dahulu",
	S3_PUT_FAILED:                        "Gagal simpan data ke storage server",
	S3_GET_FAILED:                        "Gagal simpan data dari storage server",
	FCM_CANNOT_INIT_APP:                  "Gagal inisialisasi FCM",
	FCM_CANNOT_SEND_MESSAGE:              "Gagal untuk mengirim notifikasi ke pengguna",
	OTP_NEED_TO_WAIT:                     "Mohon tunggu untuk mengirimkan OTP berikutnya",
	OTP_NUMBER_BLOCKED:                   "OTP akan tersedia dalam 1 jam",
	OTP_NOT_CORRECT:                      "OTP tidak sesuai",
	OTP_TOO_MANY_VERIFICATION_ATTEMPT:    "Terlalu banyak memasukan OTP",
	GENERATE_VA_REACH_MAX:                "Pembuatan custom VA mencapai batas maksimum",
	PROVIDER_PRODUCT_PROBLEM:             "Produk sedang dalam masalah",
	PROVIDER_SERVER_DOWN:                 "Provider server sedang ada masalah",
	PROVIDER_FAILED_UNKNOWN:              "Suplier tidak diketahui",
	PROVIDER_NOT_FOUND:                   "Provider tidak ditemukan",
	CUSTOMER_NOT_FOUND:                   "No Customer tidak ditemukan",
	ALL_PRODUCT_FAIL:                     "Semua produk gagal",
	PRODUCT_NOT_VALID:                    "Produk tidak sesuai",
	NEED_TO_INQUIRY:                      "Harus inquiry terlebih dahulu",
	KLIKOO_PRODUCT_FAILED_MULTIPLE_TIMES: "Produk gagal beberapa kali, coba produk lain",
	NOT_ELIGIBLE_TRANSACTION:             "Transaksi tidak memenuhi syarat",
	ISSUER_OR_SWITCH_IS_INOPERATIVE:      "Penerbit tidak beroperasi",
	WRONG_CUSTOMER_NUMBER:                "Nomor pelanggan salah",
	WALLET_NOT_ENOUGH:                    "Dompet tidak cukup",
	ADD_FAVOURITE_FAILED:                 "Menambahkan produk ke daftar yang disukai gagal",
	ADD_FAVOURITE_SUCCESS:                "Tambahkan produk ke kesuksesan favorit",
	DELETE_FAVOURITE_SUCCESS:             "Hapus produk untuk kesuksesan favorit",
	PIN_MAX_TRY_COUNT_REACHED:            "Kode PIN 5x Salah",
	PARTIIAL_PAYMENT_FAILED:              "Pembayaran sebagian gagal",
	ALREADY_SUBSCRIBE:                    "Pengguna sudah melakukan langganan",
	CREATE_PASSWORD_ERROR:                "Kata sandi harus mengandung Angka, Huruf Besar, Huruf Kecil, Simbol, dan minimal 8 karakter",
	INPUT_PASSWORD_ERROR:                 "Kata sandi salah",
	TOKEN_USED:                           "Token telah digunakan",
	PVG_EMAIL_NOT_FOUND:                  "E-mail harus berisikan @pvg.co.id",
	EMAIL_ALREADY_REGISTERED:             "E-mail sudah terdaftar",
	PHONE_ALREADY_REGISTERED:             "Nomor telp sudah terdaftar",
	REQUEST_CANNOT_NIL:                   "Permintaan tidak boleh kosong",
	BILL_NOT_FOUND:                       "Tagihan tidak ditemukan",
	PLEASE_PAY_AT_BRANCH_OFFICE:          "Silakan melakukan pembayaran dikantor cabang",
	BILLER_ERROR:                         "Tagihan sedang ada gangguan, silakan ulangi beberapa saat lagi",
	TOO_MANY_REQUEST:                     "Terlalu banyak permintaan, silakan tunggu beberapa saat lagi",
	MAINTENANCE_MODE:                     "Sistem sedang dalam perbaikan",
	CONNECTION_HTTP_ERROR:                "Terjadi kesalahan pada server, silahkan coba lagi",
	CONNECTION_PROVIDER_TIMEOUT:          "Supplier sedang gangguan",
	PRODUCT_NOT_FOUND:                    "Produk tidak ditemukan",
	TRANSACTION_NOT_FOUND:                "Transaksi tidak ditemukan",
	BILLING_NOT_FOUND:                    "No tagihan tidak ditemukan",
	INVALID_TOKEN:                        "Token tidak sesuai",
	SESSION_EXPIRED:                      "Sesi telah berakhir",
	INVALID_REFRESH_TOKEN:                "Token tidak sesuai",
	INVALID_TIMESTAMP:                    "Waktu tidak sesuai",
	NOT_MATCH_DEVICE:                     "Device tidak sesuai",
	IP_BLOCKED:                           "Anda telah diblokir sementara",
	INVALID_SIGNATURE:                    "Signature tidak sesuai",
	VALIDATION_ERROR:                     "Validasi Error",
	REQUIRED_FIELD:                       "Field tidak boleh kosong",
	REQUEST_TO_EARLY:                     "Terlalu cepat melakukan permintaan",
	DATA_DUPLICATED:                      "Duplikasi data",
	MIN_REF_ID:                           "Jumlah digit reference_no kurang",
	MAX_REF_ID:                           "Jumlah digit reference_no sudah melewati batas",
	BALANCE_INSUFFICIENT:                 "Saldo tidak cukup",
	PIN_NOT_MATCH:                        "Pin tidak sesuai",
	USER_HAVE_ACTIVE_PIN:                 "User sudah memiliki pin",
	INVALID_PHONE_NUMBER:                 "Nomor telepon tidak sesuai",
	SQL_SYNTAX_ERROR:                     "Sistem error",
	SQL_ERROR_NO_ROWS:                    "Data tidak ditemukan",
	SQL_CONNECTION_ERROR:                 "Koneksi database error",
	SQL_INSERT_ERROR:                     "Gagal menambahkan data ke database",
	SQL_UPDATE_ERROR:                     "Gagal merubah data ke database",
	SQL_GET_DATA_ERROR:                   "Gagal mengambil data dari database",
	REDIS_NIL:                            "Data tidak ditemukan",
	REDIS_CONNECTION_ERROR:               "Koneksi redis error",
	KAFKA_DIAL_ERROR:                     "Sistem error",
	KAFKA_CONNECTION_ERROR:               "Koneksi error",
	KAFKA_PUBLISH_ERROR:                  "Sistem error",
	NSQ_DIAL_ERROR:                       "Sistem error",
	NSQ_CONNECTION_ERROR:                 "Koneksi error",
	NSQ_PUBLISH_ERROR:                    "Sistem error",
	INVOICE_AMOUNT_DOES_NOT_MATCH:        "Nominal tagihan tidak seusai, silahkan masukan tagihan yang sesuai",
	INQUIRY_GENERAL_ERROR:                "Inquiry gagal, silahkan coba beberapa saat lagi",
	MAX_AMOUNT_VALIDATION_CHECK:          "Jumlah pembayaran harus lebih kecil dari jumlah maksimum plus biaya admin (jika ada biaya admin)",
	MIN_AMOUNT_VALIDATION_CHECK:          "Jumlah pembayaran harus lebih besar dari jumlah minimum plus biaya admin (jika ada biaya admin)",
	MUST_BE_PAID_AT_THE_PROVIDER_LOCKET:  "Silakan melakukan pembayaran di loket",
	TRANSACTION_REPEATED:                 "Pembelian ulang pada product yang sama",
	PASSWORD_ALREADY_SET:                 "Password sudah ada",
	INVALID_ROLE_ID:                      "Tidak dapat menemukan role yang dipilih",
	DUPLICATE_REF_ID:                     "ReferenceNo sudah digunakan, silakan gunakan reference_no yang lain",
	DATA_TOO_LARGE:                       "Permintaan terlalu besar, silakan periksa permintaan kamu",
	UNIMPLEMENTED:                        "Belum di implementasi",
	USER_NOT_SET:                         "data user belum di set",
	INVALID_IP:                           "Permintaan tidak sesuai",
	INVALID_DEVICE_ID:                    "Device tidak sesuai",
	USER_TYPE_IS_NOT_B2B:                 "user type bukan B2B",
	INVALID_UUID:                         "uuid tidak sesuai",
	INVALID_GRANT_TYPE:                   "grant type tidak sesuai",
	INVALID_CREDENTIALS:                  "kredensial tidak sesuai",
	IP_NOT_WHITELISTED:                   "IP tidak ada di dalam whitelist",
	USER_TYPE_IS_NOT_IRS:                 "user type bukan irs",
	PARSE_TIME_FAILED:                    "Format waktu tidak sesuai",
	DELETE_OWN_USER:                      "tidak dapat menghapus akun sendiri",
	DELETE_OWNER_USER:                    "tidak dapat menghapus akun dengan role owner",
	LAST_INVITATION_LESS_THAN_24_HOURS:   "undangan terakhir telah dikirim kurang dari 24 jam lalu. tidak dapat mengirim undangan ulang.",
	AMOUNT_MORE_THAN_SELLING_PRICE:       "Harga jual harus lebih besar dari harga produk",
	EMPTY_FILE:                           "Tidak dapat memproses file kosong",
	SAME_PASSWORD_AS_BEFORE:              "Password sama seperti sebelumnya",
	INVALID_FORMAT_PUBLIC_KEY:            "Format public key tidak sesuai",
	INVALID_FORMAT_URL:                   "Format URL tidak sesuai",
	FAILED_SENT_EMAIL:                    "Email gagal terkirim",
	FRAUD_NUMBER:                         "Nomor anda terindikasi fraud. Silahkan hubungi admin",
	PHONE_NUMBER_BLOCKED:                 "Akun Diblokir",
	PHONE_NUMBER_NONACTIVE:               "Akun Dinonaktifkan",
	MAX_BALANCE_EXCEEDED:                 "Maksimal saldo %s",
	DATE_RANGE_EXCEEDS_LIMIT:             "Anda melebihi batas rentang waktu yang telah ditentukan",
	EMAIL_NOT_FOUND:                      "Email tidak ditemukan",
	PASSWORD_NOT_MATCH:                   "Password tidak sama",
	SEAT_BOOKED:                          "Kursi tidak tersedia, silakan pilih kursi lain",
	TRANSACTION_FAILED_TRY_AGAIN:         "Transaksi gagal, silakan coba beberpa saat lagi",
	EDIT_OWNER_USER:                      "tidak dapat melakukan edit terhadap user owner",
	EDIT_OWN_USER:                        "tidak dapat melakukan edit terhadap akun anda sendiri",
	EDIT_TO_OWNER_ROLES:                  "tidak dapat mengubah role menjadi owner",
	EMAIL_DOMAIN_NOT_MATCH_WITH_OWNER:    "domain e-mail harus sama dengan owner",
	BOOKING_TICKET_EXPIRED:               "Kode tiket kursi sudah kadaluarsa",
	CHANNEL_TOPUP_NOT_ALOWED:             "channel tidak dapat melakukan topup",
	TOPUP_EXPIRED:                        "transaksi topup ini telah kadaluarsa",
	FAIL_TO_UPSERT_BULK:                  "gagal melakukan bulk upsert",
	PAYMENT_CHANNEL_NOT_IMPLEMENTED:      "Metode pembayaran tidak diimplementasikan untuk channel ini",
	PAYMENT_CHANNEL_IS_DISABLE:           "Metode pembayaran tidak aktif",
	PAYMENT_SOURCE_NOT_IMPLEMENTED:       "Sumber pembayaran tidak diketahui",
	TOTAL_B2B_USERS_EXCEEDS_THE_LIMIT:    "Jumlah user telah melebihi batas yang diberikan",
	PAYMENT_STATUS_IS_FINAL:              "Pembayaran tidak bisa dibatalkan",
	CANCEL_PAYMENT_FAILED:                "Gagal membatalkan pembayaran",
	TOPUP_STATUS_IS_FINAL:                "Status top up tidak bisa diubah",
	USER_HAS_TOPUP_ACTIVE:                "Anda memiliki top up yang aktif, silakan tunggu hingga proses transaksi selesai atau batalkan top up sebelumnya terlebih dahulu",
	MAXIMUM_KTP_USE:                      "Nomor KTP sudah terdaftar",
	INVALID_INTEGER_VARIABLE:             "tipe data tidak sesuai",
	USER_DONT_HAVE_TOPUP_ACTIVE:          "Pengguna tidak memiliki top up aktif",
	AMOUNT_MISMATCH:                      "Jumlah pembayaran tidak sesuai",
	SAME_EXCHANGE_CURRENCY:               "currency from dan currency to tidak boleh sama",
	INVALID_QRIS:                         "Format QRIS tidak sesuai",
	TOKEN_ALREADY_EXISTS:                 "token telah digunakan",
	INVALID_DATE_FORMAT:                  "Format tanggal tidak sesuai",
	RESOURCE_NOT_FOUND:                   "Resource tidak ditemukan",
	GROUP_BUY_FULL:                       "Grup full",
	GROUP_BUY_NOT_FOUND:                  "Grup belanja tidak ditemukan",
	GROUP_BUY_EXPIRED:                    "Grup belanja kedaluwarsa",
	GROUP_BUY_NOT_AVAILABLE:              "Grup belanja tidak tersedia",
	GROUP_BUY_CLOSED:                     "Grup belanja ditutup",
	GROUP_BUY_ALREADY_JOINED:             "Anda sudah bergabung dalam grup belanja",
	CANNOT_CANCEL_ORDER:                  "Tidak dapat membatalkan pembelian",
	PAYMENT_METHOD_NOT_ALLOWED:           "Metode pembayaran tidak diizinkan untuk produk ini",
	INVALID_REJECT_REASON:                "Invalid Reject Reason",
	INVALID_STATUS:                       "Invalid status",
	INVALID_CONTENT_TYPE:                 "tipe konten tidak sesuai",
	PRODUCT_IS_IN_MAINTENANCE:            "Sedang dalam perbaikan",
	INVALID_AMOUNT:                       "Jumlah Pembayaran harus merupakan kelipatan dari %s",
	NUMBER_ONLY_AMOUNT:                   "Amount must be in numbers only",
	MAX_BEST_SELLER_DATA:                 "Maksimal hanya dapat memilih 3 produk best seller dengan kategori dan prefix yang sama",
	INVALID_CATEGORY_BEST_SELLER:         "category tidak bisa dijadikan best seller",
	PAYMENT_METHOD_MIN_VERSION:           "Metode pembayaran tidak dapat digunakan",
	PRODUCT_INQUIRY_NOT_ALLOWED:          "Inquiry untuk produk ini belum tersedia",
	PROVIDER_PRODUCT_MAINTENANCE:         "Produk sedang dalam perbaikan",
	PROVIDER_PRODUCT_CLOSED:              "Produk sudah tidak dijual",
	MINIMUM_AMOUNT:                       "Nominal tidak sesuai",
	MAXIMUM_AMOUNT:                       "Nominal tidak sesuai",
	INVALID_MULTIPLES_AMOUNT:             "Nominal tidak sesuai",
	CONTEXT_DEADLINE_EXCEEDED:            "Request timeout",
	CONTEXT_CANCELLED:                    "Request canceled",
	OTP_EXPIRED:                          "OTP telah kedaluwarsa",
}

var mapCodeErrMoreInfo = map[CodeErr]string{
	PIN_MAX_TRY_COUNT_REACHED:  "Anda melebihi limit verifikasi kode PIN. Mohon tunggu 30 menit untuk mencoba verifikasi kembali",
	SAME_PASSWORD_AS_BEFORE:    "Gunakan password yang berbeda dengan sebelumnya",
	INVALID_FORMAT_URL:         "Format URL harus diawali dengan HTTP atau HTTPS",
	PHONE_NUMBER_BLOCKED:       "Jika ada kendala atau pertanyaan, silakan hubungi tim CS kami",
	OTP_NUMBER_BLOCKED:         "Jika ada kendala atau pertanyaan, silakan hubungi tim CS kami",
	PHONE_NUMBER_NONACTIVE:     "Untuk mengaktifkannya kembali, mohon hubungi tim CS kami",
	PAYMENT_STATUS_IS_FINAL:    "Status pembayaran sudah final, silakan tunggu hingga proses transaksi selesai",
	CANCEL_PAYMENT_FAILED:      "Silakan tunggu hingga proses transaksi selesai untuk dapat membatalkan pembayaran",
	TOPUP_STATUS_IS_FINAL:      "Silakan hubungi team CS jika ingin melakukan perubahan status top up",
	GROUP_BUY_FULL:             "Grup yang Anda coba bergabung sudah mencapai kapasitas maksimal",
	GROUP_BUY_NOT_FOUND:        "Grup belanja yang Anda cari tidak ditemukan. Cek grup lainnya yang tersedia untuk bergabung",
	GROUP_BUY_EXPIRED:          "Grup belanja sudah kedaluwarsa. Cek grup lainnya",
	GROUP_BUY_NOT_AVAILABLE:    "Grup belanja sudah tidak tersedia. Cek grup lainnya yang tersedia untuk bergabung",
	GROUP_BUY_CLOSED:           "Grup belanja sudah ditutup dan tidak menerima anggota baru. Cek grup lainnya",
	GROUP_BUY_ALREADY_JOINED:   "Silakan duduk manis dan tunggu kuota grup belanja terpenuhi. Yuk gabung ke grup belanja lainnya",
	CANNOT_CANCEL_ORDER:        "Maaf, pembelian kamu tidak dapat dibatalkan. Jika ingin tetap dibatalkan, silakan hubungi CS kami",
	PRODUCT_IS_IN_MAINTENANCE:  "Klikoo sedang dalam perbaikan sistem. Semua transaksi sementara ini tidak dapat dilakukan",
	PAYMENT_METHOD_MIN_VERSION: "Minimal versi aplikasi untuk metode pembayaran '%s' adalah v%s",
	MINIMUM_AMOUNT:             "Nominal minimal %s",
	MAXIMUM_AMOUNT:             "Nominal maksimal %s",
	INVALID_MULTIPLES_AMOUNT:   "Nominal harus kelipatan dari %s",
	OTP_EXPIRED:                "Silakan lakukan request OTP ulang",
}

func (ce CodeErr) Error() string {
	return fmt.Sprint(strings.ToLower(ce.Message()))
}

func (ce CodeErr) Errors() error {
	return NewError(errors.New(strings.ToLower(ce.Message())), ce.StatusCode(), ce.Code(), ce.Message())
}

func (ce CodeErr) Code() int {
	return int(ce)
}

func (ce CodeErr) StatusCode() int {
	val, ok := mapCodeErrStatusCode[ce]
	if !ok {
		return http.StatusInternalServerError
	}

	return val
}

func (ce CodeErr) Message() string {
	val, ok := mapCodeErrMessage[ce]
	if !ok {
		return InternalServerError
	}

	return val
}

func (ce CodeErr) MoreInfo(moreInfos ...string) string {
	if len(moreInfos) != 0 {
		return strings.Join(moreInfos, ", ")
	}

	val, ok := mapCodeErrMoreInfo[ce]
	if !ok {
		return ""
	}

	return val
}

func (ce CodeErr) MoreInfoFormatted(format ...interface{}) string {
	val, ok := mapCodeErrMoreInfo[ce]
	if !ok {
		return ""
	}

	return fmt.Sprintf(val, format...)
}
