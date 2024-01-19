package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AlexTLDR/ByteVault/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	/* SQL query for sessions
		CREATE TABLE sessions (
	    token CHAR(43) PRIMARY KEY,
	    data BLOB NOT NULL,
	    expiry TIMESTAMP(6) NOT NULL
		);

		CREATE INDEX sessions_expiry_idx ON sessions (expiry);
	*/
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "rick:morty@/bytevault?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	// For a cleaner main() function, I've placed the connection pool creation code in a separate openDB() function below,
	// passing it the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	// use the scs.New() function to initialize a new session manager

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// setting the Secure attribute to session cookies
	sessionManager.Cookie.Secure = true

	// And add the session manager to our application dependencies.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	logger.Info("starting server", "addr", *addr)
	// start the HTTPS server with ListenAndServeTLS()
	// TODO: switch to caddy probably, in prod
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function uses sql.Open() to create a sql.DB connection pool based on the provided DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
