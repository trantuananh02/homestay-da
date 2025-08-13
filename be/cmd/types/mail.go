package types

// VerificationEmailData dữ liệu email xác nhận tài khoản
type VerificationEmailData struct {
	Name            string `json:"name"`
	VerificationLink string `json:"verification_link"`
}
