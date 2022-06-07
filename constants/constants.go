package constants

import (
	"time"
)

const FetchVoiceDataInterval time.Duration = time.Duration(10) * time.Second
const ProcessVoiceDataInterval time.Duration = time.Duration(60) * time.Second

const FetchMessageDataInterval time.Duration = time.Duration(60) * time.Second
const ProcessMessageDataInterval time.Duration = time.Duration(240) * time.Second

const FetchNicknamesInterval time.Duration = time.Duration(240) * time.Second
