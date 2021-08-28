package common

type SysLanguageProvider interface {
	GetLanguageCode() string
	GetLanguageFontCode() string
}
