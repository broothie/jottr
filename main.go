package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

type Map map[string]interface{}

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[jottr] ", log.LstdFlags)
}

func main() {
	db, err := firestore.NewClient(context.Background(), "jottr-301706")
	if err != nil {
		logger.Panic(err)
		return
	}

	jots := db.Collection(collectionName())
	router := httprouter.New()

	// Serve index for all undefined routes
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	// Ping
	router.GET("/api/ping", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := fmt.Fprint(w, "pong"); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Create jot
	router.POST("/api/jots", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		now := time.Now()
		jotID := newJotCode()
		jot := Map{
			"id":           jotID,
			"read_only_id": newJotCode(),
			"created_at":   now,
			"updated_at":   now,
		}

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

	// Show jot
	router.GET("/api/jots/:jot_id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		readOnly := false
		jotID := params.ByName("jot_id")
		doc, err := jots.Doc(jotID).Get(r.Context())
		if err != nil && status.Code(err) != codes.NotFound {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !doc.Exists() {
			docs, err := jots.Where("read_only_id", "==", jotID).Documents(r.Context()).GetAll()
			if err != nil {
				if status.Code(err) != codes.NotFound {
					w.WriteHeader(http.StatusBadRequest)
				} else {
					logger.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}

				return
			}

			doc = docs[0]
			readOnly = true
		}

		jot := doc.Data()
		jot["id"] = jot["read_only_id"] // Need to hide real jot id
		jot["editable"] = !readOnly
		if err := json.NewEncoder(w).Encode(jot); err != nil {
			logger.Println(err)
		}
	})

	// Update jot
	router.PATCH("/api/jots/:jot_id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		jot := Map{"updated_at": time.Now()}
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

	// Destroy jot
	router.DELETE("/api/jots/:jot_id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		if _, err := jots.Doc(params.ByName("jot_id")).Delete(r.Context()); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	// Bulk get jots
	router.GET("/api/bulk/jots", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jotIDs := strings.Split(r.URL.Query().Get("jot_ids"), ",")
		refs := make([]*firestore.DocumentRef, len(jotIDs))
		for i, jotId := range jotIDs {
			refs[i] = jots.Doc(jotId)
		}

		var docs []*firestore.DocumentSnapshot
		if err := db.RunTransaction(r.Context(), func(_ context.Context, transaction *firestore.Transaction) error {
			var err error
			if docs, err = transaction.GetAll(refs); err != nil {
				logger.Println(err)
				return err
			}

			return nil
		}); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jots := make([]Map, 0, len(docs))
		for _, doc := range docs {
			if !doc.Exists() {
				continue
			}

			var jot Map
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

	// Jot purge job
	router.GET("/jobs/purge", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := db.RunTransaction(r.Context(), func(_ context.Context, transaction *firestore.Transaction) error {
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
		}); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run()
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

func newJotCode() string {
	return fmt.Sprintf("%s-%s-%s", randomLetters(3), randomLetters(4), randomLetters(3))
}

func randomLetters(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	runes := make([]rune, length)
	for i := 0; i < length; i++ {
		index, err := randomInt(len(alphabet))
		if err != nil {
			logger.Println(err)
			continue
		}

		runes[i] = rune(alphabet[index])
	}

	return string(runes)
}

func randomInt(max int) (int, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(i.Int64()), nil
}
