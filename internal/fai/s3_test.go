package fai

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

func Test