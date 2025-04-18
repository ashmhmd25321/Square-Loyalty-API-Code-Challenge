package handlers

import (
	"fmt"
	"net/http"
	"time"

	"loyalty-app/internal/models"
	"loyalty-app/internal/services"

	"github.com/gin-gonic/gin"
)

type LoyaltyHandler struct {
	squareService *services.SquareService
	dbService     *services.MemoryService
}

func NewLoyaltyHandler(squareService *services.SquareService, dbService *services.MemoryService) *LoyaltyHandler {
	return &LoyaltyHandler{
		squareService: squareService,
		dbService:     dbService,
	}
}

// EarnPoints handles the request to earn loyalty points
func (h *LoyaltyHandler) EarnPoints(c *gin.Context) {
	var request struct {
		Points int `json:"points" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user by ID
	user, err := h.dbService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Get or create Square customer for this user
	customer, err := h.squareService.GetOrCreateCustomer(c.Request.Context(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Square customer: %v", err)})
		return
	}

	// Get or create loyalty account for this customer
	loyaltyAccount, err := h.squareService.GetOrCreateLoyaltyAccount(c.Request.Context(), customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create loyalty account: %v", err)})
		return
	}

	// Earn points in Square
	if err := h.squareService.EarnPoints(c.Request.Context(), loyaltyAccount.ID, request.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to earn points in Square: %v", err)})
		return
	}

	// Get account ID from user for our local database
	account, err := h.dbService.GetLoyaltyAccountByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get local loyalty account"})
		return
	}

	// Record the transaction in our local database
	accountIDStr := fmt.Sprintf("%d", account.ID)
	if err := h.dbService.AddTransaction(c.Request.Context(), accountIDStr, "EARN", request.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record points locally"})
		return
	}

	// Update Square loyalty account ID in our database for future reference
	if err := h.dbService.UpdateSquareLoyaltyID(c.Request.Context(), accountIDStr, loyaltyAccount.ID); err != nil {
		// This is not critical, so just log it
		fmt.Printf("Warning: Failed to update Square loyalty ID: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Points earned successfully",
		"square_account_id": loyaltyAccount.ID,
	})
}

// RedeemPoints handles the request to redeem loyalty points
func (h *LoyaltyHandler) RedeemPoints(c *gin.Context) {
	var request struct {
		Points int `json:"points" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user by ID
	user, err := h.dbService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Get account ID from user for our local database
	account, err := h.dbService.GetLoyaltyAccountByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get local loyalty account"})
		return
	}
	accountIDStr := fmt.Sprintf("%d", account.ID)

	// Check if user has enough points locally - convert uint to string
	balance, err := h.dbService.GetBalance(c.Request.Context(), accountIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}

	if balance < request.Points {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough points"})
		return
	}

	// Get the Square loyalty account ID
	squareLoyaltyID, err := h.dbService.GetSquareLoyaltyID(c.Request.Context(), accountIDStr)
	if err != nil || squareLoyaltyID == "" {
		// If we don't have a Square account ID stored, try to get or create one
		customer, err := h.squareService.GetOrCreateCustomer(c.Request.Context(), user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Square customer: %v", err)})
			return
		}

		loyaltyAccount, err := h.squareService.GetOrCreateLoyaltyAccount(c.Request.Context(), customer.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create loyalty account: %v", err)})
			return
		}
		squareLoyaltyID = loyaltyAccount.ID

		// Update Square loyalty account ID in our database
		if err := h.dbService.UpdateSquareLoyaltyID(c.Request.Context(), accountIDStr, squareLoyaltyID); err != nil {
			fmt.Printf("Warning: Failed to update Square loyalty ID: %v\n", err)
		}
	}

	// Redeem points in Square
	if err := h.squareService.RedeemPoints(c.Request.Context(), squareLoyaltyID, request.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to redeem points in Square: %v", err)})
		return
	}

	// Record the redemption in our database
	if err := h.dbService.AddTransaction(c.Request.Context(), accountIDStr, "REDEEM", request.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record redemption locally"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Points redeemed successfully",
		"square_account_id": squareLoyaltyID,
	})
}

// GetBalance handles the request to get the current points balance
func (h *LoyaltyHandler) GetBalance(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get account ID from user
	account, err := h.dbService.GetLoyaltyAccountByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loyalty account"})
		return
	}

	// Get the balance from our database - convert uint to string
	accountIDStr := fmt.Sprintf("%d", account.ID)
	balance, err := h.dbService.GetBalance(c.Request.Context(), accountIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}

	// Try to get Square balance
	squareLoyaltyID, err := h.dbService.GetSquareLoyaltyID(c.Request.Context(), accountIDStr)
	if err == nil && squareLoyaltyID != "" {
		// If we have a Square account ID, get balance from Square
		squareBalance, err := h.squareService.GetBalance(c.Request.Context(), squareLoyaltyID)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"balance":           balance,
				"square_balance":    squareBalance,
				"square_account_id": squareLoyaltyID,
			})
			return
		}
	}

	// Fallback to just local balance
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// GetHistory handles the request to get the transaction history
func (h *LoyaltyHandler) GetHistory(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get account ID from user
	account, err := h.dbService.GetLoyaltyAccountByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loyalty account"})
		return
	}

	// Get transaction history from our database - convert uint to string
	accountIDStr := fmt.Sprintf("%d", account.ID)
	transactions, err := h.dbService.GetTransactionHistory(c.Request.Context(), accountIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction history"})
		return
	}

	// Process transactions for frontend display
	type transactionResponse struct {
		ID            uint      `json:"id"`
		AccountID     uint      `json:"account_id"`
		Type          string    `json:"type"`
		Points        int       `json:"points"`
		TransactionAt time.Time `json:"transaction_at"`
		CreatedAt     time.Time `json:"created_at"`
	}

	var response []transactionResponse
	for _, tx := range transactions {
		// Use transaction timestamp as transaction time
		txTime := tx.Timestamp

		response = append(response, transactionResponse{
			ID:            tx.ID,
			AccountID:     tx.AccountID,
			Type:          tx.Type,
			Points:        tx.Points,
			TransactionAt: txTime,
			CreatedAt:     tx.CreatedAt,
		})
	}

	// Try to get Square transactions if we have a Square loyalty ID
	var squareTransactions []models.Transaction
	squareLoyaltyID, err := h.dbService.GetSquareLoyaltyID(c.Request.Context(), accountIDStr)
	if err == nil && squareLoyaltyID != "" {
		// Get transaction history from Square
		squareTx, err := h.squareService.GetTransactionHistory(c.Request.Context(), squareLoyaltyID)
		if err == nil {
			// Convert Square transactions to our model
			for _, tx := range squareTx {
				squareTransactions = append(squareTransactions, models.Transaction{
					ID:        0, // We don't store these in our DB
					AccountID: account.ID,
					Type:      tx.Type,
					Points:    tx.Points,
					Timestamp: tx.Timestamp,
					CreatedAt: tx.Timestamp,
					UpdatedAt: tx.Timestamp,
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"transactions":        response,
				"square_transactions": squareTransactions,
				"square_account_id":   squareLoyaltyID,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"transactions": response})
}
