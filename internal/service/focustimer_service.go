package service

import (
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/sriram15/progressor-todo-app/internal/events"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

// IFocusTimerService defines the interface for the focus timer service.
type IFocusTimerService interface {
	RegisterEventHandlers()
	Shutdown()
	ResumeTimer()
	StopAndDeactivate()
}

// FocusTimerService manages the focus timer for active cards.
type FocusTimerService struct {
	app 		  *application.App
	eventBus       *events.EventBus
	cardService    ICardService
	settingService ISettingService

	timer        *time.Timer
	startTime    time.Time
	activeCardID int64
	mu           sync.Mutex
}

// NewFocusTimerService creates a new FocusTimerService.
func NewFocusTimerService(cs ICardService, ss ISettingService, bus *events.EventBus, app *application.App) *FocusTimerService {
	return &FocusTimerService{
		cardService:    cs,
		settingService: ss,
		eventBus:       bus,
		app: 		  app,
	}
}

// validateNotification checks and requests notification authorization.
func (s *FocusTimerService) validateNotification() bool {
	if runtime.GOOS == "darwin" {
		log.Println("Notifications are disabled on macOS.")
		return false
	}

	notifier := notifications.New()
	authorized, err := notifier.CheckNotificationAuthorization()
	if err != nil {
		log.Printf("Error checking notification authorization: %v", err)
		return false
	}

	if !authorized {
		log.Println("Notification authorization not granted. Requesting...")
		authorized, err = notifier.RequestNotificationAuthorization()
		if err != nil {
			log.Printf("Error requesting notification authorization: %v", err)
			return false
		}
		if !authorized {
			log.Println("Notification authorization denied by user.")
			return false
		}
		log.Println("Notification authorization granted.")
	}
	return true
}

// RegisterEventHandlers subscribes the service to necessary events.
func (s *FocusTimerService) RegisterEventHandlers() {
	s.eventBus.Subscribe(events.CardStartedTopic, s.handleCardStarted)
	s.eventBus.Subscribe(events.CardStoppedTopic, s.handleCardStopped)
}

// handleCardStarted is the event handler for when a card is started.
func (s *FocusTimerService) handleCardStarted(eventData interface{}) {
	event, ok := eventData.(events.CardStartedEvent)
	if !ok {
		log.Printf("Error: received non-CardStartedEvent for topic %s", events.CardStartedTopic)
		return
	}
	log.Printf("Received CardStartedEvent: %+v", event)

	if s.validateNotification() {
		notification := notifications.New()
		notificationOption := notifications.NotificationOptions{
			ID: "focus-timer-started",
			Title: "Focus Timer Started",
		}

		err := notification.SendNotification(notificationOption);
		if err != nil {
			log.Println("Error sending notification:", err)
		}
	} else {
		log.Println("Notification not authorized, skipping send.")
	}

	s.startTimer(event.CardID)
}

// handleCardStopped is the event handler for when a card is stopped.
func (s *FocusTimerService) handleCardStopped(eventData interface{}) {
	event, ok := eventData.(events.CardStoppedEvent)
	if !ok {
		log.Printf("Error: received non-CardStoppedEvent for topic %s", events.CardStoppedTopic)
		return
	}
	log.Printf("Received CardStoppedEvent: %+v", event)
	s.stopTimer()
}

// startTimer starts the focus timer for a given card.
func (s *FocusTimerService) startTimer(cardID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.timer != nil {
		s.timer.Stop()
	}

	s.activeCardID = cardID
	s.startTime = time.Now()

	var activeCardTimeout int = 30
	activeCardSettingStr, err := s.settingService.GetSetting("active_card_timeout")
	if err != nil {
		log.Printf("Error getting active card timeout setting: %v", err)
	} else {
		activeCardTimeout, _ = strconv.Atoi(activeCardSettingStr)
	}

	focusDuration := time.Duration(activeCardTimeout) * time.Minute
	s.timer = time.AfterFunc(focusDuration, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		log.Println("Focus timer completed for card:", s.activeCardID)

		if s.validateNotification() {
			notificationCategoryId := "active-card-timer-complete"
			category := notifications.NotificationCategory{
				ID: notificationCategoryId,
				Actions: []notifications.NotificationAction{
					{
            ID:    "CONTINUE",
            Title: "Continue",
        },{
					ID:    "STOP",
					Title: "Stop",
				},
				},
			}

			notification := notifications.New()
			notification.RegisterNotificationCategory(category)

			notification.OnNotificationResponse(func(result notifications.NotificationResult) {
   	 			response := result.Response
				if response.ActionIdentifier == "STOP" {
					s.StopAndDeactivate()
				}

			})

			notificationOption := notifications.NotificationOptions{
				ID:    "focus-timer-completed",
				Title: "Focus Timer Completed",
				CategoryID: notificationCategoryId,
			}

			err := notification.SendNotificationWithActions(notificationOption)
			if err != nil {
				log.Println("Error sending notification:", err)
			}
		} else {
			log.Println("Notification not authorized, skipping send.")
		}

		s.app.Event.Emit("active_card_timer_complete", s.activeCardID)
	})

	log.Println("Focus timer started for card:", cardID)
}

// stopTimer stops the focus timer.
func (s *FocusTimerService) stopTimer() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
		log.Println("Focus timer stopped for card:", s.activeCardID)
		s.activeCardID = 0
	}
}

// Shutdown is called on application exit to gracefully handle any active timer.
func (s *FocusTimerService) Shutdown() {
	s.stopTimer()
}

// ResumeTimer resumes the timer after a user confirmation.
func (s *FocusTimerService) ResumeTimer() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resumeTimerInternal()
}

func (s *FocusTimerService) resumeTimerInternal() {
	cardID := s.activeCardID
	if cardID != 0 {
		log.Println("Resuming timer for card:", cardID)
		// Drop the lock to avoid deadlock with startTimer
		s.mu.Unlock()
		s.startTimer(cardID)
		s.mu.Lock()
	}
}

// StopAndDeactivate stops the timer and deactivates the card.
func (s *FocusTimerService) StopAndDeactivate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopAndDeactivateInternal()
}

func (s *FocusTimerService) stopAndDeactivateInternal() {
	cardID := s.activeCardID
	if cardID != 0 {
		// Drop the lock to avoid deadlock with stopTimer
		s.mu.Unlock()
		s.stopTimer()
		s.mu.Lock()
		// Assuming project ID is 1, this might need to be retrieved differently
		if err := s.cardService.StopCard(1, uint(cardID)); err != nil {
			log.Printf("Error deactivating card %d: %v", cardID, err)
		}
	}
}
