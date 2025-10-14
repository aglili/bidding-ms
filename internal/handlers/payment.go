package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/utils"
	"github.com/gin-gonic/gin"
)




type PaymentHandler struct{
	cfg *config.Config
}




func NewPaymentHandler(cfg *config.Config) *PaymentHandler{
	return  &PaymentHandler{cfg: cfg}
}





func (h *PaymentHandler) WebhookEndpoint(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.RespondWithError(ctx,err,"failed to read request body")
		return
	}

	signature := ctx.GetHeader("x-paystack-signature")
	if signature == ""{
		utils.RespondWithError(ctx,err,"missing signature header")
		return
	}

	if !verifyPaystackSignature(h.cfg.PaystackSecretKey,body,signature){
		utils.RespondWithError(ctx,fmt.Errorf("invalid signature"),"invalid signature")
		return
	}

	var event map[string]any
	if err := json.Unmarshal(body,&event);err != nil {
		utils.RespondWithError(ctx,err,"failed to parse request json")
		return
	}


	eventType := event["event"].(string)
	log.Printf("event of type:[%v] has been received",eventType)

	switch eventType{
	case "charge.success":
		data := event["data"].(map[string]any)
		log.Printf("data received: %v",data)
		// TODO: handle payments later with an event
	default:
		fmt.Println("Unhandled event:", eventType)
	}



	ctx.Status(http.StatusOK)
	}









func verifyPaystackSignature(secretKey string, body []byte,signature string) bool  {
	mac := hmac.New(sha512.New,[]byte(secretKey))
	mac.Write(body)
	expectedMac := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMac)
	return  hmac.Equal([]byte(signature),[]byte(expectedSignature))
	
}