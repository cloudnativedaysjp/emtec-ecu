package model

type Notification interface {
	Current() Talk
	Next() Talk
}

//
// NotificationOnDkTimetable
//

type NotificationOnDkTimetable struct {
	current Talk
	next    Talk
}

var _ Notification = (*NotificationOnDkTimetable)(nil)

func NewNotificationOnDkTimetable(current, next Talk) *NotificationOnDkTimetable {
	return &NotificationOnDkTimetable{current, next}
}

func (m NotificationOnDkTimetable) Current() Talk {
	return m.current
}

func (m NotificationOnDkTimetable) Next() Talk {
	return m.next
}

func (m NotificationOnDkTimetable) TrackId() int32 {
	if m.current.TrackId != 0 {
		return m.current.TrackId
	}
	return m.next.TrackId
}

func (m NotificationOnDkTimetable) TrackName() string {
	if m.current.TrackName != "" {
		return m.current.TrackName
	}
	return m.next.TrackName
}
