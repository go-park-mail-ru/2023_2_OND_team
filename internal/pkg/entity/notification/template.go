package notification

var notifyTypeTemplate = map[NotifyType]string{
	NotifyComment: `Пользователь {{.Username}} оставил комментарий под пином "{{.TitlePin}}".`,
}
