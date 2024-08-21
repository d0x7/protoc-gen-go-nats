package main

import (
	"fmt"
	pb "gitea.xiam.li/Hydria/protoc-gen-go-nats/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"strconv"
	"strings"
)

const (
	natsPkg   = protogen.GoImportPath("github.com/nats-io/nats.go")
	microPkg  = protogen.GoImportPath("github.com/nats-io/nats.go/micro")
	timePkg   = protogen.GoImportPath("time")
	protoPkg  = protogen.GoImportPath("google.golang.org/protobuf/proto")
	protocPkg = protogen.GoImportPath("gitea.xiam.li/Hydria/protoc-gen-go-nats")
)

var (
	natsConn       = natsPkg.Ident("Conn")
	microConfig    = microPkg.Ident("Config")
	microRequest   = microPkg.Ident("Request")
	timeDuration   = timePkg.Ident("Duration")
	protoMessage   = protoPkg.Ident("Message")
	protoMarshal   = protoPkg.Ident("Marshal")
	protoUnmarshal = protoPkg.Ident("Unmarshal")

	protocErrMarshalling   = protocPkg.Ident("ErrMarshallingFailed")
	protocErrUnmarshalling = protocPkg.Ident("ErrUnmarshallingFailed")
)

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Services) == 0 {
		return
	}
	filename := file.GeneratedFilenamePrefix + "_nats.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-nats. DO NOT EDIT.")
	g.P("// Versions:")
	g.P("// - protoc-gen-go-nats v", version)
	g.P("// - protoc        v", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	for _, service := range file.Services {
		generateService(g, service)
	}
}

func generateServer(g *protogen.GeneratedFile, service *protogen.Service) {
	srvName := service.GoName + "NATSServer"

	// Generate server interface
	g.P("type ", srvName, " interface {")
	for _, method := range service.Methods {
		g.P(method.Comments.Leading, methodSignature(g, method))
	}
	g.P("}")
	g.P()

	// Generate server options
	g.P("type ServerOption func(config *", g.QualifiedGoIdent(microConfig), ")")

	// Generate server option functions
	generateServerOptionHandler(g, "StatsHandler", "StatsHandler")
	generateServerOptionHandler(g, "DoneHandler", "DoneHandler")
	generateServerOptionHandler(g, "ErrorHandler", "ErrHandler")

	// Generate NewServer function
	g.P("func New", srvName, "(nc *", natsConn, ", impl ", srvName, ", opts ...ServerOption) {")
	g.P("config := ", microConfig, "{")

	if serviceOptions := service.Desc.Options().(*descriptorpb.ServiceOptions); serviceOptions != nil {
		natsOptions := proto.GetExtension(serviceOptions, pb.E_Nats).(*pb.NATSServiceOptions)
		if natsOptions.GetName() == "" || natsOptions.GetVersion() == "" {
			panic("Name and Version are required for NATSServiceOptions")
		}
		g.P("Name: ", strconv.Quote(natsOptions.GetName())+",")
		g.P("Version: ", strconv.Quote(natsOptions.GetVersion())+",")
		if natsOptions.GetDescription() != "" {
			g.P("Description: ", strconv.Quote(natsOptions.GetDescription())+",")
		}
	} else {
		g.P("Name: ", strconv.Quote(service.GoName), ",")
		g.P("Version: ", strconv.Quote("1.0.0-DEV"), ",") // TODO: make this use actual option nats_in
		g.P("Description: ", strconv.Quote("Using default version and auto derived name - please manually set details using NATSServiceOptions service in protobuf file."), ",")
	}

	g.P("}")
	g.P()
	g.P("for _, opt := range opts {")
	g.P("opt(&config)")
	g.P("}")
	g.P()
	g.P("service, err := micro.AddService(nc, config)")
	g.P("if err != nil {")
	g.P("panic(err) // TODO: Update this to proper error handling")
	g.P("}")
	g.P()

	// Generate service endpoints
	for _, method := range service.Methods {
		g.P("service.AddEndpoint(", strconv.Quote(subjectName(service, method)), ", ", g.QualifiedGoIdent(microPkg.Ident("HandlerFunc")), "(func(request ", g.QualifiedGoIdent(microRequest), ") {")
		g.P("var req ", g.QualifiedGoIdent(method.Input.GoIdent))
		g.P("if err := ", g.QualifiedGoIdent(protoUnmarshal), "(request.Data(), &req); err != nil {")
		g.P("request.Error(", strconv.Quote("560"), ", ", strconv.Quote("Failed to unmarshal proto message"), ", []byte(err.Error()))")
		g.P("return")
		g.P("}")
		g.P()
		g.P("response, err := impl.", method.GoName, "(&req)")
		g.P("if err != nil {")
		g.P("if natsErr, ok := err.(*", g.QualifiedGoIdent(protocPkg.Ident("NATSError")), "); ok {")
		g.P("natsErr.RespondWith(request)")
		g.P("} else {")
		g.P("request.Error(", strconv.Quote("500"), ", ", strconv.Quote("Internal server error"), ", []byte(err.Error()))")
		g.P("}")
		g.P("return")
		g.P("}")
		g.P()
		g.P("data, err := ", protoMarshal, "(response)")
		g.P("if err != nil {")
		g.P("request.Error(", strconv.Quote("560"), ", ", strconv.Quote("Failed to marshal proto message"), ", []byte(err.Error()))")
		g.P("return")
		g.P("}")
		g.P()
		g.P("request.Respond(data)")
		g.P("}))")
		g.P()
	}

	g.P("}")
	g.P()
}

