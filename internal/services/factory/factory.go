package factory

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Woland-prj/dilemator/config"
	"github.com/Woland-prj/dilemator/internal/repo/dilemma_repo"
	"github.com/Woland-prj/dilemator/internal/repo/sessions_repo"
	"github.com/Woland-prj/dilemator/internal/repo/users_repo"
	"github.com/Woland-prj/dilemator/internal/services/dilemma_service"
	"github.com/Woland-prj/dilemator/internal/services/sessions_service"
	"github.com/Woland-prj/dilemator/internal/services/users_service"
	"github.com/Woland-prj/dilemator/pkg/hashing"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/Woland-prj/dilemator/pkg/postgres"
)

type ServiceFactory struct {
	cfg    *config.Config
	logger logger.Interface
	pg     *postgres.Postgres
	hash   *hashing.HashProvider
	// s3Repo *s3.FileS3RepositoryAdapter

	// Кэширование сервисов (опционально)
	userService    users_service.UserService
	sessionService sessions_service.SessionService
	dilemmaService dilemma_service.DilemmaService
	// fileService     files_service.FileService
}

var ErrUnsupportedService = errors.New("unsupported service type")

func NewServiceFactory(cfg *config.Config, l logger.Interface) (*ServiceFactory, error) {
	pg, err := postgres.New(getPgConf(cfg), postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres: %w", err)
	}

	hash := hashing.NewHashProvider(cfg.App.PassCost)

	// s3Repo, err := getS3Repo(cfg)
	// if err != nil {
	// 	pg.Close()

	// 	return nil, fmt.Errorf("failed to init S3 repo: %w", err)
	// }

	return &ServiceFactory{
		cfg:    cfg,
		logger: l,
		pg:     pg,
		hash:   hash,
	}, nil
}

// Close free resources.
func (f *ServiceFactory) Close() {
	if f.pg != nil {
		f.pg.Close()
	}
}

func (f *ServiceFactory) instantiateUsersService() users_service.UserService {
	userRepo := users_repo.NewUserRepositoryAdapter(f.pg)

	if f.userService == nil {
		f.userService = users_service.NewUserService(f.logger, userRepo, f.hash)
	}

	return f.userService
}

func (f *ServiceFactory) instantiateSessionsService() (sessions_service.SessionService, error) {
	if f.sessionService == nil {
		sessionRepo := sessions_repo.NewSessionRepositoryAdapter(f.pg)

		us, err := InstantiateService[users_service.UserService](f)
		if err != nil {
			return nil, err
		}

		f.sessionService = sessions_service.NewSessionService(
			f.logger,
			sessionRepo,
			f.hash,
			us,
			f.cfg.HTTP.CookieMaxAgeDays,
			f.cfg.App.TgBotToken,
		)
	}

	return f.sessionService, nil
}

func (f *ServiceFactory) instantiateDilemmaService() dilemma_service.DilemmaService {
	if f.dilemmaService == nil {
		dilemmaRepo := dilemma_repo.NewDilemmaRepositoryAdapter(f.pg)
		f.dilemmaService = dilemma_service.NewDilemmaService(f.logger, dilemmaRepo)
	}

	return f.dilemmaService
}

// func (f *ServiceFactory) instantiateFileService() filecontracts.FileService {
// 	if f.fileService == nil {
// 		f.fileService = files.NewFileService(f.s3Repo, f.logger)
// 	}

// 	return f.fileService
// }

type serviceConstructor func(*ServiceFactory) (any, error)

func getServiceRegistry() map[reflect.Type]serviceConstructor {
	return map[reflect.Type]serviceConstructor{
		reflect.TypeOf((*users_service.UserService)(nil)).Elem(): func(f *ServiceFactory) (any, error) {
			return f.instantiateUsersService(), nil
		},
		reflect.TypeOf((*sessions_service.SessionService)(nil)).Elem(): func(f *ServiceFactory) (any, error) {
			return f.instantiateSessionsService()
		},
		reflect.TypeOf((*dilemma_service.DilemmaService)(nil)).Elem(): func(f *ServiceFactory) (any, error) {
			return f.instantiateDilemmaService(), nil
		},
		// reflect.TypeOf((*files_service.FileService)(nil)).Elem(): func(f *ServiceFactory) (any, error) {
		// 	return f.instantiateFileService(), nil
		// },
	}
}

// InstantiateService returns service by type T.
func InstantiateService[T any](f *ServiceFactory) (T, error) {
	var zero T

	t := reflect.TypeOf((*T)(nil)).Elem()

	constructor, ok := getServiceRegistry()[t]
	if !ok {
		return zero, fmt.Errorf("%w: %T", ErrUnsupportedService, zero)
	}

	svc, err := constructor(f)
	if err != nil {
		return zero, err
	}

	val, ok := svc.(T)
	if !ok {
		return zero, fmt.Errorf("%w: expected %T, got %T", ErrUnsupportedService, zero, svc)
	}

	return val, nil
}

func getPgConf(cfg *config.Config) *postgres.Config {
	return &postgres.Config{
		User:     cfg.PG.User,
		Password: cfg.PG.Password,
		DbName:   cfg.PG.DbName,
		SslMode:  cfg.PG.SslMode,
		Host:     cfg.PG.Host,
		Port:     cfg.PG.Port,
		LogEnv:   cfg.Log.Env,
	}
}

// func getS3Repo(cfg *config.Config) (*s3.FileS3RepositoryAdapter, error) {
// 	s3FileAdapter, err := s3.NewFileS3Repository(
// 		cfg.S3.AccessKey,
// 		cfg.S3.SecretKey,
// 		cfg.S3.BucketName,
// 		cfg.S3.Region,
// 		cfg.S3.Endpoint,
// 		cfg.S3.BucketDomain,
// 		cfg.S3.PresignLifetimeHours,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s3FileAdapter, nil
// }
