package cnfg

import "errors"

var (
	ErrConfigRead    = errors.New("ReadInConfig")
	ErrUnmarshalRead = errors.New("err to unmarshal config ")
	ErrEnvRead       = errors.New("read env error")
	ErrUnknownDB     = errors.New("unknown datebase")
)
