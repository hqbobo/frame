// 自定义protobuf 代码生成器
//
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/micro/protoc-gen-micro/generator"
	_ "github.com/micro/protoc-gen-micro/plugin/micro"
	"io/ioutil"
	"os"
	"strings"
)

const head_format = "// %s \n" +
	"type %s struct {\n" +
	"	svc    proto.%sService\n" +
	"}\n\n"

const new_func = "func new%s(cli service.Client) *%s {\n" +
	"	s := new(%s)\n" +
	"	s.svc = proto.New%sService(define.Svc%s, cli)\n" +
	"	return s\n" +
	"}\n\n"

const func_format = "// %s %s函数\n" +
	"func (cli * %s) %s(ctx context.Context, req *%s)(*%s, error) {\n" +
	"	rsp, err := cli.svc.%s(ctx, req)\n" +
	"	if err != nil {\n" +
	"		log.Error(\"failed \", err)\n" +
	"		return nil, err\n" +
	"	}\n" +
	"	return rsp, err\n" +
	"}\n"

func removeFirstDot(in *string) string {
	out := []byte(*in)[1:]
	return string(out)
}

func capitaiFirst(in *string) string {
	f := string([]byte(*in)[0])
	out := string([]byte(*in)[1:])
	return strings.ToUpper(f) + out
}

func generateClient(g *generator.Generator) {
	var out string
	out += "package cli\n\n"
	env := os.Getenv("PWD")
	in := string([]byte(env)[strings.Index(env, "github.com/hqbobo/frame/"):])
	out += "import (\n"
	out += "	\"" + in + "\"\n"
	out += "	\"context\"\n"
	out += "	\"github.com/hqbobo/frame/common/service\"\n"
	out += "	\"github.com/hqbobo/frame/common/log\"\n"
	out += ")\n\n"
	for _, v := range g.Request.ProtoFile[0].Service {
		name := capitaiFirst(v.Name)
		svc := name + "Service"
		out += fmt.Sprintf(head_format, svc, svc, name)
		out += fmt.Sprintf(new_func, svc, svc, svc, name, name)
		for _, f := range v.Method {
			//fmt.Printf("%+v\n", *f)
			out += fmt.Sprintf(func_format, *f.Name, *f.Name, svc, *f.Name,
				removeFirstDot(f.InputType), removeFirstDot(f.OutputType),
				*f.Name)
			out += "\n"
		}
	}
	os.Mkdir("cli", os.ModePerm)
	ioutil.WriteFile("cli/cli.go", []byte(out), os.ModePerm)
}

const shead_format = "type %s struct {\n" +
	"}\n\n"

const sfunc_format = "// %s %s函数\n" +
	"func (%s * %s) %s(ctx context.Context, req *%s , rsp *%s) error {\n" +
	"	return nil\n" +
	"}\n"

func generateSvc(g *generator.Generator) {
	var out string
	out += "package svc\n\n"
	env := os.Getenv("PWD")
	in := string([]byte(env)[strings.Index(env, "github.com/hqbobo/frame/"):])
	out += "import (\n"
	out += "	\"" + in + "\"\n"
	out += "	\"context\"\n"
	out += "	\"github.com/hqbobo/frame/common/service\"\n"
	out += "	\"github.com/hqbobo/frame/common/log\"\n"
	out += ")\n\n"
	for _, v := range g.Request.ProtoFile[0].Service {
		name := capitaiFirst(v.Name)
		first := strings.ToLower(string([]byte(*v.Name)[0]))
		svc := name + "Service"
		out += fmt.Sprintf(shead_format, svc)
		for _, f := range v.Method {
			out += fmt.Sprintf(sfunc_format, *f.Name, *f.Name, first, svc, *f.Name,
				removeFirstDot(f.InputType), removeFirstDot(f.OutputType))
			out += "\n"
		}
	}
	os.Mkdir("server", os.ModePerm)
	ioutil.WriteFile("server/impl.go", []byte(out), os.ModePerm)
}

func main() {
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}
	//fmt.Println(*g.Request.ProtoFile[0])
	generateClient(g)
	generateSvc(g)
	//g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	//g.WrapTypes()

	//g.SetPackageNames()
	//g.BuildTypeNameMap()

	//g.GenerateAllFiles()

	// Send back the results.
	//data, err = proto.Marshal(g.Response)
	//if err != nil {
	//	g.Error(err, "failed to marshal output proto")
	//}
	//_, err = os.Stdout.Write(data)
	//if err != nil {
	//	g.Error(err, "failed to write output proto")
	//}
}
