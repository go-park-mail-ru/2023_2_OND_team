package notification

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

type NotifyType uint8

const _defaultCapBuffer = 128

const (
	_ NotifyType = iota
	NotifyComment

	_notifyCustom
)

type notify struct {
	NotifyType NotifyType
	buf        *sync.Pool
	tmp        *template.Template
}

func NewWithTemplate(tmp *template.Template) notify {
	return notify{
		NotifyType: _notifyCustom,
		buf: &sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, _defaultCapBuffer)) },
		},
		tmp: tmp,
	}
}

func NewWithType(t NotifyType) (notify, error) {
	content, ok := notifyTypeTemplate[t]
	if !ok {
		return notify{}, fmt.Errorf("new notify with type %s: %w", TypeString(t), ErrUnknownNotifyType)
	}

	res := notify{
		NotifyType: t,
		buf: &sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, _defaultCapBuffer)) },
		},
	}

	tmp, err := template.New(TypeString(t)).Parse(content)
	if err != nil {
		return notify{}, fmt.Errorf("new notify with type %s: %w", TypeString(t), err)
	}

	res.tmp = tmp
	return res, nil
}

func (n notify) Type() NotifyType {
	return n.NotifyType
}

func (n notify) BuildNotifyMessage(data any) (*NotifyMessage, error) {
	content, err := n.FormatContent(data)
	if err != nil {
		return nil, fmt.Errorf("build notify message: %w", err)
	}

	return NewNotifyMessage(n.NotifyType, content), nil
}

func (n notify) FormatContent(data any) (string, error) {
	buf := n.buf.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		n.buf.Put(buf)
	}()

	err := n.tmp.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("")
	}

	return buf.String(), nil
}
