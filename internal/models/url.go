package models

import "errors"

var ErrNameInvalid = errors.New("invalid parameter names in json body")

type ShortenedUrl struct {
	ID       uint64 `json:"id"`
	ShortUrl string `json:"shortUrl"`
	LongUrl  string `json:"longUrl"`
}

type LongUrl struct {
	LongUrl string `json:"longUrl"`
}

func (l LongUrl) Validation() error {
	switch {
	case len(l.LongUrl) == 0:
		return ErrNameInvalid
	default:
		return nil
	}
}
