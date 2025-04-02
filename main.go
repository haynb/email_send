package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	IsHtml  bool     `json:"is_html,omitempty"`
}

type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type EmailConfig struct {
	SMTPServer   string
	SMTPPort     int
	SenderEmail  string
	SenderName   string
	AuthPassword string
}

func sendEmail(config EmailConfig, req EmailRequest) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.SenderEmail, config.SenderName))
	m.SetHeader("To", req.To...)
	m.SetHeader("Subject", req.Subject)

	if req.IsHtml {
		m.SetBody("text/html", req.Body)
	} else {
		m.SetBody("text/plain", req.Body)
	}

	d := gomail.NewDialer(config.SMTPServer, config.SMTPPort, config.SenderEmail, config.AuthPassword)

	return d.DialAndSend(m)
}

func main() {
	// 获取环境变量
	smtpServer := getEnv("SMTP_SERVER", "smtp.163.com")
	smtpPortStr := getEnv("SMTP_PORT", "465")
	smtpPort, _ := strconv.Atoi(smtpPortStr)
	senderEmail := getEnv("SENDER_EMAIL", "")
	senderName := getEnv("SENDER_NAME", "Alert Service")
	authPassword := getEnv("AUTH_PASSWORD", "")

	if senderEmail == "" {
		log.Fatal("环境变量 SENDER_EMAIL 必须设置")
	}

	if authPassword == "" {
		log.Fatal("环境变量 AUTH_PASSWORD 必须设置")
	}

	config := EmailConfig{
		SMTPServer:   smtpServer,
		SMTPPort:     smtpPort,
		SenderEmail:  senderEmail,
		SenderName:   senderName,
		AuthPassword: authPassword,
	}

	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 发送邮件的API端点
	http.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
			return
		}

		var req EmailRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			sendErrorResponse(w, "无效的请求格式: "+err.Error())
			return
		}

		// 验证请求
		if len(req.To) == 0 {
			sendErrorResponse(w, "收件人不能为空")
			return
		}

		if req.Subject == "" {
			sendErrorResponse(w, "邮件主题不能为空")
			return
		}

		if req.Body == "" {
			sendErrorResponse(w, "邮件内容不能为空")
			return
		}

		// 发送邮件
		err := sendEmail(config, req)
		if err != nil {
			log.Printf("发送邮件失败: %v", err)
			sendErrorResponse(w, "发送邮件失败: "+err.Error())
			return
		}

		// 发送成功响应
		resp := EmailResponse{
			Success: true,
			Message: "邮件发送成功",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("邮件服务启动，监听端口: %s", port)
	log.Fatal(server.ListenAndServe())
}

func sendErrorResponse(w http.ResponseWriter, message string) {
	resp := EmailResponse{
		Success: false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
