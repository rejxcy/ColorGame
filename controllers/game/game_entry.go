package game

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func GameEntry(ctx *gin.Context) {
	localIP, err := GetLocalIP()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error getting local IP")
		return
	}

	url := fmt.Sprintf("http://%s:8080/player_name", localIP)

	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error creating QR code")
		return
	}

	qrCode.DisableBorder = true

	qrCodeImage, err := qrCode.PNG(256)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error encoding QR code to PNG")
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeImage)

	ctx.HTML(http.StatusOK, "game_entry.html", gin.H{
		"QRCodeBase64": template.HTML(qrCodeBase64),
	})
}
