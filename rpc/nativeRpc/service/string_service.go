package service

import (
	"errors"
	"strings"
	"time"
)

type StringRequest struct {
	A string
	B string
}

type Service interface {
	// Concat 连接两个字符串.
	Concat(req StringRequest, ret *string) error
	// Diff 寻找相同的字串.
	Diff(req StringRequest, ret *string) error
}

type StringService struct {
}

var (
	StrMaxSize = 1024
	ErrMaxSize = errors.New("the string is to long")
)

func (s StringService) Concat(req StringRequest, ret *string) error {
	if len(req.A)+len(req.B) > StrMaxSize {
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}

func (s StringService) Diff(req StringRequest, ret *string) error {
	time.Sleep(5 * time.Second)
	if len(req.A) < 1 || len(req.B) < 1 {
		*ret = ""
		return nil
	}

	res := ""
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.A, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.B, string(char)) {
				res = res + string(char)
			}
		}
	}
	*ret = res
	return nil
}
