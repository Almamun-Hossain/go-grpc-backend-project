package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebsiteHandler struct {
	serviceURL string
}

func NewWebsiteHandler(serviceURL string) *WebsiteHandler {
	return &WebsiteHandler{
		serviceURL: serviceURL,
	}
}

func (h *WebsiteHandler) GetSiteSettings(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/site/settings")
}

func (h *WebsiteHandler) UpdateSiteSettings(c *gin.Context) {
	h.proxyRequest(c, "PUT", "/v1/site/settings")
}

func (h *WebsiteHandler) UpdateSiteTitle(c *gin.Context) {
	h.proxyRequest(c, "PUT", "/v1/site/title")
}

func (h *WebsiteHandler) UploadLogo(c *gin.Context) {
	h.proxyRequest(c, "POST", "/v1/site/logo")
}

func (h *WebsiteHandler) GetLogo(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/site/logo")
}

func (h *WebsiteHandler) DeleteLogo(c *gin.Context) {
	h.proxyRequest(c, "DELETE", "/v1/site/logo")
}

func (h *WebsiteHandler) GetDynamicSettings(c *gin.Context) {
	h.proxyRequest(c, "GET", "/v1/site/dynamic-settings")
}

func (h *WebsiteHandler) UpdateDynamicSetting(c *gin.Context) {
	key := c.Param("key")
	h.proxyRequest(c, "PUT", "/v1/site/dynamic-settings/"+key)
}

func (h *WebsiteHandler) DeleteDynamicSetting(c *gin.Context) {
	key := c.Param("key")
	h.proxyRequest(c, "DELETE", "/v1/site/dynamic-settings/"+key)
}

func (h *WebsiteHandler) proxyRequest(c *gin.Context, method, path string) {
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
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach website service"})
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

