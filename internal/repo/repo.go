package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"simple-service/internal/config"
)

// Слой репозитория, здесь должны быть все методы, связанные с базой данных

// SQL-запрос на вставку задачи
const (
	insertTaskQuery  = `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id;`
	getTaskByIdQuery = `SELECT id, title, description, created_at, updated_at, status FROM tasks WHERE id = $1;`
)

type repository struct {
	pool *pgxpool.Pool
}

// Repository - интерфейс с методом создания задачи
type Repository interface {
	CreateTask(ctx context.Context, task Task) (int, error)   // Создание задачи
	GetTaskById(ctx context.Context, id int) (GetTask, error) // Получение задачи
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

// CreateTask - вставка новой задачи в таблицу tasks
func (r *repository) CreateTask(ctx context.Context, task Task) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert task")
	}
	return id, nil
}

func (r *repository) GetTaskById(ctx context.Context, id int) (GetTask, error) {
	var task GetTask
	err := r.pool.QueryRow(ctx, getTaskByIdQuery, id).Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Status)
	if err != nil {
		return GetTask{}, errors.Wrap(err, "failed to query task")
	}

	return task, nil
}
