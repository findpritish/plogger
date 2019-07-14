package loger

import (
	"regexp"
	"text/template"
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

// Config contains the config for stern
type Config struct {
	KubeConfig            string
	ContextName           string
	Namespace             string
	PodQuery              *regexp.Regexp
	Timestamps            bool
	ContainerQuery        *regexp.Regexp
	ExcludeContainerQuery *regexp.Regexp
	ContainerState        ContainerState
	Exclude               []*regexp.Regexp
	Since                 time.Duration
	AllNamespaces         bool
	LabelSelector         labels.Selector
	TailLines             *int64
	Template              *template.Template
}
