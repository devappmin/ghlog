module main

go 1.17

replace main.com/github => ./github

replace main.com/auth => ./auth

replace main.com/util => ./util

require (
	main.com/auth v0.0.0-00010101000000-000000000000
	main.com/github v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/go-github/v41 v41.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)
