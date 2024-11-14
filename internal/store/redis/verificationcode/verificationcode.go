package verificationcode

import (
	"context"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/redis"
	rdb "github.com/redis/go-redis/v9"
)

type storage struct {
	*redis.Redis
}

func New(rdb *redis.Redis) core.VerificationStore {
	return &storage{rdb}
}

func (s *storage) DeleteVerificationCode(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := s.Client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (s *storage) GetVerificationCode(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userID, err := s.Client.Get(ctx, key).Result()
	if err == rdb.Nil {
		return "", core.ErrVerificationCodeNotValid
	} else if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *storage) SetVerificationCode(ctx context.Context, key string, val string, expiresIn time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := s.Client.Set(ctx, key, val, expiresIn).Err(); err != nil {
		return err
	}

	return nil
}
