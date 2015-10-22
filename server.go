/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type task struct {
	Id int `json:"id"`
	Content string `json:"content"`
	Other   string `json:"other"`
}

const dataFile = "./tasks.json"

var taskMutex = new(sync.Mutex)

// Handle tasks
func handleTasks(w http.ResponseWriter, r *http.Request) {
	// Since multiple requests could come in at once, ensure we have a lock
	// around all file operations
	taskMutex.Lock()
	defer taskMutex.Unlock()

	// Stat the file, so we can find its current permissions
	fi, err := os.Stat(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	// Read the comments from the file.
	taskData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		// Decode the JSON data
		var tasks []task
		if err := json.Unmarshal(taskData, &tasks); err != nil {
			http.Error(w, fmt.Sprintf("Unable to Unmarshal tasks from data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		// Add a new task to the in memory slice of tasks
		tasks = append(tasks, task{Id: len(tasks) , Content: r.FormValue("content"), Other: r.FormValue("other")})

		// Marshal the tasks to indented json.
		taskData, err = json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to marshal tasks to json: %s", err), http.StatusInternalServerError)
			return
		}

		// Write out the tasks to the file, preserving permissions
		err := ioutil.WriteFile(dataFile, taskData, fi.Mode())
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to write tasks to data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		io.Copy(w, bytes.NewReader(taskData))

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		// stream the contents of the file to the response
		io.Copy(w, bytes.NewReader(taskData))

	default:
		// Don't know the method, so error
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.HandleFunc("/api/tasks", handleTasks)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
