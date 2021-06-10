package actions

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
)

// FeedIndex default implementation.
func FeedIndex(c buffalo.Context) error {
	feeds := randomFeeds()
	c.Set("feeds", feeds)

	if isTurboFrame(c.Request(), "feed-frame") {
		return c.Render(http.StatusOK, r.Func("text/html", createTurboPlain("feed/feed.plush.html")))
	}

	return c.Render(http.StatusOK, r.HTML("feed/index.plush.html"))
}

// FeedIndexSlow default implementation.
func FeedIndexSlow(c buffalo.Context) error {
	time.Sleep(3 * time.Second)

	feeds := randomFeeds()
	c.Set("feeds", feeds)

	return c.Render(http.StatusOK, r.HTML("feed/index_slow.plush.html"))
}

// FeedIndexWithSSE feeds with SSE.
func FeedIndexWithSSE(c buffalo.Context) error {
	feeds := randomFeeds()
	c.Set("feeds", feeds)

	return c.Render(http.StatusOK, r.HTML("feed/index_with_sse.plush.html"))
}

// FeedIndexWithWebsocket feeds with WebSocket.
func FeedIndexWithWebsocket(c buffalo.Context) error {
	feeds := randomFeeds()
	c.Set("feeds", feeds)

	return c.Render(http.StatusOK, r.HTML("feed/index_with_websocket.plush.html"))
}

// FeedIndex default implementation.
func FeedFrame(c buffalo.Context) error {
	time.Sleep(3 * time.Second)

	feeds := randomFeeds()
	c.Set("feeds", feeds)

	return c.Render(http.StatusOK, r.Func("text/html", createTurboPlain("feed/feed.plush.html")))
}

var names []string
var actions []string
var tasks []string

func init() {
	names = []string{
		"Chuck",
		"Rob",
		"Chip",
		"Buck",
		"Martha",
		"Fred",
		"Rose",
		"Anne",
	}

	actions = []string{
		"created",
		"deleted",
		"updated",
	}

	tasks = []string{
		"Washing the dishes",
		"Buy 3 tomatoes",
		"Clean kitchen",
		"Call mum",
		"Pay rent",
		"Go Running",
	}

}

func RandomFeed(target string) string {
	bgColors := []string{
		"bg-info",
		"bg-light",
		"bg-secondary",
		"bg-primary",
		"bg-danger",
		"bg-warning",
	}
	bgColor := bgColors[rand.Intn(len(bgColors))]

	feedTemplate := `<turbo-stream action="prepend" target="$$TARGET$$"><template><div class="card"><div class="card-body $$BGCOLOR$$">$$TASK$$</div></div></template></turbo-stream>`
	feed := strings.Replace(feedTemplate, "$$TARGET$$", target, -1)
	feed = strings.Replace(feed, "$$BGCOLOR$$", bgColor, -1)
	return strings.Replace(feed, "$$TASK$$", randomTask(), -1)
}

func randomTask() string {
	return names[rand.Intn(len(names))] + " " + actions[rand.Intn(len(actions))] + " task \"" + tasks[rand.Intn(len(tasks))] + "\""
}

func randomFeeds() []string {
	feeds := make([]string, 10)

	for i := 0; i <= len(feeds)-1; i++ {
		feeds[i] = randomTask()
	}

	return feeds
}
