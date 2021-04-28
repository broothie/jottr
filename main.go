package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PublicJot struct {
	ID        string      `json:"id" firestore:"id"`
	Delta     interface{} `json:"delta" firestore:"delta"`
	Title     string      `json:"title" firestore:"title"`
	CreatedAt time.Time   `json:"-" firestore:"created_at"`
	UpdatedAt time.Time   `json:"-" firestore:"updated_at"`
}

type Jot struct {
	PublicJot
	ReadOnlyID string `json:"read_only_id" firestore:"read_only_id"`
}

func main() {
	logger := log.New(os.Stdout, "[jottr] ", log.LstdFlags)
	db, err := firestore.NewClient(context.Background(), "jottr-301706")
	if err != nil {
		log.Panic(err)
		return
	}

	jots := db.Collection(collectionName())
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	router.GET("/api/ping", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := fmt.Fprint(w, "pong"); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	router.POST("/api/jots", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		now := time.Now()
		jotID := newJotCode()
		jot := Jot{PublicJot: PublicJot{ID: jotID, CreatedAt: now, UpdatedAt: now}, ReadOnlyID: newJotCode()}

		if _, err := jots.Doc(jotID).Set(r.Context(), jot); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(jot); err != nil {
			logger.Println(err)
		}
	})

	router.GET("/api/jots/:jot_id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		doc, err := jots.Doc(params.ByName("jot_id")).Get(r.Context())
		if err != nil {
			logger.Println(err)
			code := http.StatusInternalServerError
			if status.Code(err) == codes.NotFound {
				code = http.StatusBadRequest
			}

			w.WriteHeader(code)
			return
		}

		var jot Jot
		if err := doc.DataTo(&jot); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(jot); err != nil {
			logger.Println(err)
		}
	})

	router.PATCH("/api/jots/:jot_id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		jot := map[string]interface{}{"updated_at": time.Now()}
		if err := json.NewDecoder(r.Body).Decode(&jot); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := jots.Doc(params.ByName("jot_id")).Set(r.Context(), jot, firestore.MergeAll); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	router.DELETE("/api/jots/:jot_id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		if _, err := jots.Doc(params.ByName("jot_id")).Delete(r.Context()); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	router.GET("/api/bulk/jots", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jotIDs := strings.Split(r.URL.Query().Get("jot_ids"), ",")
		refs := make([]*firestore.DocumentRef, len(jotIDs))
		for i, jotId := range jotIDs {
			refs[i] = jots.Doc(jotId)
		}

		var docs []*firestore.DocumentSnapshot
		err := db.RunTransaction(r.Context(), func(ctx context.Context, transaction *firestore.Transaction) error {
			var err error
			if docs, err = transaction.GetAll(refs); err != nil {
				logger.Println(err)
				return err
			}

			return nil
		})
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jots := make([]Jot, 0, len(docs))
		for _, doc := range docs {
			if !doc.Exists() {
				continue
			}

			var jot Jot
			if err := doc.DataTo(&jot); err != nil {
				logger.Println(err)
				continue
			}

			jots = append(jots, jot)
		}

		if err := json.NewEncoder(w).Encode(jots); err != nil {
			logger.Println(err)
		}
	})

	router.GET("/jobs/purge", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		err := db.RunTransaction(r.Context(), func(ctx context.Context, transaction *firestore.Transaction) error {
			docs, err := transaction.Documents(jots.Where("title", "==", "")).GetAll()
			if err != nil {
				logger.Println(err)
				return err
			}

			for _, doc := range docs {
				logger.Printf("deleting doc '%s'", doc.Ref.ID)
				if err := transaction.Delete(doc.Ref); err != nil {
					logger.Printf("error deleting doc '%s: %v", doc.Ref.ID, err)
				}
			}

			return nil
		})
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run()
}

func newJotCode() string {
	return fmt.Sprintf("%s-%s-%s", randomLetters(3), randomLetters(4), randomLetters(3))
}

func randomLetters(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	runes := make([]rune, length)
	for i := 0; i < length; i++ {
		runes[i] = rune(alphabet[rand.Intn(len(alphabet))])
	}

	return string(runes)
}

func collectionName() string {
	if os.Getenv("ENVIRONMENT") == "production" {
		return "production.jots"
	} else {
		devName := os.Getenv("DEV_NAME")
		if devName == "" {
			devName = "everyone"
		}

		return fmt.Sprintf("development.%s.jots", devName)
	}
}
