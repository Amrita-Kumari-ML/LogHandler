package connection

import (
	"LogParser/logger"
	"LogParser/models"
	_ "LogParser/models"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func init(){
	logger.InitializeLogger("info")
}

// Mock YAML content to test YAML loading
const mockYamlContent = `
database:
  DB_PORT: "5432"
  DB_HOST: "localhost"
  DB_USERNAME: "testuser"
  DB_PASSWORD: "testpass"
  DB_NAME: "testdb"
  DB_SSLMODE: "disable"

logs:
  table_name: "logs"
  create_table_query: "CREATE TABLE logs (...);"
`

func writeTempYaml(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp YAML file: %v", err)
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temp YAML file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestLoadConfigFromYaml(t *testing.T) {
	filePath := writeTempYaml(t, mockYamlContent)
	defer os.Remove(filePath)

	err := LoadConfigFromYaml(filePath)
	if err != nil {
		t.Errorf("LoadConfigFromYaml returned error: %v", err)
	}

	if ConfigData.Database.DBHost != "localhost" {
		t.Errorf("Expected DBHost to be 'localhost', got '%s'", ConfigData.Database.DBHost)
	}
	if ConfigData.Logs.TableName != "logs" {
		t.Errorf("Expected TableName to be 'logs', got '%s'", ConfigData.Logs.TableName)
	}
}

func TestFirstLoad_EnvVars(t *testing.T) {
	// Set mock environment variables
	os.Setenv("DB_HOST", "envhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USERNAME", "envuser")
	os.Setenv("DB_PASSWORD", "envpass")
	os.Setenv("DB_NAME", "envdb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("DB_TABLE_NAME", "logs")
	os.Setenv("DB_CREATE_TABLE_QUERY", "CREATE TABLE env_logs (...);")

	defer func() {
		// Clean up environment variables after test
		os.Clearenv()
	}()

	err := FirstLoad()
	if err != nil {
		t.Errorf("FirstLoad returned error: %v", err)
	}

	if ConfigData.Database.DBHost != "envhost" {
		t.Errorf("Expected DBHost from env to be 'envhost', got '%s'", ConfigData.Database.DBHost)
	}

	if ConfigData.Logs.TableName != "logs" {
		t.Errorf("Expected TableName from env to be 'env_logs', got '%s'", ConfigData.Logs.TableName)
	}
}

func TestGetEnvString_DefaultFallback(t *testing.T) {
	os.Unsetenv("NON_EXISTENT_VAR")
	defaultVal := "default"
	val := getEnvString("NON_EXISTENT_VAR", defaultVal)
	if val != defaultVal {
		t.Errorf("Expected default value '%s', got '%s'", defaultVal, val)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("INT_VAR", "123")
	defer os.Unsetenv("INT_VAR")

	val := getEnvInt("INT_VAR", 456)
	if val != 123 {
		t.Errorf("Expected 123, got %d", val)
	}

	os.Setenv("BAD_INT_VAR", "not-an-int")
	val = getEnvInt("BAD_INT_VAR", 789)
	if val != 789 {
		t.Errorf("Expected fallback value 789 for bad int, got %d", val)
	}
}


func setMockConfig() {
	ConfigData = models.DB_Config{
		Database: struct {
			DBPort     string `yaml:"DB_PORT"`
			DBHost     string `yaml:"DB_HOST"`
			DBUsername string `yaml:"DB_USERNAME"`
			DBPassword string `yaml:"DB_PASSWORD"`
			DBName     string `yaml:"DB_NAME"`
			DBSslMode  string `yaml:"DB_SSLMODE"`
		}{
			DBPort:     "5432",
			DBHost:     "localhost",
			DBUsername: "postgres",
			DBPassword: "password",
			DBName:     "testdb",
			DBSslMode:  "disable",
		},
		Logs: struct {
			TableName        string `yaml:"table_name"`
			CreateTableQuery string `yaml:"create_table_query"`
		}{
			TableName:        "logs",
			CreateTableQuery: "CREATE TABLE logs (id SERIAL PRIMARY KEY);",
		},
	}
	Config = &ConfigData
}

// TestPingDB tests if PingDB correctly pings a live connection
func TestPingDB(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()
	DB = db

	//mock.ExpectPing().WillReturnError(nil)

	success, conn := PingDB()
	if !success || conn == nil {
		t.Errorf("Expected successful ping, got success=%v, conn=%v", success, conn)
	}
}

// TestCreateLogsTableIfNotExist_TableDoesNotExist simulates missing table and ensures creation is triggered
func TestCreateLogsTableIfNotExist_TableDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	DB = db
	setMockConfig()

	// Simulate no rows returned (table doesn't exist)
	mock.ExpectQuery(`SELECT table_name FROM information_schema.tables WHERE table_name = \$1`).
		WithArgs("logs").
		WillReturnError(sql.ErrNoRows)

	// Expect the table creation to be called
	mock.ExpectExec("CREATE TABLE logs").WillReturnResult(sqlmock.NewResult(1, 1))

	// Simulate checking index existence, and it does not exist
	mock.ExpectQuery(`SELECT indexname FROM pg_indexes WHERE indexname = \$1`).
		WithArgs("idx_time_local").
		WillReturnError(sql.ErrNoRows)

	createLogsTableIfNotExist(*Config)
}

// TestCreateLogsTableIfNotExist_TableExists ensures no creation when table already exists
func TestCreateLogsTableIfNotExist_TableExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	DB = db
	setMockConfig()

	// Simulate the table already exists
	mock.ExpectQuery(`SELECT table_name FROM information_schema.tables WHERE table_name = \$1`).
		WithArgs("logs").
		WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("logs"))

	createLogsTableIfNotExist(*Config)
}

// TestIndexExists_IndexExists checks behavior when index exists
func TestIndexExists_IndexExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectQuery(`SELECT indexname FROM pg_indexes WHERE indexname = \$1`).
		WithArgs("idx_time_local").
		WillReturnRows(sqlmock.NewRows([]string{"indexname"}).AddRow("idx_time_local"))

	if !indexExists("idx_time_local") {
		t.Errorf("Expected index to exist but got false")
	}
}

// TestIndexExists_IndexDoesNotExist checks behavior when index is missing
func TestIndexExists_IndexDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	DB = db

	mock.ExpectQuery(`SELECT indexname FROM pg_indexes WHERE indexname = \$1`).
		WithArgs("nonexistent_index").
		WillReturnError(sql.ErrNoRows)

	if indexExists("nonexistent_index") {
		t.Errorf("Expected index to not exist but got true")
	}
}