package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/segmentio/ksuid"
)

type requestBody struct {
	RequestVersion string `json:"request_version,omitempty"`
	RequestOpts    string `json:"request_opts,omitempty"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Only POST allowed")
			return
		}

		workspace, cleanup, err := setupNewWorkspace()
		defer cleanup()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to create workspace: %s", err.Error())
			return
		}

		var b requestBody
		err = json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Error decoding body: %s", err.Error())
			return
		} else if b.RequestVersion == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing request_version")
			return
		} else if b.RequestOpts == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing request_opts")
			return
		}

		if _, err := execCmdInDir(workspace, "npm", "init", "-y"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error init npm: %s", err.Error())
			return
		}
		if _, err := execCmdInDir(workspace, "npm", "install", fmt.Sprintf("request@%s", b.RequestVersion)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error init npm: %s", err.Error())
			return
		}
		output, err := execCmdInDir(workspace, "node", "index.js", b.RequestOpts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error executing request: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(output))
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func execCmdInDir(dir, prog string, args ...string) (string, error) {
	cmd := exec.Command(prog, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func setupNewWorkspace() (path string, cleanup func(), err error) {
	path = ksuid.New().String()
	cleanup = func() {}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return
	}
	cleanup = func() {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Printf("Failed to clean up workspace %s", path)
		}
	}

	_, err = exec.Command("cp", "index.js", path).Output()
	return
}
