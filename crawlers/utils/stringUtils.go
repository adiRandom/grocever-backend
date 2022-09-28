package utils

type CssSelector string

func (s CssSelector) SubSelector(selector CssSelector) CssSelector {
	return CssSelector(string(s) + " " + string(selector))
}

func (s CssSelector) String() string {
	return string(s)
}
