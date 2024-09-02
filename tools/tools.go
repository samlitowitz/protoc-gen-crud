//go:build tools

package tools

import (
	_ "github.com/jackc/pgx-zap"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)
