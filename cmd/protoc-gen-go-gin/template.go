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
	{{.Name}}(*gin.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}GinServer(s *gin.Engine, srv {{.ServiceType}}GinServer,ms ...gin.HandlerFunc) {
	s.Use(ms...)
	{{- range .Methods}}
	s.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_Gin_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_Gin_Handler(srv {{$svrType}}GinServer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var in {{.Request}}
		{{- if .HasBody}}
		if err := c.ShouldBindJSON(&in{{.Body}}); err != nil {
			return 
		}
		
		{{- if not (eq .Body "")}}
		if err := c.ShouldBind(&in); err != nil {
			return 
		}
		{{- end}}
		{{- else}}
		if err := c.BindQuery(&in{{.Body}}); err != nil {
			return 
		}
		{{- end}}
		{{- if .HasVars}}
		if err := ctx.BindVars(&in); err != nil {
			return 
		}
		{{- end}}
		out,err := srv.{{.Name}}(c, &in)
		if err != nil {
			return 
		}
		c.JSON(200,out)
		return
		//reply := out.(*{{.Reply}})
		//return ctx.Result(200, reply{{.ResponseBody}})
	}
}
{{end}}
`
