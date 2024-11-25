module main

go 1.22.9

replace agent => ./agent

replace server => ./server

require server v0.0.0-00010101000000-000000000000 // indirect
