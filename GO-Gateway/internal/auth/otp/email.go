package otp

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// SendEmail ส่ง OTP ผ่าน Gmail SMTP แบบสวยงาม
func SendEmail(to string, code string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	// สร้าง HTML Template สำหรับ OTP Email
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f4f4;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 10px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
        
        <!-- Header -->
        <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 600;">🔐 OTP Verification</h1>
            <p style="color: #e8eaff; margin: 10px 0 0 0; font-size: 16px;">Secure Account Verification</p>
        </div>
        
        <!-- Content -->
        <div style="padding: 40px 30px;">
            <h2 style="color: #333333; margin: 0 0 20px 0; font-size: 24px; font-weight: 600;">Hello there! 👋</h2>
            
            <p style="color: #666666; line-height: 1.6; margin: 0 0 25px 0; font-size: 16px;">
                We received a request to verify your account. Please use the OTP code below to proceed with your verification.
            </p>
            
            <!-- OTP Code Box -->
            <div style="text-align: center; margin: 35px 0;">
                <div style="display: inline-block; background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 50%, #fecfef 100%); padding: 3px; border-radius: 15px;">
                    <div style="background-color: #ffffff; padding: 25px 40px; border-radius: 12px;">
                        <p style="margin: 0 0 10px 0; color: #666666; font-size: 14px; font-weight: 500;">Your OTP Code</p>
                        <div style="font-size: 36px; font-weight: bold; color: #4a5568; letter-spacing: 8px; font-family: 'Courier New', monospace;">{{OTP_CODE}}</div>
                    </div>
                </div>
            </div>
            
            <!-- Warning Box -->
            <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 8px; padding: 20px; margin: 30px 0;">
                <div style="display: flex; align-items: flex-start;">
                    <div style="margin-right: 15px; font-size: 24px;">⚠️</div>
                    <div>
                        <p style="margin: 0 0 10px 0; color: #856404; font-weight: 600; font-size: 16px;">Important Security Notice:</p>
                        <ul style="margin: 0; padding-left: 20px; color: #856404; line-height: 1.6;">
                            <li>This OTP code will expire in <strong>5 minutes</strong></li>
                            <li>Never share this code with anyone</li>
                            <li>If you didn't request this verification, please ignore this email</li>
                        </ul>
                    </div>
                </div>
            </div>
            
            <p style="color: #666666; line-height: 1.6; margin: 25px 0 0 0; font-size: 14px;">
                If you have any questions or need assistance, feel free to contact our support team anytime.
            </p>
        </div>
    </div>
    
    <!-- Mobile Responsive -->
    <style>
        @media only screen and (max-width: 600px) {
            .container { width: 100% !important; }
            .content { padding: 20px !important; }
            .otp-code { font-size: 28px !important; letter-spacing: 4px !important; }
        }
    </style>
</body>
</html>`

	// แทนที่ OTP Code ใน Template
	htmlBody := strings.Replace(htmlTemplate, "{{OTP_CODE}}", code, 1)

	// สร้าง Email Message แบบ Multipart
	boundary := "boundary123456789"
	
	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: =?UTF-8?Q?🔐_Your_OTP_Verification_Code?=\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=\"%s\"\r\n"+
		"\r\n"+
		"--%s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: 8bit\r\n"+
		"\r\n"+
		"Your OTP code is: %s\r\n"+
		"\r\n"+
		"This code will expire in 5 minutes.\r\n"+
		"Please do not share this code with anyone.\r\n"+
		"\r\n"+
		"--%s\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: 8bit\r\n"+
		"\r\n"+
		"%s\r\n"+
		"--%s--\r\n",
		to, boundary, boundary, code, boundary, htmlBody, boundary)

	// ตั้งค่า SMTP Auth
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	// ส่งอีเมลผ่าน Gmail SMTP Server
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log ว่าส่ง OTP ไปที่อีเมลแล้ว
	log.Printf("OTP email sent to %s successfully", to)
	return nil
}