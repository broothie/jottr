package server

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/broothie/jottr/config"
	"github.com/broothie/jottr/logger"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"google.golang.org/api/option"
)

var whitespaceMatcher = regexp.MustCompile(`^\s+$`)

type Jot struct {
	ID       string      `json:"id" firestore:"id"`
	Body     string      `json:"body" firestore:"body"`
	Contents interface{} `json:"contents" firestore:"contents"`
}

func (j Jot) LinkText() string {
	if whitespaceMatcher.MatchString(j.Body) {
		return j.ID
	}

	return strings.Split(j.Body, "\n")[0]
}

type Server struct {
	*negroni.Negroni
	config config.Config
	render *render.Render
	log    *logger.Logger
	db     *firestore.Client
}

func New(cfg config.Config, log *logger.Logger) (*Server, error) {
	var options []option.ClientOption
	if cfg.Environment != "production" {
		options = append(options, option.WithCredentialsFile("gcloud-key.json"))
	}

	client, err := firestore.NewClient(context.Background(), "jottr-301706", options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	server := &Server{
		Negroni: negroni.Classic(),
		config:  cfg,
		render: render.New(render.Options{
			Layout:     "layout",
			Extensions: []string{".html"},
			Funcs:      []template.FuncMap{{"json": marshalJSONString}},
		}),
		log: log,
		db:  client,
	}

	server.UseHandler(server.routes())
	return server, nil
}

func (s *Server) Error(w http.ResponseWriter, err error, message string, code int, logFields ...logger.Fieldser) {
	s.log.Err(err, message, logger.Field("user_message", message))
	http.Error(w, message, code)
}

func noCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func marshalJSONString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}
