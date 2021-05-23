package actions

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/validate/v3"
)

var taskItems []*TaskItem

func init() {
	taskItems = []*TaskItem{}
}

// TaskItem for a single task
type TaskItem struct {
	ID   int
	Name string
	Done bool
}

// TaskIndex default implementation.
func TaskIndex(c buffalo.Context) error {
	openTaskItems := make([]*TaskItem, 0)
	for _, taskItem := range taskItems {
		if !taskItem.Done {
			openTaskItems = append(openTaskItems, taskItem)
		}
	}
	c.Set("taskItem", &TaskItem{})
	c.Set("taskItems", openTaskItems)
	return c.Render(http.StatusOK, r.HTML("task/index.plush.html"))
}

// TaskCompleted default implementation.
func TaskCompleted(c buffalo.Context) error {
	completedTaskItems := make([]*TaskItem, 0)
	for _, taskItem := range taskItems {
		if taskItem.Done {
			completedTaskItems = append(completedTaskItems, taskItem)
		}
	}
	c.Set("taskItems", completedTaskItems)
	return c.Render(http.StatusOK, r.HTML("task/completed.plush.html"))
}

// TaskToggle default implementation.
func TaskCheck(c buffalo.Context) error {
	taskItemID := c.Params().Get("ID")

	for _, taskItem := range taskItems {
		if strconv.Itoa(taskItem.ID) == taskItemID {
			taskItem.Done = !taskItem.Done
			c.Set("taskItem", taskItem)
			break
		}
	}

	if acceptsTurboStream(c.Request()) {
		id := "task_item_" + taskItemID
		return c.Render(http.StatusOK, r.Func("text/vnd.turbo-stream.html", createTurboWriter("task/item.plush.html", "remove", id)))
	}

	return c.Redirect(302, "/")
}

// TaskCreate default implementation.
func TaskCreate(c buffalo.Context) error {
	newTaskItem := &TaskItem{ID: rand.Intn(1000)}
	if err := c.Bind(newTaskItem); err != nil {
		return err
	}
	// Handle form errors
	if newTaskItem.Name == "" {
		c.Set("taskItem", newTaskItem)

		verrs := validate.NewErrors()
		verrs.Add("name", "Name missing.")
		c.Set("errors", verrs)

		if acceptsTurboStream(c.Request()) {
			turboAction := "replace"
			turboDomID := "task_new_form"
			return c.Render(http.StatusOK, r.Func("text/vnd.turbo-stream.html", createTurboWriter("task/new.plush.html", turboAction, turboDomID)))
		}

		return c.Render(http.StatusOK, r.HTML("task/new.plush.html"))
	}

	taskItems = append(taskItems, newTaskItem)
	return c.Redirect(302, "/")
}

// TaskNew default implementation.
func TaskNew(c buffalo.Context) error {
	c.Set("taskItem", &TaskItem{})
	return c.Render(http.StatusOK, r.HTML("task/new.plush.html"))
}

func acceptsTurboStream(request *http.Request) bool {
	for _, acceptValue := range request.Header["Accept"] {
		if strings.Contains(acceptValue, "text/vnd.turbo-stream.html") {
			return true
		}
	}
	return false
}

func createTurboWriter(template, action, target string) render.RendererFunc {
	return func(w io.Writer, d render.Data) error {
		d["action"] = action
		d["target"] = target
		r.HTML(template, "turbo/turbo_stream.plush.html").Render(w, d)
		return nil
	}
}
