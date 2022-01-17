package main

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	connStr = kingpin.Arg(
		"conn", "PostgreSQL connection string in URL format").Required().String()
	schema = kingpin.Flag(
		"schema", "PostgreSQL schema name").Default("public").Short('s').String()
	pkgName         = kingpin.Flag("package", "package name").Default("main").Short('p').String()
	typeMapFilePath = kingpin.Flag("typemap", "column type and go type map file path").Short('t').String()
	exTbls          = kingpin.Flag("exclude", "table names to exclude").Short('x').Strings()
	customTmpl      = kingpin.Flag("template", "custom template path").Default("template/methods.gohtml").String()
	outPath         = kingpin.Flag("output", "output directory path").Short('o').String()
	onlyInterface   = kingpin.Flag("only-interface", "output without Criteria and order types").Short('i').Bool()
)

//go:generate go get golang.org/x/tools/cmd/goimports

func main() {
	kingpin.Parse()

	conn, err := OpenDB(*connStr)
	if err != nil {
		log.Fatal(err)
	}

	oInfo, err := os.Stat(*outPath)
	if !oInfo.IsDir() {
		log.Fatalln("Out path is not a directory")
	}

	builder := NewPgOrmBuilder(conn, *schema, *typeMapFilePath, *pkgName, *exTbls)

	tables, err := builder.GetPgStruct()
	if err != nil {
		log.Fatal(err)
	}

	outPutPath := path.Join(*outPath, *schema)
	err = os.MkdirAll(outPutPath, 0775)
	if err != nil {
		log.Fatal(err)
	}

	src, err := builder.RenderPgCustomTmpl(nil, "template/queryInterface.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	err = writeFile(outPutPath, "queryInterface", src)
	if err != nil {
		log.Fatalf("failed to create output file %s: %s", "queryInterface", err)
	}

	if !*onlyInterface {
		src, err := builder.RenderPgCustomTmpl(nil, "template/queryTypes.gohtml")
		if err != nil {
			log.Fatal(err)
		}

		err = writeFile(outPutPath, "queryTypes", src)
		if err != nil {
			log.Fatalf("failed to create output file %s: %s", "queryTypes", err)
		}
	}

	for _, tbl := range tables {
		src, err := builder.RenderPgCustomTmpl(tbl, *customTmpl)
		if err != nil {
			log.Fatal(err)
		}

		err = writeFile(outPutPath, tbl.Name, src)
		if err != nil {
			log.Fatalf("failed to create output file %s: %s", tbl.Name, err)
		}
	}
}

func writeFile(outPutPath string, tblName string, src []byte) (err error) {

	outFile := path.Join(outPutPath, tblName+".go")

	out, err := os.Create(outFile)
	if err != nil {
		return errors.Wrap(err, "failed to create output file: ")
	}
	if _, err := out.Write(src); err != nil {
		return err
	}
	params := []string{"-w", outFile}
	if err = exec.Command("goimports", params...).Run(); err != nil {
		err = errors.Wrap(err, "failed to goimports: ")
	}
	return err
}
