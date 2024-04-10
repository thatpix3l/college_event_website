package gen_sql

import (
	s "github.com/go-jet/jet/v2/postgres"
	m "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/model"
	t "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/table"
)

func CreateTag() s.InsertStatement {
	return t.Tag.INSERT(t.Tag.MutableColumns)
}

func CreateTaggedEvent() s.InsertStatement {
	return t.Taggedevent.INSERT(t.Taggedevent.AllColumns)
}

func CreateEvent() s.InsertStatement {
	return t.Baseevent.INSERT(t.Baseevent.MutableColumns)
}

func ReadEvents() s.SelectStatement {
	return s.SELECT(
		t.Baseevent.AllColumns,
		t.Rsoevent.AllColumns,
		t.Publicevent.AllColumns,
		t.Privateevent.AllColumns,
		t.Tag.AllColumns,
		t.Comment.AllColumns,
		t.Rating.AllColumns,
	).FROM(
		t.Baseevent.LEFT_JOIN(
			t.Rsoevent, t.Baseevent.ID.EQ(t.Rsoevent.ID),
		).LEFT_JOIN(
			t.Publicevent, t.Baseevent.ID.EQ(t.Publicevent.ID),
		).LEFT_JOIN(
			t.Privateevent, t.Baseevent.ID.EQ(t.Privateevent.ID),
		).LEFT_JOIN(
			t.University, t.Baseevent.UniversityID.EQ(t.University.ID),
		).LEFT_JOIN(
			t.Taggedevent, t.Baseevent.ID.EQ(t.Taggedevent.BaseEventID),
		).LEFT_JOIN(
			t.Tag, t.Tag.ID.EQ(t.Taggedevent.TagID).AND(t.Taggedevent.BaseEventID.EQ(t.Baseevent.ID)),
		).LEFT_JOIN(
			t.Comment, t.Comment.BaseEventID.EQ(t.Baseevent.ID),
		).LEFT_JOIN(
			t.Rating, t.Rating.BaseEventID.EQ(t.Baseevent.ID),
		),
	)
}

func CreatePublicEvent() s.InsertStatement {
	return t.Publicevent.INSERT(t.Publicevent.AllColumns)
}

type Event struct {
	m.Baseevent
	*m.Rsoevent
	*m.Publicevent
	*m.Privateevent
	*m.University
	Tags     []m.Tag
	Ratings  []m.Rating
	Comments []m.Comment
}

func CreateBaseUser() s.InsertStatement { return t.Baseuser.INSERT(t.Baseuser.MutableColumns) }

func ReadUsers() s.SelectStatement {
	return s.SELECT(
		t.Baseuser.AllColumns,
		t.Student.AllColumns,
		t.Superadmin.AllColumns,
		t.Rsomember.AllColumns,
	).FROM(
		t.Baseuser.LEFT_JOIN(
			t.Student, t.Baseuser.ID.EQ(t.Student.ID),
		).LEFT_JOIN(
			t.Superadmin, t.Baseuser.ID.EQ(t.Superadmin.ID),
		).LEFT_JOIN(
			t.Rsomember, t.Baseuser.ID.EQ(t.Rsomember.ID),
		),
	)
}

type User struct {
	m.Baseuser
	*m.Student
	*m.Superadmin
	*m.Rsomember
}

func CreateUniversity() s.InsertStatement { return t.University.INSERT(t.University.MutableColumns) }

func ReadUniversities() s.SelectStatement {
	return s.SELECT(
		t.University.AllColumns,
	).FROM(
		t.University,
	)
}

func CreateStudent() s.InsertStatement { return t.Student.INSERT(t.Student.ID) }

func ReadRsos() s.SelectStatement { return t.Rso.SELECT(t.Rso.AllColumns).FROM(t.Rso) }

func CreateComment() s.InsertStatement {
	return t.Comment.INSERT(t.Comment.MutableColumns)
}
