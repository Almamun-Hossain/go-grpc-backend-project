package router

import (
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupMembershipRoutes sets up routes for membership service
func SetupMembershipRoutes(r *gin.Engine, h *handler.MembershipHandler) {
	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		// Protected routes
		users := v1.Group("/users")
		users.Use(middleware.Auth())
		{
			users.GET("/profile", h.GetProfile)
			users.PUT("/profile", h.UpdateProfile)
			users.GET("/assets", h.GetAssets)
			users.PUT("/assets/:asset_id", h.UpdateAsset)
			users.GET("/referrals", h.GetReferrals)
			users.POST("/withdrawals", h.RequestWithdrawal)
			users.GET("/withdrawals", h.GetWithdrawals)
		}
	}
}

// SetupWebsiteRoutes sets up routes for website service
func SetupWebsiteRoutes(r *gin.Engine, h *handler.WebsiteHandler) {
	v1 := r.Group("/api/v1")
	{
		site := v1.Group("/site")
		{
			// Public routes
			site.GET("/settings", h.GetSiteSettings)
			site.GET("/logo", h.GetLogo)
			site.GET("/dynamic-settings", h.GetDynamicSettings)

			// Protected routes (admin only)
			admin := site.Group("")
			admin.Use(middleware.Auth()) // TODO: Add admin role check
			{
				admin.PUT("/settings", h.UpdateSiteSettings)
				admin.PUT("/title", h.UpdateSiteTitle)
				admin.POST("/logo", h.UploadLogo)
				admin.DELETE("/logo", h.DeleteLogo)
				admin.PUT("/dynamic-settings/:key", h.UpdateDynamicSetting)
				admin.DELETE("/dynamic-settings/:key", h.DeleteDynamicSetting)
			}
		}
	}
}

// SetupBankingRoutes sets up routes for banking service
func SetupBankingRoutes(r *gin.Engine, h *handler.BankingHandler) {
	v1 := r.Group("/api/v1")
	v1.Use(middleware.Auth())
	{
		banking := v1.Group("/banking")
		{
			// Card routes
			banking.GET("/cards", h.ListCards)
			banking.GET("/cards/:card_id", h.GetCardInfo)
			banking.GET("/cards/:card_id/balance", h.GetCreditCardBalance)
			banking.GET("/cards/:card_id/physical", h.GetPhysicalCardDetails)

			// Remittance routes
			banking.POST("/remittance", h.InitiateRemittance)
			banking.GET("/remittance", h.ListRemittances)
			banking.GET("/remittance/:transaction_id", h.GetRemittanceStatus)

			// Exchange rate
			banking.GET("/exchange-rate", h.GetExchangeRate)
		}
	}
}

