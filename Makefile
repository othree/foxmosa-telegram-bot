
all: foxmosa

foxmosa: foxmosa.go offset.go
	go build foxmosa.go offset.go
