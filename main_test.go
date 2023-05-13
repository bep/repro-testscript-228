package main

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bep/helpers/envhelpers"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestScripts(t *testing.T) {
	setup := testSetupFunc()
	testscript.Run(t, testscript.Params{
		Dir: "testscripts",
		//..

		// UpdateScripts: true, // Uncomment to rewrite the test scripts with
		// TestWork: true, // Uncomment to keep the test work dir.
		Setup: func(env *testscript.Env) error {
			return setup(env)
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"sleep": func(ts *testscript.TestScript, neg bool, args []string) {
				i := 1
				if len(args) > 0 {
					var err error
					i, err = strconv.Atoi(args[0])
					if err != nil {
						i = 1
					}
				}
				time.Sleep(time.Duration(i) * time.Second)
			},
		},
	})
}

func TestMain(m *testing.M) {
	os.Exit(
		testscript.RunMain(m, map[string]func() int{
			"runserver": func() int {
				runServer(os.Args[1:])
				return 0
			},
		}),
	)
}

func testSetupFunc() func(env *testscript.Env) error {
	sourceDir, _ := os.Getwd()
	return func(env *testscript.Env) error {
		var keyVals []string
		// Add some environment variables to the test script.
		keyVals = append(keyVals, "SOURCE", sourceDir)
		envhelpers.SetEnvVars(&env.Vars, keyVals...)

		return nil
	}
}
