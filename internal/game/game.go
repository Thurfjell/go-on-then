package game

import (
	"embed"
	"fmt"
	"gogopowerrangers/internal/gogopowerrangers"
	"gogopowerrangers/internal/webroot"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

const eventMood string = "mood"
const eventClick string = "click"
const eventLike string = "like"

type State interface {
	Get() *gogopowerrangers.State
	Like()
	Click()
}

type event struct {
	EventType string
	Value     string
}

type Service struct {
	template        *template.Template
	state           State
	eventChannelMap map[string]chan (event)
	eventLock       sync.RWMutex
}

func (service *Service) page() webroot.RouteMeta {
	return webroot.RouteMeta{
		Path: "GET /game",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			state := service.state.Get()
			err := service.template.ExecuteTemplate(w, "page.html", state)

			if err != nil {
				fmt.Fprint(w, "oh no.. :(")
			}
		}),
	}
}

func (service *Service) gameEvents() webroot.RouteMeta {
	return webroot.RouteMeta{
		Path: "GET /game-events",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := uuid.NewString()
			service.eventLock.Lock()
			service.eventChannelMap[id] = make(chan event)
			service.eventLock.Unlock()

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			f := w.(http.Flusher)

		L:
			for {
				select {
				case <-r.Context().Done():
					service.eventLock.Lock()
					close(service.eventChannelMap[id])
					delete(service.eventChannelMap, id)
					service.eventLock.Unlock()
					f.Flush()
					break L
				case event, ok := <-service.eventChannelMap[id]:
					if ok {
						fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.EventType, event.Value)
						f.Flush()
					}
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}

		}),
	}
}

func (service *Service) clickEvent() webroot.RouteMeta {
	return webroot.RouteMeta{
		Path: "POST /click-event",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := service.state.Get()
			oldMood := state.Mood
			service.state.Click()
			shouldMoodEvent := oldMood != state.Mood

			go func() {
				for a := range service.eventChannelMap {
					service.eventChannelMap[a] <- event{
						EventType: eventClick,
						Value:     strconv.Itoa(state.Clicks),
					}

					if shouldMoodEvent {
						service.eventChannelMap[a] <- event{
							EventType: eventMood,
							Value:     state.Mood,
						}
					}
				}
			}()
		}),
	}
}

func (service *Service) likeEvent() webroot.RouteMeta {
	return webroot.RouteMeta{
		Path: "POST /like-event",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := service.state.Get()
			oldMood := state.Mood
			service.state.Like()
			shouldMoodEvent := oldMood != state.Mood

			go func() {
				for a := range service.eventChannelMap {
					service.eventChannelMap[a] <- event{
						EventType: eventLike,
						Value:     strconv.Itoa(state.Likes),
					}
					if shouldMoodEvent {
						service.eventChannelMap[a] <- event{
							EventType: eventMood,
							Value:     state.Mood,
						}
					}
				}
			}()
		}),
	}
}

func (service *Service) Routes() []webroot.RouteMeta {
	return []webroot.RouteMeta{
		service.page(),
		service.gameEvents(),
		service.clickEvent(),
		service.likeEvent(),
	}
}

//go:embed templates/*
var content embed.FS

func NewService(state State) (*Service, error) {
	template, err := template.ParseFS(content, "templates/*.html")

	if err != nil {
		return nil, err
	}

	return &Service{
		template:        template,
		state:           state,
		eventChannelMap: make(map[string]chan event),
		eventLock:       sync.RWMutex{},
	}, err
}
