package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/util/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Salt = "15m.games"
)

func addSalt(pw string) string {
	return pw + Salt
}

func saltPw(pw string) string {
	bytes := md5.Sum([]byte(addSalt(pw)))
	return hex.EncodeToString(bytes[:])
}

func (s *Svc) OnSignUp(pkt kiwi.IRcvRequest, req *pb.SignUpReq, res *pb.SignUpRes) {
	if idLen := len(req.Id); idLen < 4 || idLen > 32 {
		pkt.Fail(EcSignUp_IdWrong)
		return
	}
	if pwLen := len(req.Password); pwLen < 8 || pwLen > 32 {
		pkt.Fail(EcSignUp_PwWrong)
		return
	}
	_, err := mgo.InsertOne(SchemaUser, &pb.User{
		Id:       req.Id,
		Password: saltPw(req.Password),
		RoleMask: util.GenMask(common.RPlayer),
	})
	if err != nil {
		pkt.Fail(EcSignUp_IdExist)
		return
	}
	addr, _ := util.MGet[string](pkt.Head(), "addr")
	tkn, _ := common.GenToken(req.Id, addr)
	res.Token = tkn
	common.ReqGateAddrUpdate(pkt.Tid(), pkt.SenderId(), addr, util.M{
		common.HdMask:   util.GenMask(common.RPlayer),
		common.HdUserId: req.Id,
	}, util.M{
		common.CcToken: tkn,
	})
	pkt.Ok(res)
}

func (s *Svc) OnSignIn(pkt kiwi.IRcvRequest, req *pb.SignInReq, res *pb.SignInRes) {
	if idLen := len(req.Id); idLen < 4 || idLen > 32 {
		pkt.Fail(EcSignIn_WrongIdOrPassword)
		return
	}
	if pwLen := len(req.Password); pwLen < 8 || pwLen > 32 {
		pkt.Fail(EcSignIn_WrongIdOrPassword)
		return
	}
	addr, _ := util.MGet[string](pkt.Head(), common.HdGateAddr)
	var user pb.User
	e := mgo.FindOne(SchemaUser, bson.D{
		{Id, req.Id},
	}, &user)
	if e != nil {
		if errors.Is(e, mongo.ErrNoDocuments) {
			pkt.Fail(EcSignIn_WrongIdOrPassword)
			return
		}
		pkt.Err(util.WrapErr(util.EcDbErr, e))
		return
	}
	if user.Password != saltPw(req.Password) {
		pkt.Fail(EcSignIn_WrongIdOrPassword)
		return
	}
	tkn, _ := common.GenToken(req.Id, addr)
	res.Token = tkn
	common.ReqGateAddrUpdate(pkt.Tid(), pkt.SenderId(), addr, util.M{
		common.HdMask:   user.RoleMask,
		common.HdUserId: req.Id,
	}, util.M{
		common.CcToken: tkn,
	})
	pkt.Ok(res)
}

func (s *Svc) OnSignOut(pkt kiwi.IRcvRequest, req *pb.SignOutReq, res *pb.SignOutRes) {
	pkt.Ok(res)

	_ = core.Ntf(pkt.Tid(), nil, &pb.SignOutNtc{Id: pkt.HeadId()})
}
