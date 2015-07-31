
all: foxmosa

foxmosa: foxmosa.go writeoffset.go
	go build foxmosa.go writeoffset.go
