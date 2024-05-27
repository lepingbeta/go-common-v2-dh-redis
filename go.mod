module github.com/lepingbeta/go-common-v2-dh-redis

replace (
	github.com/lepingbeta/go-common-v2-dh-json => ../go-common-v2-dh-json
	github.com/lepingbeta/go-common-v2-dh-log => ../go-common-v2-dh-log
// github.com/lepingbeta => ../
)

go 1.22.1

require (
	github.com/gomodule/redigo v1.9.2
	github.com/lepingbeta/go-common-v2-dh-log v0.0.0-00010101000000-000000000000
)
