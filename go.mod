module isucon9final

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/sessions v1.2.1
	github.com/hirosuzuki/go-isucon-tracer v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.2.0
	goji.io v2.0.2+incompatible
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
)

replace github.com/hirosuzuki/go-isucon-tracer => ../go-isucon-tracer
