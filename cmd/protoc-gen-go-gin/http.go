// Package main
// @author: xs
// @date: 2022/3/3
// @Description: main gin-http
package main

import (
	"fmt"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"regexp"
	"strings"
)

const (
	contextPkg = protogen.GoImportPath("context")
	ginPkg     = protogen.GoImportPath("github.com/gin-gonic/gin")
	ginSev     = protogen.GoImportPath("github.com/china-xs/gin-tpl")
	ginBind    = protogen.GoImportPath("github.com/gin-gonic/gin/binding")
	pkgStr     = protogen.GoImportPath("strings")
)

var methodSets = make(map[string]int)

// generateFile generates a _gin.pb.go file
func generateFile(gen *protogen.Plugin, file *protogen.File, omitempty bool) *protogen.GeneratedFile {
	if len(file.Services) == 0 || (omitempty && !hasHTTPRule(file.Services)) {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_gin.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-gin. DO NOT EDIT.")
	g.P(fmt.Sprintf("// versions:%s", release))
	g.P(fmt.Sprintf("// protoc-gen-go-gin %s", release))
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(gen, file, g, omitempty)
	return g
}

// generateFileContent generates the kratos errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, omitempty bool) {
	if len(file.Services) == 0 {
		return
	}
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the kratos package it is being compiled against.")
	g.P("// ",
		ginPkg.Ident(""),
		contextPkg.Ident(""),
		ginSev.Ident(""),
		ginBind.Ident(""),
		pkgStr.Ident(""),
	)
	//g.P("var _ = new(", contextPackage.Ident("Context"), ")")
	//g.P("var _ = ", bindingPackage.Ident("EncodeURL"))
	//g.P("const _ = ", transportHTTPPackage.Ident("SupportPackageIsVersion1"))
	//g.P()

	for _, service := range file.Services {
		genService(gen, file, g, service, omitempty)
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	// HTTP Server.
	sd := &serviceDesc{
		ServiceType: service.GoName,
		ServiceName: string(service.Desc.FullName()),
		Metadata:    file.Desc.Path(),
	}
	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}
		rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			for _, bind := range rule.AdditionalBindings {
				sd.Methods = append(sd.Methods, buildHTTPRule(g, method, bind))
			}
			sd.Methods = append(sd.Methods, buildHTTPRule(g, method, rule))
		} else if !omitempty {
			path := fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name())
			sd.Methods = append(sd.Methods, buildMethodDesc(g, method, "POST", path))
		}
	}
	if len(sd.Methods) != 0 {
		g.P(sd.execute())
	}
}

func hasHTTPRule(services []*protogen.Service) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
			if rule != nil && ok {
				return true
			}
		}
	}
	return false
}

func buildHTTPRule(g *protogen.GeneratedFile, m *protogen.Method, rule *annotations.HttpRule) *methodDesc {
	var (
		path         string
		method       string
		body         string
		responseBody string
		hasVars      bool
	)
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		if hasPathParams(path) {
			hasVars = true
			path = initPathParams(path)
		}
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		if hasPathParams(path) {
			hasVars = true
			path = initPathParams(path)
		}
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	body = rule.Body
	responseBody = rule.ResponseBody
	md := buildMethodDesc(g, m, method, path)
	if hasVars && !md.HasVars {
		md.HasVars = true
	}
	if method == "GET" || method == "DELETE" {
		if body != "" {
			_, _ = fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s body should not be declared.\n", method, path)
		}
		md.HasBody = false
	} else if body == "*" {
		md.HasBody = true
		md.Body = ""
	} else if body != "" {
		md.HasBody = true
		md.Body = "." + camelCaseVars(body)
	} else {
		md.HasBody = false
		_, _ = fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s does not declare a body.\n", method, path)
	}
	if responseBody == "*" {
		md.ResponseBody = ""
	} else if responseBody != "" {
		md.ResponseBody = "." + camelCaseVars(responseBody)
	}
	return md
}

//
// hasPathParams
// @Description: 判断路由是否有参数，暂时仅对get|delete 做此操作调整，骚操作，post、pust 暂不处理
// @param path
// @return bool
//
func hasPathParams(path string) bool {
	paths := strings.Split(path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}
	return false
}

//
// initPathParams
// @Description: 转换路由参数
// @param oldPath
// @return string
//
func initPathParams(oldPath string) string {
	paths := strings.Split(oldPath, "/")
	for i, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}
	return strings.Join(paths, "/")
}

func buildMethodDesc(g *protogen.GeneratedFile, m *protogen.Method, method, path string) *methodDesc {
	defer func() { methodSets[m.GoName]++ }()

	vars := buildPathVars(path)
	fields := m.Input.Desc.Fields()

	for v, s := range vars {
		if s != nil {
			path = replacePath(v, *s, path)
		}
		for _, field := range strings.Split(v, ".") {
			if strings.TrimSpace(field) == "" {
				continue
			}
			if strings.Contains(field, ":") {
				field = strings.Split(field, ":")[0]
			}
			fd := fields.ByName(protoreflect.Name(field))
			if fd == nil {
				fmt.Fprintf(os.Stderr, "\u001B[31mERROR\u001B[m: The corresponding field '%s' declaration in message could not be found in '%s'\n", v, path)
				os.Exit(2)
			}
			if fd.IsMap() {
				fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: The field in path:'%s' shouldn't be a map.\n", v)
			} else if fd.IsList() {
				fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: The field in path:'%s' shouldn't be a list.\n", v)
			} else if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
				fields = fd.Message().Fields()
			}
		}
	}
	return &methodDesc{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: g.QualifiedGoIdent(m.Input.GoIdent),
		Reply:   g.QualifiedGoIdent(m.Output.GoIdent),
		Path:    path,
		Method:  method,
		HasVars: len(vars) > 0,
	}
}

func camelCaseVars(s string) string {
	vars := make([]string, 0)
	subs := strings.Split(s, ".")
	for _, sub := range subs {
		vars = append(vars, camelCase(sub))
	}
	return strings.Join(vars, ".")
}

func buildPathVars(path string) (res map[string]*string) {
	res = make(map[string]*string)
	pattern := regexp.MustCompile(`(?i){([a-z\.0-9_\s]*)=?([^{}]*)}`)
	matches := pattern.FindAllStringSubmatch(path, -1)
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if len(name) > 1 && len(m[2]) > 0 {
			res[name] = &m[2]
		} else {
			res[name] = nil
		}
	}
	return
}

func replacePath(name string, value string, path string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?i){([\s]*%s[\s]*)=?([^{}]*)}`, name))
	idx := pattern.FindStringIndex(path)
	if len(idx) > 0 {
		path = fmt.Sprintf("%s{%s:%s}%s",
			path[:idx[0]], // The start of the match
			name,
			strings.ReplaceAll(value, "*", ".*"),
			path[idx[1]:],
		)
	}
	return path
}

// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
