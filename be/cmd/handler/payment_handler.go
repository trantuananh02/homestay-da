package handler

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	vnpTmnCode    = "DBZCRTPR"
	vnpHashSecret = "EHL12KYWUG52L5NM8KHVT7LD3MYXZESG"
	vnpUrl        = "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"
	vnpReturnUrl  = "http://localhost:3000/payment-success"
)

func CreateVnpayPayment(c *gin.Context) {
	amount := c.Query("amount") // Số tiền (VND)
	orderId := c.Query("orderId") // Mã đơn hàng
	orderInfo := c.Query("orderInfo") // Thông tin đơn hàng

	if amount == "" || orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing amount or orderId"})
		return
	}

	// Đảm bảo số tiền là số nguyên, nhân 100, không có ký tự lạ
	var vnpAmount string
	{
		var amtInt int64
		_, err := fmt.Sscanf(amount, "%d", &amtInt)
		if err != nil || amtInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Số tiền không hợp lệ"})
			return
		}
		vnpAmount = fmt.Sprintf("%d", amtInt*100)
	}

	vnpParams := map[string]string{
		"vnp_Version":     "2.1.0",
		"vnp_Command":     "pay",
		"vnp_TmnCode":     vnpTmnCode,
		"vnp_Amount":      vnpAmount,
		"vnp_CurrCode":    "VND",
		"vnp_TxnRef":      orderId,
		"vnp_OrderInfo":   orderInfo,
		"vnp_OrderType":   "other",
		"vnp_Locale":      "vn",
		"vnp_ReturnUrl":   vnpReturnUrl,
		"vnp_IpAddr":      c.ClientIP(),
		"vnp_CreateDate":  time.Now().Format("20060102150405"),
	}

	// Sắp xếp key
	var keys []string
	for k := range vnpParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Tạo chuỗi dữ liệu ký đúng chuẩn VNPAY
	var signData string
	for i, k := range keys {
		if i > 0 {
			signData += "&"
		}
		signData += fmt.Sprintf("%s=%s", k, vnpParams[k])
	}

	// Tạo query string cho URL (dùng url.Values, không encode giá trị)
	var query url.Values = url.Values{}
	for _, k := range keys {
		query.Add(k, vnpParams[k])
	}

	// Tạo checksum
	h := hmac.New(sha512.New, []byte(vnpHashSecret))
	h.Write([]byte(signData))
	vnpSecureHash := hex.EncodeToString(h.Sum(nil))
	query.Add("vnp_SecureHash", vnpSecureHash)

	paymentUrl := vnpUrl + "?" + query.Encode()

	c.JSON(http.StatusOK, gin.H{
		"paymentUrl": paymentUrl,
	})
}
