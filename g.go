// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import "os"

var cmdGenerate = &Command{
	UsageLine: "generate [Command]",
	Short:     "generate code based on application",
	Long: `
bee generate scaffold [scaffoldname] [-fields=""]	
	The generate scaffold command will do a number of things for you。
	-fields: a list of database fields.
	example: bee generate scaffold post -fields="title:string,body:text"

bee generate model [modelname] [-fields=""]
	generate RESTFul model based on fields

bee generate controller [controllerfile]
    generate RESTFul controllers             

bee generate view [viewpath]
    generate CRUD view in viewpath

bee generate migration [migrationfile]
    generate migration file for making database schema update
	
bee generate docs
    generate swagger doc file

bee generate test [routerfile]
    generate testcase

bee generate appcode [-tables=""] [-driver=mysql] [-conn=root:@tcp(127.0.0.1:3306)/test] [-level=1]
    generate appcode based on an existing database
    -tables: a list of table names separated by ',', default is empty, indicating all tables
    -driver: [mysql | postgresql | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test
    -level:  [1 | 2 | 3], 1 = model; 2 = models,controller; 3 = models,controllers,router		
`,
}

var driver docValue
var conn docValue
var level docValue
var tables docValue
var fields docValue

func init() {
	cmdGenerate.Run = generateCode
	cmdGenerate.Flag.Var(&tables, "tables", "specify tables to generate model")
	cmdGenerate.Flag.Var(&driver, "driver", "database driver: mysql, postgresql, etc.")
	cmdGenerate.Flag.Var(&conn, "conn", "connection string used by the driver to connect to a database instance")
	cmdGenerate.Flag.Var(&level, "level", "1 = models only; 2 = models and controllers; 3 = models, controllers and routers")
	cmdGenerate.Flag.Var(&fields, "fields", "specify the fields want to generate.")
}

func generateCode(cmd *Command, args []string) {
	curpath, _ := os.Getwd()
	if len(args) < 1 {
		ColorLog("[ERRO] command is missing\n")
		os.Exit(2)
	}

	gopath := os.Getenv("GOPATH")
	Debugf("gopath:%s", gopath)
	if gopath == "" {
		ColorLog("[ERRO] $GOPATH not found\n")
		ColorLog("[HINT] Set $GOPATH in your environment vairables\n")
		os.Exit(2)
	}

	gcmd := args[0]
	switch gcmd {
	case "scaffold":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate scaffold [scaffoldname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate scaffold [scaffoldname] [-fields=\"title:string,body:text\"]\n")
			os.Exit(2)
		}
		sname := args[1]
		ColorLog("[INFO] Using '%s' as scaffold name\n", sname)
		generateScaffold(sname, fields.String(), curpath)
	case "docs":
		generateDocs(curpath)
	case "appcode":
		cmd.Flag.Parse(args[1:])
		if driver == "" {
			driver = "mysql"
		}
		if conn == "" {
			conn = "root:@tcp(127.0.0.1:3306)/test"
		}
		if level == "" {
			level = "1"
		}
		ColorLog("[INFO] Using '%s' as 'driver'\n", driver)
		ColorLog("[INFO] Using '%s' as 'conn'\n", conn)
		ColorLog("[INFO] Using '%s' as 'tables'", tables)
		ColorLog("[INFO] Using '%s' as 'level'\n", level)
		generateAppcode(string(driver), string(conn), string(level), string(tables), curpath)
	case "migration":
		if len(args) == 2 {
			mname := args[1]
			ColorLog("[INFO] Using '%s' as migration name\n", mname)
			generateMigration(mname, curpath)
		} else {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate migration [migrationname]\n")
			os.Exit(2)
		}
	case "controller":
		if len(args) == 2 {
			cname := args[1]
			generateController(cname, curpath)
		} else {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate controller [controllername]\n")
			os.Exit(2)
		}
	case "model":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate model [modelname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate model [modelname] [-fields=\"title:string,body:text\"]\n")
			os.Exit(2)
		}
		sname := args[1]
		ColorLog("[INFO] Using '%s' as model name\n", sname)
		generateModel(sname, fields.String(), curpath)
	case "view":
		if len(args) == 2 {
			cname := args[1]
			generateView(cname, curpath)
		} else {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate view [viewpath]\n")
			os.Exit(2)
		}
	default:
		ColorLog("[ERRO] command is missing\n")
	}
	ColorLog("[SUCC] generate successfully created!\n")
}