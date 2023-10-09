package controllers

type ControllerMetaData struct {
	Name            string
	Path            []string
	LoggerName      string
	JsonContentType string
	AllowedMethods  []string
}
