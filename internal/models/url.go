package models

import "errors"

var ErrNameInvalid = errors.New("invalid parameter names in json body")

type Url struct {
	Id       uint64 `json:"id"`
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
