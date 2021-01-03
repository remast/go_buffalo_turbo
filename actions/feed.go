package actions

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
)

// FeedIndex default implementation.
func FeedIndex(c buffalo.Context) error {
	time.Sleep(1 * time.Second)

	names := []string{
		"Chuck",
		"Rob",
		"Chip",
		"Buck",
		"Martha",
		"Fred",
		"Rose",
		"Anne",
	}

	actions := []string{
		"created",
		"deleted",
		"updated",
	}

	tasks := []string{
		"Washing the dishes",
		"Buy 3 tomatoes",
		"Clean kitchen",
		"Call mum",
		"Pay rent",
		"Go Running",
	}

	feeds := [10]string{}

	for i := 0; i <= len(feeds)-1; i++ {
		feeds[i] = names[rand.Intn(len(names))] + " " + actions[rand.Intn(len(actions))] + " task \"" + tasks[rand.Intn(len(tasks))] + "\""
	}

	c.Set("feeds", feeds)
	return c.Render(http.StatusOK, r.HTML("feed/index.plush.html"))
}
