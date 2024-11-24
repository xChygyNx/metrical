module main

go 1.22.9

replace agent => ./agent

replace server => ./server

replace agent/senders => ./agent/senders

require (
	agent/senders v0.0.0-00010101000000-000000000000 // indirect
	server v0.0.0-00010101000000-000000000000 // indirect
)
