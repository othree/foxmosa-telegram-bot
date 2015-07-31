
all: foxmosa

foxmosa: foxmosa.go offset.go telegram_to_pierc.go foxmosa-sticker.go
	go build foxmosa.go offset.go telegram_to_pierc.go foxmosa-sticker.go
