package module

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	authzv1 "github.com/gh1st/protoc-gen-authorize/gen/authz/v1"
)

type module struct {
	*pgs.ModuleBase
	pgsgo.Context

	authorizer string
}

func New() pgs.Module {
	return &module{ModuleBase: &pgs.ModuleBase{}}
}

func (m *module) Name() string {
	return "authorize"
}

func (m *module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
	m.authorizer = strings.ToLower(m.authorizer)
}

func (m *module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		if f.BuildTarget() {
			m.generate(f)
		}
	}
	return m.Artifacts()
}

func (m *module) generate(f pgs.File) {
	var rules = map[string]*authzv1.AuthOptions{}
	for _, s := range f.Services() {
		for _, method := range s.Methods() {
			var ruleSet authzv1.AuthOptions
			ok, err := method.Extension(authzv1.E_AuthOptions, &ruleSet)
			if err != nil {
				m.AddError(err.Error())
				continue
			}
			if !ok {
				continue
			}
			// EchoService_Echo_FullMethodName
			name := fmt.Sprintf("%s_%s_FullMethodName", s.Name().UpperCamelCase(), method.Name().UpperCamelCase())
			rules[name] = &ruleSet
		}
	}
	if len(rules) == 0 {
		return
	}
	name := f.InputPath().SetExt(".pb.authz.go").String()

	t, err := template.New("authz").Parse(tmpl)
	if err != nil {
		m.AddError(err.Error())
		return
	}

	buffer := &bytes.Buffer{}
	if err := t.Execute(buffer, templateData{
		Package: m.Context.PackageName(f).String(),
		Rules:   rules,
	}); err != nil {
		m.AddError(err.Error())
		return
	}
	m.AddGeneratorFile(name, buffer.String())
}

type templateData struct {
	Package string
	Rules   map[string]*authzv1.AuthOptions
}

var tmpl = `
package {{ .Package }}

import (
	authzv1 "github.com/gh1st/protoc-gen-authorize/gen/authz/v1"
)

func NewAuthorizer() (map[string]*authzv1.AuthOptions, error) {
	return map[string]*authzv1.AuthOptions{
	{{- range $key, $value := .Rules }}
		{{$key}}: &authzv1.AuthOptions {
			Public: {{ $value.Public }},
			RequiredRoles: []authzv1.Role {
			{{- range $value.RequiredRoles }}
				authzv1.NewRole("{{ . }}"),
			{{- end }}
			},
		},
	{{- end }}
	}, nil
}
`
