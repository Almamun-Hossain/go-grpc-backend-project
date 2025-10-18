package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	serviceURL string
}

func NewMembershipHandler(serviceURL string) *MembershipHandler {
	return &MembershipHandler{
		serviceURL: serviceURL,
	}
}

func (h *MembershipHandler) Register(c *gin.Context) {
	h.proxyRequest(c, "POST", "/v1/auth/register")
}

func (h *MembershipHandler) Login(c *gin.Context) {
	h.proxyRequest(c, "POST", "/v1/auth/login")
}

func (h *MembershipHandler) RefreshToken(c *gin.Context) {
	h.proxyRequest(c, "POST", "/v1/auth/refresh")
}

func (h *MembershipHandler) GetProfile(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/users/profile")
}

func (h *MembershipHandler) UpdateProfile(c *gin.Context) {
	h.proxyRequest(c, "PUT", "/v1/users/profile")
}

func (h *MembershipHandler) GetAssets(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/users/assets")
}

func (h *MembershipHandler) UpdateAsset(c *gin.Context) {
	assetID := c.Param("asset_id")
	h.proxyRequest(c, "PUT", "/v1/users/assets/"+assetID)
}

func (h *MembershipHandler) GetReferrals(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/users/referrals")
}

func (h *MembershipHandler) RequestWithdrawal(c *gin.Context) {
	h.proxyRequest(c, "POST", "/v1/users/withdrawals")
}

func (h *MembershipHandler) GetWithdrawals(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/users/withdrawals")
}

func (h *MembershipHandler) proxyRequest(c *gin.Context, method, path string) {
	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Create request to backend service
	url := h.serviceURL + path
	if c.Request.URL.RawQuery != "" {
		url += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	req.Header = c.Request.Header.Clone()

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach membership service"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Forward response
	var jsonResp interface{}
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
		return
	}

	c.JSON(resp.StatusCode, jsonResp)
}

