package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/smthjapanese/avito_pvz/internal/app"
	"github.com/smthjapanese/avito_pvz/internal/config"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

type dummyLoginRequest struct {
	Role models.UserRole `json:"role"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type createPVZRequest struct {
	City models.City `json:"city"`
}

type pvzResponse struct {
	ID               uuid.UUID   `json:"id"`
	RegistrationDate time.Time   `json:"registrationDate"`
	City             models.City `json:"city"`
}

type createReceptionRequest struct {
	PVZID uuid.UUID `json:"pvzId"`
}

type receptionResponse struct {
	ID       uuid.UUID              `json:"id"`
	DateTime time.Time              `json:"dateTime"`
	PVZID    uuid.UUID              `json:"pvzId"`
	Status   models.ReceptionStatus `json:"status"`
}

type createProductRequest struct {
	Type  models.ProductType `json:"type"`
	PVZID uuid.UUID          `json:"pvzId"`
}

type productResponse struct {
	ID          uuid.UUID          `json:"id"`
	DateTime    time.Time          `json:"dateTime"`
	Type        models.ProductType `json:"type"`
	ReceptionID uuid.UUID          `json:"receptionId"`
}

func TestPVZWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Log("Starting integration test")

	cfg, err := config.Load("../../configs/test_config.yaml")
	require.NoError(t, err, "Failed to load test configuration")

	setupTestDatabase(t, cfg)

	application, err := app.NewApp(cfg)
	require.NoError(t, err, "Failed to create application")

	err = application.Run()
	require.NoError(t, err, "Failed to run application")

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := application.Shutdown(ctx)
		if err != nil {
			t.Logf("Error shutting down application: %v", err)
		}
	}()

	baseURL := fmt.Sprintf("http://localhost:%s", cfg.Server.HTTPPort)

	time.Sleep(2 * time.Second)
	t.Log("Application should be up and running now")

	t.Log("Step 1: Getting moderator token")
	moderatorToken := getModeratorToken(t, baseURL)
	require.NotEmpty(t, moderatorToken, "Moderator token should not be empty")

	t.Log("Step 2: Creating PVZ with moderator role")
	pvz := createPVZ(t, baseURL, moderatorToken, models.CityMoscow)
	require.NotNil(t, pvz, "PVZ should not be nil")
	require.NotEqual(t, uuid.Nil, pvz.ID, "PVZ ID should not be nil")
	t.Logf("Created PVZ with ID: %s", pvz.ID)

	t.Log("Step 3: Getting employee token")
	employeeToken := getEmployeeToken(t, baseURL)
	require.NotEmpty(t, employeeToken, "Employee token should not be empty")

	t.Log("Step 4: Creating reception with employee role")
	reception := createReception(t, baseURL, employeeToken, pvz.ID)
	require.NotNil(t, reception, "Reception should not be nil")
	require.NotEqual(t, uuid.Nil, reception.ID, "Reception ID should not be nil")
	t.Logf("Created reception with ID: %s", reception.ID)

	t.Log("Step 5: Adding 50 products")
	productTypes := []models.ProductType{
		models.ProductTypeElectronics,
		models.ProductTypeClothes,
		models.ProductTypeShoes,
	}

	for i := 0; i < 50; i++ {
		productType := productTypes[i%len(productTypes)]
		product := createProduct(t, baseURL, employeeToken, pvz.ID, productType)
		require.NotNil(t, product, "Product should not be nil")
		require.NotEqual(t, uuid.Nil, product.ID, "Product ID should not be nil")
		t.Logf("Created product %d/%d with ID: %s", i+1, 50, product.ID)
	}

	t.Log("Step 6: Closing reception")
	closedReception := closeReception(t, baseURL, employeeToken, pvz.ID)
	require.NotNil(t, closedReception, "Closed reception should not be nil")
	require.Equal(t, models.ReceptionStatusClose, closedReception.Status, "Reception status should be 'close'")
	t.Logf("Closed reception with ID: %s", closedReception.ID)

	t.Cleanup(func() {
		db, err := sql.Open("postgres", fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
			cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
		))
		if err != nil {
			t.Logf("Failed to connect to database for cleanup: %v", err)
			return
		}
		defer db.Close()

		_, err = db.Exec("DELETE FROM products WHERE reception_id IN (SELECT id FROM receptions WHERE pvz_id = $1)", pvz.ID)
		if err != nil {
			t.Logf("Failed to clean up products: %v", err)
		}

		_, err = db.Exec("DELETE FROM receptions WHERE pvz_id = $1", pvz.ID)
		if err != nil {
			t.Logf("Failed to clean up receptions: %v", err)
		}

		_, err = db.Exec("DELETE FROM pvzs WHERE id = $1", pvz.ID)
		if err != nil {
			t.Logf("Failed to clean up pvz: %v", err)
		}
	})

	t.Log("Test completed successfully!")
}

func setupTestDatabase(t *testing.T, cfg *config.Config) {
	connStrPostgres := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.SSLMode,
	)

	dbPostgres, err := sql.Open("postgres", connStrPostgres)
	require.NoError(t, err, "Failed to connect to postgres database")
	defer dbPostgres.Close()

	// Проверка соединения
	err = dbPostgres.Ping()
	require.NoError(t, err, "Failed to ping postgres database")

	t.Log("Successfully connected to the postgres database")

	var exists bool
	err = dbPostgres.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.Database.DBName).Scan(&exists)
	require.NoError(t, err, "Failed to check if test database exists")

	if !exists {
		t.Logf("Creating test database: %s", cfg.Database.DBName)
		_, err = dbPostgres.Exec("CREATE DATABASE " + cfg.Database.DBName)
		require.NoError(t, err, "Failed to create test database")
		t.Logf("Created test database: %s", cfg.Database.DBName)
	} else {
		t.Logf("Test database already exists: %s", cfg.Database.DBName)
	}

	dbPostgres.Close()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "Failed to connect to test database")
	defer db.Close()

	// Проверка соединения с тестовой базой
	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	t.Log("Successfully connected to the test database")

	_, err = db.Exec(`
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
                CREATE TYPE user_role AS ENUM ('employee', 'moderator');
            END IF;
            
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'city_type') THEN
                CREATE TYPE city_type AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
            END IF;
            
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type') THEN
                CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');
            END IF;
            
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reception_status') THEN
                CREATE TYPE reception_status AS ENUM ('in_progress', 'close');
            END IF;
        END
        $$;

        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            email VARCHAR(255) NOT NULL UNIQUE,
            password_hash VARCHAR(255) NOT NULL,
            role user_role NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS pvzs (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            city city_type NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS receptions (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            date_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            pvz_id UUID NOT NULL REFERENCES pvzs(id),
            status reception_status NOT NULL DEFAULT 'in_progress',
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS products (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            date_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            type product_type NOT NULL,
            reception_id UUID NOT NULL REFERENCES receptions(id),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE INDEX IF NOT EXISTS idx_receptions_pvz_id ON receptions(pvz_id);
        CREATE INDEX IF NOT EXISTS idx_receptions_status ON receptions(status);
        CREATE INDEX IF NOT EXISTS idx_products_reception_id ON products(reception_id);
        CREATE INDEX IF NOT EXISTS idx_products_date_time ON products(date_time);
    `)
	if err != nil {
		t.Logf("Warning during schema setup: %v", err)
		// Не останавливаем тест, так как ошибка может быть из-за того, что объекты уже существуют
	}

	t.Log("Database schema setup completed")
}

// Вспомогательные функции для работы с HTTP-запросами

func getModeratorToken(t *testing.T, baseURL string) string {
	return getToken(t, baseURL, models.ModeratorRole)
}

func getEmployeeToken(t *testing.T, baseURL string) string {
	return getToken(t, baseURL, models.EmployeeRole)
}

func getToken(t *testing.T, baseURL string, role models.UserRole) string {
	req := dummyLoginRequest{
		Role: role,
	}
	reqBody, err := json.Marshal(req)
	require.NoError(t, err, "Failed to marshal request body")

	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(reqBody))
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	if resp.StatusCode != http.StatusOK {
		t.Logf("Error response: %s", string(body))
		require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, got %d", resp.StatusCode)
		return ""
	}

	bodyStr := string(body)
	t.Logf("Token response body: %s", bodyStr)

	if len(bodyStr) >= 2 && bodyStr[0] == '"' && bodyStr[len(bodyStr)-1] == '"' {
		return bodyStr[1 : len(bodyStr)-1]
	}

	var response tokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Logf("Failed to unmarshal as JSON, treating as raw token: %v", err)
		return bodyStr
	}

	return response.Token
}

func createPVZ(t *testing.T, baseURL string, token string, city models.City) *pvzResponse {
	req := createPVZRequest{
		City: city,
	}
	reqBody, err := json.Marshal(req)
	require.NoError(t, err, "Failed to marshal request body")

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, baseURL+"/pvz", bytes.NewBuffer(reqBody))
	require.NoError(t, err, "Failed to create request")

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	if resp.StatusCode != http.StatusCreated {
		t.Logf("Error response: %s", string(body))
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, got %d", resp.StatusCode)

	var response pvzResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to unmarshal response body")

	return &response
}

func createReception(t *testing.T, baseURL string, token string, pvzID uuid.UUID) *receptionResponse {
	req := createReceptionRequest{
		PVZID: pvzID,
	}
	reqBody, err := json.Marshal(req)
	require.NoError(t, err, "Failed to marshal request body")

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, baseURL+"/receptions", bytes.NewBuffer(reqBody))
	require.NoError(t, err, "Failed to create request")

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	if resp.StatusCode != http.StatusCreated {
		t.Logf("Error response: %s", string(body))
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, got %d", resp.StatusCode)

	var response receptionResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to unmarshal response body")

	return &response
}

func createProduct(t *testing.T, baseURL string, token string, pvzID uuid.UUID, productType models.ProductType) *productResponse {
	req := createProductRequest{
		Type:  productType,
		PVZID: pvzID,
	}
	reqBody, err := json.Marshal(req)
	require.NoError(t, err, "Failed to marshal request body")

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, baseURL+"/products", bytes.NewBuffer(reqBody))
	require.NoError(t, err, "Failed to create request")

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	if resp.StatusCode != http.StatusCreated {
		t.Logf("Error response: %s", string(body))
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, got %d", resp.StatusCode)

	var response productResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to unmarshal response body")

	return &response
}

func closeReception(t *testing.T, baseURL string, token string, pvzID uuid.UUID) *receptionResponse {
	client := &http.Client{}
	url := fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err, "Failed to create request")

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	if resp.StatusCode != http.StatusOK {
		t.Logf("Error response: %s", string(body))
	}
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, got %d", resp.StatusCode)

	var response receptionResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to unmarshal response body")

	return &response
}
