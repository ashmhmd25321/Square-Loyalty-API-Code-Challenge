package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"loyalty-app/internal/models"
)

// MemoryService provides in-memory storage for loyalty program data
type MemoryService struct {
	mu               sync.RWMutex
	users            map[string]*models.User
	loyaltyAccounts  map[uint]*models.LoyaltyAccount
	transactions     map[uint][]models.Transaction
	squareLoyaltyIDs map[uint]string // Stores Square loyalty IDs
	nextUserID       uint
	nextAccountID    uint
	nextTxnID        uint
}

// NewMemoryService creates a new MemoryService instance
func NewMemoryService() *MemoryService {
	return &MemoryService{
		users:            make(map[string]*models.User),
		loyaltyAccounts:  make(map[uint]*models.LoyaltyAccount),
		transactions:     make(map[uint][]models.Transaction),
		squareLoyaltyIDs: make(map[uint]string),
		nextUserID:       1,
		nextAccountID:    1,
		nextTxnID:        1,
	}
}

// UserExists checks if a user with the given email exists
func (s *MemoryService) UserExists(ctx context.Context, email string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.users[email]
	return exists, nil
}

// CreateUser creates a new user with the given email and password hash
func (s *MemoryService) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[email]; exists {
		return nil, errors.New("user already exists")
	}

	now := time.Now()
	user := &models.User{
		ID:           s.nextUserID,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.users[email] = user
	s.nextUserID++

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *MemoryService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *MemoryService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	for _, user := range s.users {
		if user.ID == uint(userIDUint) {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetLoyaltyAccountByUserID retrieves a loyalty account by user ID
func (s *MemoryService) GetLoyaltyAccountByUserID(ctx context.Context, userID string) (*models.LoyaltyAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	for _, account := range s.loyaltyAccounts {
		if account.UserID == uint(userIDUint) {
			return account, nil
		}
	}

	return nil, errors.New("loyalty account not found")
}

// CreateLoyaltyAccount creates a new loyalty account for a user
func (s *MemoryService) CreateLoyaltyAccount(ctx context.Context, userID string) (*models.LoyaltyAccount, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	for _, account := range s.loyaltyAccounts {
		if account.UserID == uint(userIDUint) {
			return nil, errors.New("loyalty account already exists for this user")
		}
	}

	var user *models.User
	for _, u := range s.users {
		if u.ID == uint(userIDUint) {
			user = u
			break
		}
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	now := time.Now()
	account := &models.LoyaltyAccount{
		ID:        s.nextAccountID,
		UserID:    uint(userIDUint),
		User:      *user,
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.loyaltyAccounts[account.ID] = account
	s.nextAccountID++

	return account, nil
}

// AddTransaction adds a transaction to the database and updates the account balance
func (s *MemoryService) AddTransaction(ctx context.Context, accountID string, transactionType string, points int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid account ID: %v", err)
	}

	account, exists := s.loyaltyAccounts[uint(accountIDUint)]
	if !exists {
		return errors.New("account not found")
	}

	// Update balance
	if transactionType == "EARN" {
		account.Balance += points
	} else if transactionType == "REDEEM" {
		if account.Balance < points {
			return errors.New("insufficient points")
		}
		account.Balance -= points
	} else {
		return errors.New("invalid transaction type")
	}

	now := time.Now()

	// Create transaction
	transaction := models.Transaction{
		ID:        s.nextTxnID,
		AccountID: uint(accountIDUint),
		Type:      transactionType,
		Points:    points,
		Timestamp: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Store transaction
	if s.transactions[uint(accountIDUint)] == nil {
		s.transactions[uint(accountIDUint)] = make([]models.Transaction, 0)
	}
	s.transactions[uint(accountIDUint)] = append(s.transactions[uint(accountIDUint)], transaction)
	s.nextTxnID++

	return nil
}

// GetBalance retrieves the current points balance for a loyalty account
func (s *MemoryService) GetBalance(ctx context.Context, accountID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid account ID: %v", err)
	}

	account, exists := s.loyaltyAccounts[uint(accountIDUint)]
	if !exists {
		return 0, errors.New("account not found")
	}

	return account.Balance, nil
}

// GetTransactionHistory retrieves the transaction history for a loyalty account
func (s *MemoryService) GetTransactionHistory(ctx context.Context, accountID string) ([]models.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %v", err)
	}

	transactions, exists := s.transactions[uint(accountIDUint)]
	if !exists {
		return []models.Transaction{}, nil
	}

	return transactions, nil
}

// UpdateSquareLoyaltyID updates the Square loyalty ID for a loyalty account
func (s *MemoryService) UpdateSquareLoyaltyID(ctx context.Context, accountID string, squareLoyaltyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid account ID: %v", err)
	}

	_, exists := s.loyaltyAccounts[uint(accountIDUint)]
	if !exists {
		return errors.New("account not found")
	}

	// Store Square loyalty ID in the map
	s.squareLoyaltyIDs[uint(accountIDUint)] = squareLoyaltyID

	return nil
}

// GetSquareLoyaltyID gets the Square loyalty ID for a loyalty account
func (s *MemoryService) GetSquareLoyaltyID(ctx context.Context, accountID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid account ID: %v", err)
	}

	_, exists := s.loyaltyAccounts[uint(accountIDUint)]
	if !exists {
		return "", errors.New("account not found")
	}

	loyaltyID, exists := s.squareLoyaltyIDs[uint(accountIDUint)]
	if !exists {
		return "", nil
	}

	return loyaltyID, nil
}
