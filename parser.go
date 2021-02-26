package sugar

import (
	"encoding/json"
	"reflect"

	"github.com/vmihailenco/msgpack"
)

type IMsgParser interface {
	C2S() interface{}
	S2C() interface{}
	C2SData() []byte
	S2CData() []byte
	C2SString() string
	S2CString() string
}

type MsgParser struct {
	s2c     interface{}
	c2s     interface{}
	c2sFunc ParseFunc
	s2cFunc ParseFunc
	parser  IParser
}

func (r *MsgParser) C2S() interface{} {
	if r.c2s == nil && r.c2sFunc != nil {
		r.c2s = r.c2sFunc()
	}
	return r.c2s
}

func (r *MsgParser) S2C() interface{} {
	if r.s2c == nil && r.s2cFunc != nil {
		r.s2c = r.s2cFunc()
	}
	return r.s2c
}

func (r *MsgParser) C2SData() []byte {
	return r.parser.PackMsg(r.C2S())
}

func (r *MsgParser) S2CData() []byte {
	return r.parser.PackMsg(r.S2C())
}

func (r *MsgParser) C2SString() string {
	return string(r.C2SData())
}

func (r *MsgParser) S2CString() string {
	return string(r.S2CData())
}

type ParserType int

const (
	ParserTypePB   ParserType = iota // protoBuf, use this type of message to communicate with client
	ParserTypeJSON                   // json type
	ParserTypeCmd                    // cmd type, like telnet, console
	ParserTypeRaw                    // do not parse whatever
)

type ParseErrType int

const (
	ParseErrTypeSendRemind ParseErrType = iota // if message parsed failed, send message tips to sender, notice sender message send failed
	ParseErrTypeContinue                       // if message parsed failed, skip this message
	ParseErrTypeAlways                         // if message parsed failed, still go to nex logic
	ParseErrTypeClose                          // if message parsed failed, closed connection
)

type ParseFunc func() interface{}

type IParser interface {
	GetType() ParserType
	GetErrType() ParseErrType
	ParseC2S(msg *Message) (IMsgParser, error)
	PackMsg(v interface{}) []byte
	GetRemindMsg(err error, t MsgType) *Message
}

type Parser struct {
	Type    ParserType
	ErrType ParseErrType

	msgMap  map[int]MsgParser
	cmdRoot *cmdParseNode
	parser  IParser
}

func (r *Parser) Get() IParser {
	switch r.Type {
	case ParserTypePB:
		if r.parser == nil {
			r.parser = &pBParser{Parser: r}
		}
	case ParserTypeCmd:
		return &cmdParser{Parser: r}
	case ParserTypeRaw:
		return nil
	}

	return r.parser
}

func (r *Parser) GetType() ParserType {
	return r.Type
}

func (r *Parser) GetErrType() ParseErrType {
	return r.ErrType
}

func (r *Parser) RegisterFunc(cmd, act uint8, c2sFunc ParseFunc, s2cFunc ParseFunc) {
	if r.msgMap == nil {
		r.msgMap = map[int]MsgParser{}
	}

	r.msgMap[CmdAct(cmd, act)] = MsgParser{c2sFunc: c2sFunc, s2cFunc: s2cFunc}
}

func (r *Parser) Register(cmd, act uint8, c2s interface{}, s2c interface{}) {
	if r.msgMap == nil {
		r.msgMap = map[int]MsgParser{}
	}

	p := MsgParser{}
	if c2s != nil {
		c2sType := reflect.TypeOf(c2s).Elem()
		p.c2sFunc = func() interface{} {
			return reflect.New(c2sType).Interface()
		}
	}
	if s2c != nil {
		s2cType := reflect.TypeOf(s2c).Elem()
		p.s2cFunc = func() interface{} {
			return reflect.New(s2cType).Interface()
		}
	}

	r.msgMap[CmdAct(cmd, act)] = p
}

func (r *Parser) RegisterMsgFunc(c2sFunc ParseFunc, s2cFunc ParseFunc) {
	if r.cmdRoot == nil {
		r.cmdRoot = &cmdParseNode{}
	}
	registerCmdParser(r.cmdRoot, c2sFunc, s2cFunc)
}

func (r *Parser) RegisterMsg(c2s interface{}, s2c interface{}) {
	var c2sFunc ParseFunc
	var s2cFunc ParseFunc
	if c2s != nil {
		c2sType := reflect.TypeOf(c2s).Elem()
		c2sFunc = func() interface{} {
			return reflect.New(c2sType).Interface()
		}
	}
	if s2c != nil {
		s2cType := reflect.TypeOf(s2c).Elem()
		s2cFunc = func() interface{} {
			return reflect.New(s2cType).Interface()
		}
	}

	if r.cmdRoot == nil {
		r.cmdRoot = &cmdParseNode{}
	}
	registerCmdParser(r.cmdRoot, c2sFunc, s2cFunc)
}

func JSONUnPack(data []byte, msg interface{}) error {
	if data == nil || msg == nil {
		return ErrJSONUnPack
	}

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	return nil
}

func JSONPack(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, ErrJSONPack
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MsgPackUnPack(data []byte, msg interface{}) error {
	err := msgpack.Unmarshal(data, msg)
	return err
}

func MsgPackPack(msg interface{}) ([]byte, error) {
	data, err := msgpack.Marshal(msg)
	return data, err
}
