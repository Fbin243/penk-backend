package cron

import (
	"context"
	"log"

	cron "tenkhours/pkg/cron"
	"tenkhours/services/notification/business"
	messagequeue "tenkhours/services/notification/message-queue"
)

type ReminderCron struct {
	ReminderRepo      business.IReminderRepo
	NotificationQueue *messagequeue.NotificationQueue
}

func NewReminderCron(reminderRepo business.IReminderRepo, notificationQueue *messagequeue.NotificationQueue) *ReminderCron {
	return &ReminderCron{
		ReminderRepo:      reminderRepo,
		NotificationQueue: notificationQueue,
	}
}

func (r *ReminderCron) Start() {
	c := cron.NewCron()
	c.RunEveryHours(r.ProcessReminders)
	c.Start()
}

// ProcessReminders fetches reminders for the next hour and publishes them to Kafka.
func (r *ReminderCron) ProcessReminders() {
	ctx := context.Background()

	// Fetch reminders that are due within the next hour
	reminders, err := r.ReminderRepo.GetUpcomingReminders(ctx)
	if err != nil {
		log.Printf("Error fetching reminders: %v\n", err)
		return
	}

	// If no reminders are found, log and return
	if len(reminders) == 0 {
		log.Println("No reminders found for the next hour.")
		return
	}

	// Publish each reminder to Kafka
	for _, reminder := range reminders {
		err := r.NotificationQueue.PublishReminder(ctx, reminder)
		if err != nil {
			log.Printf("Failed to publish reminder %s to Kafka: %v\n", reminder.ID, err)
		} else {
			log.Printf("Reminder %s published to Kafka\n", reminder.ID)
		}
	}
}
