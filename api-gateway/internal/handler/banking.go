package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BankingHandler struct {
	serviceURL string
}

func NewBankingHandler(serviceURL string) *BankingHandler {
	return &BankingHandler{
		serviceURL: serviceURL,
	}
}

func (h *BankingHandler) GetCardInfo(c *gin.Context) {
	cardID := c.Param("card_id")
	h.proxyRequest(c, "GET", "/v1/banking/cards/"+cardID)
}

func (h *BankingHandler) ListCards(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/banking/cards")
}

func (h *BankingHandler) GetCreditCardBalance(c *gin.Context) {
	cardID := c.Param("card_id")
	h.proxyRequest(c, "GET", "/v1/banking/cards/"+cardID+"/balance")
}

func (h *BankingHandler) GetPhysicalCardDetails(c *gin.Context) {
	cardID := c.Param("card_id")
	h.proxyRequest(c, "GET", "/v1/banking/cards/"+cardID+"/physical")
}

func (h *BankingHandler) InitiateRemittance(c *gin.Context) {
	h.proxyRequest(c, "POST", "/v1/banking/remittance")
}

func (h *BankingHandler) GetRemittanceStatus(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	h.proxyRequest(c, "GET", "/v1/banking/remittance/"+transactionID)
}

func (h *BankingHandler) ListRemittances(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/banking/remittance")
}

func (h *BankingHandler) GetExchangeRate(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/banking/exchange-rate")
}

func (h *BankingHandler) proxyRequest(c *gin.Context, method, path string) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	url := h.serviceURL + path
	if c.Request.URL.RawQuery != "" {
		url += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach banking service"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var jsonResp interface{}
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
		return
	}

	c.JSON(resp.StatusCode, jsonResp)
}

