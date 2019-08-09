package order

import (
	"net/http"

	"github.com/SachinMaharana/isabella/util"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	serviceName       = "order"
	nextIDKey         = "nextID"
	orderKeyNamespace = serviceName
)

// Server is a wrapper for a HTTP server, with dependencies attached.
type Server struct {
	address     string
	endpoint    string
	itemService string
	redis       *redis.Client
	redisOps    uint64
	server      *http.Server
	router      *mux.Router
	logger      *util.Logger
	promReg     *prometheus.Registry
}

// ServerOptions sets options when creating a new server.
type ServerOptions func(*Server) error

// NewServer creates a new Server according to options.
func NewServer(options ...ServerOptions) (*Server, error) {
	// Create default logger
	logger, err := util.NewLogger("info", serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Logger")
	}
	// Sane defaults
	rc, _ := NewRedisClient("redis://127.0.0.1:6380/0")

	s := &Server{
		address:     ":8090",
		endpoint:    "http://127.0.0.1:8091",
		itemService: "http://127.0.0.1:8080",
		redis:       rc,
		logger:      logger,
		router:      util.NewRouter(),
		promReg:     prometheus.NewRegistry(),
	}

	// for  _ ,fn := rane options {
	// 	if err := fn(s); err != nil {
	// 		return nil, errors.Wrap(err, "failed to set server options")
	// 	}
	// }
}

// NewRedisClient creates a new go-redis/redis client according to passed options.
// Address needs to be a valid redis URL, e.g. redis://127.0.0.1:6379/0 or redis://:qwerty@localhost:6379/1
func NewRedisClient(addr string) (*redis.Client, error) {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}

	c := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})

	return c, nil
}

// Run starts a Server and shuts it down properly on a SIGINT and SIGTERM.
func (s *Server) Run() error {

}

// InitPromReg initializes a custom Prometheus registry with Collectors.
func (s *Server) InitPromReg() {
	s.promReg.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		rm.InFlightGauge, rm.Counter, rm.Duration, rm.ResponseSize,
	)
}
