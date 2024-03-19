templ:
	templ generate --watch --proxy=http://localhost:8080

schema:
	sqlite3 test.db .read schema.sql
