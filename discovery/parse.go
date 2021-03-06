package discovery

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/appc/spec/schema/types"
)

type App struct {
	Name   types.ACIdentifier
	Labels map[types.ACIdentifier]string
}

func NewApp(name string, labels map[types.ACIdentifier]string) (*App, error) {
	if labels == nil {
		labels = make(map[types.ACIdentifier]string, 0)
	}
	acn, err := types.NewACIdentifier(name)
	if err != nil {
		return nil, err
	}
	return &App{
		Name:   *acn,
		Labels: labels,
	}, nil
}

// NewAppFromString takes a command line app parameter and returns a map of labels.
//
// Example app parameters:
// 	example.com/reduce-worker:1.0.0
// 	example.com/reduce-worker,channel=alpha,label=value
func NewAppFromString(app string) (*App, error) {
	var (
		name   string
		labels map[types.ACIdentifier]string
	)

	app = strings.Replace(app, ":", ",version=", -1)
	app = "name=" + app
	v, err := url.ParseQuery(strings.Replace(app, ",", "&", -1))
	if err != nil {
		return nil, err
	}
	labels = make(map[types.ACIdentifier]string, 0)
	for key, val := range v {
		if len(val) > 1 {
			return nil, fmt.Errorf("label %s with multiple values %q", key, val)
		}
		if key == "name" {
			name = val[0]
			continue
		}
		labelName, err := types.NewACIdentifier(key)
		if err != nil {
			return nil, err
		}
		labels[*labelName] = val[0]
	}
	a, err := NewApp(name, labels)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Copy() *App {
	ac := &App{
		Name:   a.Name,
		Labels: make(map[types.ACIdentifier]string, 0),
	}
	for k, v := range a.Labels {
		ac.Labels[k] = v
	}
	return ac
}

// String returns the URL-like image name
func (a *App) String() string {
	img := a.Name.String()
	for n, v := range a.Labels {
		img += fmt.Sprintf(",%s=%s", n, v)
	}
	return img
}
