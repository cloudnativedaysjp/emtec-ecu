package model

type Notification interface {
	Current() Talk
	Next() Talk
}

func getTrackId(current, next Talk) int32 {
	if current.TrackId != 0 {
		return current.TrackId
	}
	return next.TrackId
}

func getTrackName(current, next Talk) string {
	if current.TrackName != "" {
		return current.TrackName
	}
	return next.TrackName
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
	return getTrackId(m.current, m.next)
}

func (m NotificationOnDkTimetable) TrackName() string {
	return getTrackName(m.current, m.next)
}

//
// NotificationSceneMovedToNext
//

type NotificationSceneMovedToNext struct {
	current Talk
	next    Talk
}

var _ Notification = (*NotificationSceneMovedToNext)(nil)

func NewNotificationSceneMovedToNext(current, next Talk) *NotificationSceneMovedToNext {
	return &NotificationSceneMovedToNext{current, next}
}

func (m NotificationSceneMovedToNext) Current() Talk {
	return m.current
}

func (m NotificationSceneMovedToNext) Next() Talk {
	return m.next
}

func (m NotificationSceneMovedToNext) TrackId() int32 {
	return getTrackId(m.current, m.next)
}

func (m NotificationSceneMovedToNext) TrackName() string {
	return getTrackName(m.current, m.next)
}
