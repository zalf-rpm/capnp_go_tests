module example.com/test

go 1.17

require (
	capnproto.org/go/capnp/v3 v3.0.0-alpha.1
	example.com/capnp_schemas v0.0.0-00010101000000-000000000000
	example.com/greetings v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	rsc.io/quote v1.5.2
)

require (
	golang.org/x/text v0.3.3 // indirect
	rsc.io/sampler v1.3.0 // indirect
)

replace example.com/greetings => ./greetings

replace example.com/capnp_schemas => ./capnp_schemas
