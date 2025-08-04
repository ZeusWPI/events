package storage

import (
	"fmt"
	"time"

	"github.com/ZeusWPI/events/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/minio"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

var S fiber.Storage

func Init(pool *pgxpool.Pool) error {
	provider := config.GetDefaultString("storage.provider", "postgres")

	switch provider {

	case "minio":
		S = minio.New(minio.Config{
			Bucket:   config.GetDefaultString("minio.bucket", "events"),
			Endpoint: config.GetString("minio.endpoint"),
			Secure:   config.GetDefaultBool("minio.secure", false),
			Credentials: minio.Credentials{
				AccessKeyID:     config.GetString("minio.username"),
				SecretAccessKey: config.GetString("minio.password"),
			},
		})

	case "postgres":
		S = postgres.New(postgres.Config{
			DB:         pool,
			Table:      "events_files",
			Reset:      false,
			GCInterval: 10 * time.Second,
		})

	default:
		return fmt.Errorf("unsupported storage provider %s", provider)
	}

	return nil
}