func generateServerOptionHandler(g *protogen.GeneratedFile, field, typ string) {
	g.P("func With", field, "(handler ", g.QualifiedGoIdent(microPkg.Ident(typ)), ") ServerOption {")
	g.P("return func(config *", g.QualifiedGoIdent(microConfig), ") {")
	g.P("config.", field, " = handler")
	g.P("}")
	g.P("}")
	g.P()
}

func generateClient(g *protogen.GeneratedFile, service *protogen.Service) {
	cliName := service.GoName + "NATSClient"

	// Generate client interface
	g.AnnotateSymbol(cliName, protogen.Annotation{Location: service.Location}) // TODO: Find out when to annotate symbols
	g.P("type ", cliName, " interface {")
	for _, method := range service.Methods {
		g.AnnotateSymbol(cliName+"."+method.GoName, protogen.Annotation{Location: method.Location})
		g.P(method.Comments.Leading, methodSignature(g, method))
	}
	g.P("SetTimeout(", timeDuration, ")")
	g.P("}")
	g.P()

	// Generate client struct implementation
	g.P("type ", unexport(cliName), " struct {")
	g.P("nc *", natsConn)
	g.P("timeout ", timeDuration)
	g.P("}")
	g.P()

	// Client struct functions

	//		Generate SetTimeout function
	g.P("func (c *", unexport(cliName), ") SetTimeout(timeout ", timeDuration, ") {")
	g.P("c.timeout = timeout")
	g.P("}")
	g.P()

	// 		Generate handle function
	g.P("func (c *", unexport(cliName), ") handle(req ", protoMessage, ", subject string, out ", protoMessage, ") error {")
	g.P("data, err := ", protoMarshal, "(req)")
	g.P("if err != nil {")
	g.P("return ", g.QualifiedGoIdent(protocErrMarshalling))
	g.P("}")
	g.P("msg, err := c.nc.Request(subject, data, c.timeout)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P("if err := ", g.QualifiedGoIdent(protoUnmarshal), "(msg.Data, out); err != nil {")
	g.P("return ", g.QualifiedGoIdent(protocErrUnmarshalling))
	g.P("}")
	g.P("return nil")
	g.P("}")
	g.P()

	// Generate NewClient function
	g.P("func New", cliName, "(nc *", natsConn, ") ", cliName, " {")
	g.P("return &", unexport(cliName), "{nc: nc, timeout: ", timePkg.Ident("Second"), " * 10}")
	g.P("}")
	g.P()

	// Generate client methods
	for _, method := range service.Methods {
		g.P("func (c *", unexport(cliName), ") ", method.GoName, "(req *", g.QualifiedGoIdent(method.Input.GoIdent), ") (*", g.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
		g.P("var response ", g.QualifiedGoIdent(method.Output.GoIdent))
		g.P("err := c.handle(req, ", strconv.Quote(subjectName(service, method)), ", &response)")
		g.P("return &response, err")
		g.P("}")
		g.P()
	}
}

func generateService(g *protogen.GeneratedFile, service *protogen.Service) {
	generateClient(g, service)
	generateServer(g, service)
	return
}

func subjectName(service *protogen.Service, method *protogen.Method) string {
	return service.GoName + "." + method.GoName
}

func methodSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	s := method.GoName
	s += "(req *" + g.QualifiedGoIdent(method.Input.GoIdent) + ") "
	s += "(*" + g.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
	return s
}
