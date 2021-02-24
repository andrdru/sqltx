module github.com/andrdru/sqltx/examples/sqlitetx

go 1.15

replace github.com/andrdru/sqltx v0.0.0 => ./../../../sqltx

require (
	github.com/andrdru/sqltx v0.0.0
	github.com/mattn/go-sqlite3 v1.14.6
)
