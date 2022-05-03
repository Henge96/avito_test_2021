package config

import (
	"testing"
)

func TestTakeConfigFromYaml(t *testing.T) {

	type checkRes struct {
		pathToFail  string
		expectedRes *Config
	}

	tableCase := []checkRes{
		{
			"./testdata/testconfig2.yml",
			&Config{
				Server: Server{
					Host: "0.0.0.0",
				},
			},
		},
		{"./testdata/testconfig1.yml",
			&Config{
				Server: Server{
					Host: "0.0.0.0",
					Port: Port{
						Http: "5050",
					},
				},
				Db: Db{
					Driver:   "postgres",
					User:     "postgres",
					Password: "postgres",
					HostDb:   "localhost",
					PortDb:   "5432",
					Dbname:   "postgres",
					Mode:     "disable",
				},
			},
		}, {"./testdata/testconfig3.yml",
			&Config{
				Server: Server{
					Host: "0.0.0.0",
					Port: Port{
						Http: "5050",
					},
				},
				Db: Db{
					Driver:   "postgres",
					User:     "postgres",
					Password: "postgres",
					HostDb:   "localhost",
					PortDb:   "5432",
					Dbname:   "postgres",
					Mode:     "disable",
				},
			},
		},
	}

	for _, val := range tableCase {
		result, err := TakeConfigFromYaml(val.pathToFail)
		if err != nil {
			t.Fatal(err)
		}

		switch {

		case result.Server.Host == "":
			result.Server.Host = "0.0.0.0"

		case result.Server.Port.Http == "":
			result.Server.Port.Http = "8080"

		case result.Db.Driver == "":
			result.Db.Driver = "postgres"

		case result.Db.User == "":
			result.Db.User = "postgres"

		case result.Db.Password == "":
			result.Db.Password = "postgres"

		case result.Db.HostDb == "":
			result.Db.HostDb = "localhost"

		case result.Db.PortDb == "":
			result.Db.PortDb = "5432"

		case result.Db.Dbname == "":
			result.Db.Dbname = "postgres"

		case result.Db.Mode == "":
			result.Db.Mode = "disable"

		}

		if *result != *val.expectedRes {
			t.Errorf("Error. Expected %v, got %v", result, val.expectedRes)
		}

	}

}
