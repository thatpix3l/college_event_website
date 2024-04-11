package gen_sql

import (
	pg "github.com/go-jet/jet/v2/postgres"
	m "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/model"
	t "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/table"
)

func CreateTag() pg.InsertStatement {
	return t.Tag.INSERT(t.Tag.MutableColumns)
}

func CreateTaggedEvent() pg.InsertStatement {
	return t.Taggedevent.INSERT(t.Taggedevent.AllColumns)
}

func CreateEvent() pg.InsertStatement {
	return t.Baseevent.INSERT(t.Baseevent.MutableColumns)
}

func ReadEvents() pg.SelectStatement {
	return pg.SELECT(
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

func CreatePublicEvent() pg.InsertStatement {
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

func CreateBaseUser() pg.InsertStatement { return t.Baseuser.INSERT(t.Baseuser.MutableColumns) }

func FullUserColumns() pg.ProjectionList {
	return []pg.Projection{
		t.Baseuser.AllColumns,
		t.Student.AllColumns,
		t.Superadmin.AllColumns,
		t.Rsomember.AllColumns,
	}
}

func FullUserTable() pg.ReadableTable {
	return t.Baseuser.LEFT_JOIN(
		t.Student, t.Baseuser.ID.EQ(t.Student.ID),
	).LEFT_JOIN(
		t.Superadmin, t.Baseuser.ID.EQ(t.Superadmin.ID),
	).LEFT_JOIN(
		t.Rsomember, t.Baseuser.ID.EQ(t.Rsomember.ID),
	)
}

func ReadUsers() pg.SelectStatement {
	return pg.SELECT(FullUserColumns()).FROM(FullUserTable())
}

type User struct {
	m.Baseuser
	*m.Student
	*m.Superadmin
	*m.Rsomember
}

func CreateUniversity() pg.InsertStatement { return t.University.INSERT(t.University.MutableColumns) }

func ReadUniversities() pg.SelectStatement {
	return pg.SELECT(
		t.University.AllColumns,
	).FROM(
		t.University,
	)
}

func CreateStudent() pg.InsertStatement { return t.Student.INSERT(t.Student.ID) }

type Rso struct {
	m.Rso
	m.University
	Tags    []m.Tag
	Members []User
}

// Query that selects all RSOs, their associated university data and members
func ReadRsos() pg.SelectStatement {
	rsoUniversityBool := t.Rso.UniversityID.EQ(t.University.ID)                           // Match university with RSO
	rsoMemberBool := t.Rso.ID.EQ(t.Rsomember.RsoID).AND(t.Rsomember.ID.EQ(t.Baseuser.ID)) // Match Rso member to Rso

	// Table that MUST include Rso,
	// MUST include its associated university,
	// MAY include rso members,
	// MAY include tags
	table := t.Rso.INNER_JOIN(
		t.University, rsoUniversityBool,
	).LEFT_JOIN(
		FullUserTable(), rsoMemberBool,
	).LEFT_JOIN(
		t.Taggedrso, t.Taggedrso.RsoID.EQ(t.Rso.ID),
	).LEFT_JOIN(
		t.Tag, t.Tag.ID.EQ(t.Taggedrso.TagID).AND(t.Taggedrso.RsoID.EQ(t.Rso.ID)),
	)

	return t.Rso.SELECT(t.Rso.AllColumns, t.University.AllColumns, t.Tag.AllColumns, FullUserColumns()).FROM(table)
}

// Query that selects all RSOs and their associated university data, that have at least 5 members
func ReadRsosValid() pg.SelectStatement {

	// Count of RSO members that are part of the same RSO
	rsoMemberCount := pg.COUNT(t.Rsomember.RsoID.EQ(t.Rso.ID)).GT(pg.Int(4))

	return ReadRsos().WHERE(rsoMemberCount)
}

func CreateComment() pg.InsertStatement {
	return t.Comment.INSERT(t.Comment.MutableColumns)
}

func ReadComment() pg.SelectStatement {
	return t.Comment.SELECT(t.Comment.AllColumns)
}

func CreateRsoMember() pg.InsertStatement {
	return t.Rsomember.INSERT(t.Rsomember.ID, t.Rsomember.RsoID)
}
