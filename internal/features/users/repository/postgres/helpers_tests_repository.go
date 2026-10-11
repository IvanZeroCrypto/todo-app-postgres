package users_postgres_repository

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
	"github.com/stretchr/testify/require"
)

var (
	emptyFullName   = ""
	fullName        = "Ivan Ivanov"
	shortFullName   = "Iv"
	longFullName    = strings.Repeat("a", 101)
	fullNameNumbers = "463473563433"
	fullNameSymbols = "%^$&#($&%^##*$)"
	fullNameArray   = "[Ivan Ivanov]"
	fullNameObject  = "{name:Ivan}"

	phoneNumber         = "+77777777777"
	shortPhoneNumber    = "+78987"
	longPhoneNumber     = "+7896786786786786786"
	symbolsPhoneNumber  = "+%^$*!@#"
	lettersPhoneNumber  = "+werfdsazxcf"
	arrayPhoneNumber    = "[+789789789789]"
	objectPhoneNumber   = "{phoneNumber:+789789789789}"
	noPrefixPhoneNumber = "77777777777"
)

// TestPool реализует ваш интерфейс core_postgres_pool.Pool для тестов
type TestPool struct {
	pool      *pgxpool.Pool
	opTimeout time.Duration
}

func (p *TestPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *TestPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, arguments...)
}

func (p *TestPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *TestPool) Close() {
	p.pool.Close()
}

func (p *TestPool) OpTimeout() time.Duration {
	return p.opTimeout
}

// NewTestPool создаёт тестовый Pool с изолированной БД
func NewTestPool(t *testing.T) *TestPool {
	t.Helper()

	// Конфигурация pgtestdb
	conf := pgtestdb.Config{
		DriverName: "pgx",
		Host:       "localhost",
		Port:       "5433",
		User:       "testuser",
		Password:   "testpass",
		Database:   "postgres",
		Options:    "sslmode=disable",
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current file path")
	}

	// Поднимаемся на 6 уровней от postgres/ до корня проекта
	// postgres -> repository -> users -> features -> internal -> todoApp/
	projectRoot := filepath.Dir(filename)   // .../postgres
	projectRoot = filepath.Dir(projectRoot) // .../repository
	projectRoot = filepath.Dir(projectRoot) // .../users
	projectRoot = filepath.Dir(projectRoot) // .../features
	projectRoot = filepath.Dir(projectRoot) // .../internal
	projectRoot = filepath.Dir(projectRoot) // .../todoApp (корень)

	migrationsPath := filepath.Join(projectRoot, "migrations")
	t.Logf("Project root: %s", projectRoot)
	t.Logf("Migrations path: %s", migrationsPath)
	migrator := golangmigrator.New(migrationsPath)

	// pgtestdb создаёт новую БД и возвращает *sql.DB
	sqlDB := pgtestdb.New(t, conf, migrator)

	// Получаем имя созданной БД (оно уникально для каждого теста)
	var dbName string
	err := sqlDB.QueryRowContext(context.Background(), "SELECT current_database()").Scan(&dbName)
	require.NoError(t, err)

	// Строим connection string для pgxpool
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, dbName,
	)

	// Создаём pgxpool
	pgxConfig, err := pgxpool.ParseConfig(connString)
	require.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	require.NoError(t, err)

	// Проверяем подключение
	err = pool.Ping(context.Background())
	require.NoError(t, err)

	// Закрываем ресурсы после теста
	t.Cleanup(func() {
		pool.Close()
		sqlDB.Close()
	})

	// Таймаут для тестов (можно маленький для проверки)
	return &TestPool{
		pool:      pool,
		opTimeout: 5 * time.Second,
	}
}

func NewTestPoolWithTimeout(t *testing.T, timeout time.Duration) *TestPool {
	pool := NewTestPool(t)
	pool.opTimeout = timeout
	return pool
}
