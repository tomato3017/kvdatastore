package io

import (
	"github.com/rs/zerolog/log"
	"io"
)

func SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Err(err).Msg("Unable to close stream")
	}
}
