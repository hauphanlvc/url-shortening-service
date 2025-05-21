package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	db "url-shortening-service/db/sqlc"
)

type Service interface {
	GenerateService
}
