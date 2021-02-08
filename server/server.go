package server

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/broothie/jottr/config"
	"github.com/broothie/jottr/logger"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"google.golang.org/api/option"
)

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
		render:  render.New(render.Options{Layout: "layout", Extensions: []string{".html"}}),
		log:     log,
		db:      client,
	}

	server.UseHandler(server.routes())
	return server, nil
}
