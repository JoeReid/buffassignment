package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/jmoiron/sqlx"
)

var _ model.Store = &Store{}

type Store struct {
	host     string
	port     int
	database string
	user     string
	password string

	// Database connection options
	maxOpenCons    int
	maxIdleCons    int
	maxConLifetime time.Duration
	connectTimeout time.Duration

	db *sqlx.DB
}

var (
	videoStreamTable  = "video_streams"
	videoStreamFields = []string{"id", "title", "created", "updated"}

	questionTable  = "questions"
	questionFields = []string{"id", "stream", "text"}

	answerTable  = "answers"
	answerFields = []string{"id", "question", "text", "correct"}

	buffFields = []string{
		"questions.id", "questions.stream", "questions.text",
		"answers.id", "answers.question", "answers.text", "answers.correct",
	}
)

type StoreOption func(*Store) error

func NewStore(options ...StoreOption) (*Store, error) {
	const (
		defaultPort           = 5432
		defaultMaxOpenCons    = 10
		defaultMaxIdleCons    = 10
		defaultMaxConLifetime = time.Hour
		defaultConnectTimeout = 10 * time.Second
	)

	s := &Store{
		port:           defaultPort,
		maxOpenCons:    defaultMaxOpenCons,
		maxIdleCons:    defaultMaxIdleCons,
		maxConLifetime: defaultMaxConLifetime,
		connectTimeout: defaultConnectTimeout,
	}

	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	switch {
	case s.user == "":
		return nil, errors.New("required config not set: user")
	case s.database == "":
		return nil, errors.New("required config not set: database name")
	case s.password == "":
		return nil, errors.New("required config not set: password")
	case s.host == "":
		return nil, errors.New("required config not set: hostname")
	default:

		connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%d", s.user, s.database, s.password, s.host, s.port)
		db, err := sqlx.Open("postgres", connectionString)
		if err != nil {
			return nil, err
		}

		then := time.Now()
		for {
			err := db.Ping()
			if err == nil {
				break
			}

			if then.Add(s.connectTimeout).Before(time.Now()) {
				return nil, err
			}
		}

		db.DB.SetConnMaxLifetime(s.maxConLifetime)
		db.DB.SetMaxIdleConns(s.maxIdleCons)
		db.DB.SetMaxOpenConns(s.maxOpenCons)

		s.db = db
		return s, nil
	}
}

func SetDBUser(user string) StoreOption {
	return func(p *Store) error {
		p.user = user
		return nil
	}
}

func SetDBPassword(pass string) StoreOption {
	return func(p *Store) error {
		p.password = pass
		return nil
	}
}

func SetDBHostname(host string) StoreOption {
	return func(p *Store) error {
		p.host = host
		return nil
	}
}

func SetDBPort(port int) StoreOption {
	return func(p *Store) error {
		p.port = port
		return nil
	}
}

func SetDBName(name string) StoreOption {
	return func(p *Store) error {
		p.database = name
		return nil
	}
}

func SetConnectTimeout(t time.Duration) StoreOption {
	return func(p *Store) error {
		p.connectTimeout = t
		return nil
	}
}

// WithMaxOpenCons is a function option for NewStore that sets
// limits on the maximum open cons the store can use when connecting to the
// database.
func WithMaxOpenCons(n int) StoreOption {
	return func(p *Store) error {
		if n <= 0 {
			return fmt.Errorf("cannot set maxOpenCons to %d", n)
		}

		p.maxOpenCons = n
		return nil
	}
}

// WithMaxIdleCons is a function option for NewStore that sets
// limits on the maximum idle cons the store can use when connecting to the
// database.
func WithMaxIdleCons(n int) StoreOption {
	return func(p *Store) error {
		if n <= 0 {
			return fmt.Errorf("cannot set maxIdleCons to %d", n)
		}

		p.maxIdleCons = n
		return nil
	}
}

// WithMaxConLifetime is a function option for NewStore that sets
// limits on the maximum lifetime of the cons the store can use when
// connecting to the database.
func WithMaxConLifetime(duration time.Duration) StoreOption {
	return func(p *Store) error {
		if duration <= 0 {
			return fmt.Errorf("cannot set maxConLifetime to %s", duration)
		}

		p.maxConLifetime = duration
		return nil
	}
}
