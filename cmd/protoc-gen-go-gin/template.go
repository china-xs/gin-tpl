// Package main
// @author: xs
// @date: 2022/3/3
// @Description: main,描述
package main

import (
	"bytes"
	"strings"
	"text/template"
)

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}

// gin 模版 暂时 不使用
var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}GinServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}GinServer(s *gin_tpl.Server, srv {{.ServiceType}}GinServer,ms ...gin.HandlerFunc) {
	route :=s.Engine.Use(ms...)
	{{- range .Methods}}
	route.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_Gin_Handler(s,srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_Gin_Handler(s *gin_tpl.Server,srv {{$svrType}}GinServer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var in {{.Request}}

		{{- if .HasBody}}
		if err := c.ShouldBind(&in{{.Body}}); err != nil {
			return 
		}
		{{- end}}
		
		{{- if eq .Method "GET" "DELETE" }}
		if err := c.ShouldBindQuery(&in); err != nil {
			return 
		}
		{{- end}}
		
		{{- if .HasVars}}
		if err := c.ShouldBindUri(&in); err != nil {
			return 
		}
		{{- end}}


		out,err := srv.{{.Name}}(c, &in)
		s.Enc(c,out,err)
		return
		//reply := out.(*{{.Reply}})
		//return ctx.Result(200, reply{{.ResponseBody}})
	}
}
{{end}}
`
