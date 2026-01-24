package middleware

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"jasamc/infrastructure/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// =========================
		// Generate Request ID
		// =========================
		requestID := uuid.New().String()

		// ⛔ PENTING: pakai context dari request asli
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		ctx := context.WithValue(timeoutCtx, "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)

		// =========================
		// Payload Logging (SAFE)
		// =========================
		var payload string
		method := c.Request.Method
		contentType := c.GetHeader("Content-Type")

		if method == "POST" || method == "PUT" || method == "PATCH" {

			// ❌ JANGAN LOG FILE / MULTIPART
			if strings.HasPrefix(contentType, "multipart/form-data") {
				payload = "[multipart/form-data hidden]"
			} else {

				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					payload = "[error reading body]"
				} else {

					// Batasi payload max 2KB
					if len(bodyBytes) > 2048 {
						payload = string(bodyBytes[:2048]) + "...(truncated)"
					} else {
						payload = string(bodyBytes)
					}
				}

				// WAJIB rewind body
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// =========================
		// Execute Request
		// =========================
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// =========================
		// Build Log
		// =========================
		logFields := logrus.Fields{
			"request_id": requestID,
			"method":     method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency":    latency.String(),
			"payload":    payload,
		}

		// =========================
		// Log Result
		// =========================
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			log.Logger.WithFields(logFields).Info("Request success")
		} else {
			log.Logger.WithFields(logFields).Warn("Request error")
		}
	}
}
