module github.com/tarik/nessmpp

go 1.21.0

require (
	github.com/tarik/nessmpp/pkg/pdu v0.0.0-00010101000000-000000000000
	golang.org/x/time v0.10.0
)

replace github.com/tarik/nessmpp/pkg/pdu => ./pkg/pdu
