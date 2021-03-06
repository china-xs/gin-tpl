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

func Register{{.ServiceType}}GinServer(s *gin_tpl.Server, srv {{.ServiceType}}GinServer,ms ...gin.HandlerFunc) {
	{{- range .Methods}}
	hdl_{{.Name}}{{.Num}} := append(ms,_{{$svrType}}_{{.Name}}{{.Num}}_Gin_Handler(s,srv))
	s.Engine.{{.Method}}("{{.Path}}", hdl_{{.Name}}{{.Num}}...)
	{{- end}}
}

func _{{.ServiceType}}_getBindBodyType(c *gin.Context) binding.BindingBody{
	b := binding.Default(c.Request.Method, c.ContentType())
	var bin binding.BindingBody
	switch b.Name() {
	case "json":
		bin = binding.JSON
	case "xml":
		bin = binding.XML
	case "yaml":
		bin = binding.YAML
	case "protobuf":
		bin = binding.ProtoBuf
	case "msgpack":
		bin = binding.MsgPack
	default:
		bin = binding.JSON
	}
	return bin
}


{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_Gin_Handler(s *gin_tpl.Server,srv {{$svrType}}GinServer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var in {{.Request}}
		switch c.Request.Method {
			case "POST","PUT":
			bin := _{{$svrType}}_getBindBodyType(c)
			if err := c.ShouldBindBodyWith(&in,bin); err!=nil{
				s.Enc(c,nil,err)
				return 
			}
			if strings.Contains(c.Request.URL.String(),"?"){
				if err := c.ShouldBindQuery(&in); err != nil {
					s.Enc(c,nil,err)
					return 
				}
			}
			case "GET","DELETE":
			if err := c.ShouldBindQuery(&in); err != nil {
				s.Enc(c,nil,err)
				return 
			}
		}
		{{- if .HasVars}}
		if err := c.ShouldBindUri(&in); err != nil {
			s.Enc(c,nil,err)
			return 
		}
		{{- end}}
		c.Set(gin_tpl.OperationKey, "/{{$svrName}}/{{.Name}}")
		h := s.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return srv.{{.Name}}(c, req.(*{{.Request}}))
		})
		out, err := h(c, &in)
		s.Enc(c,out,err)
		return
	}
}
{{end}}
`
