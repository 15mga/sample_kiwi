package common

import (
	"encoding/base64"
	"errors"
	"github.com/15mga/kiwi/util"
	"github.com/golang-jwt/jwt/v4"
)

var (
	Issuer          = "15m.games"
	TokenSecret, _  = base64.URLEncoding.DecodeString("95eh.com")
	_ErrWrongIssuer = errors.New("wrong issuer")
)

type Claims struct {
	Issuer string
	Id     string
	Addr   string
}

func (c *Claims) Valid() error {
	if c.Issuer != Issuer {
		return _ErrWrongIssuer
	}
	return nil
}

func GenToken(id, addr string) (string, *util.Err) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Issuer: Issuer,
		Id:     id,
		Addr:   addr,
	})
	tkn, e := token.SignedString(TokenSecret)
	if e != nil {
		return "", util.WrapErr(util.EcServiceErr, e)
	}
	return tkn, nil
}

func ParseToken(tkn string) (id, addr string, err *util.Err) {
	token, e := jwt.ParseWithClaims(tkn, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return TokenSecret, nil
	})
	if e != nil {
		err = util.NewErr(util.EcIllegalOp, util.M{
			"token": tkn,
			"error": e.Error(),
		})
		return
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		err = util.NewErr(util.EcIllegalOp, util.M{
			"token": tkn,
		})
		return
	}
	id, addr = claims.Id, claims.Addr
	return
}
