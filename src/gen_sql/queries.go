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

// Query for creating event with only required parameters
func CreateEvent() pg.InsertStatement {
	return t.Baseevent.INSERT(t.Baseevent.MutableColumns)
}

func FullEventColumns() pg.ProjectionList {
	return []pg.Projection{
		t.Baseevent.AllColumns,
		t.Rsoevent.AllColumns,
		t.Publicevent.AllColumns,
		t.Privateevent.AllColumns,
		t.Tag.AllColumns,
		t.Comment.AllColumns,
		t.Rating.AllColumns,
	}
}

func FullEventTable() pg.ReadableTable {
	return t.Baseevent.LEFT_JOIN(
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
	)
}

func ReadEvents() pg.SelectStatement {
	return pg.SELECT(FullEventColumns()).FROM(FullEventTable())
}

func IsPublic() pg.BoolExpression {
	return t.Baseevent.ID.EQ(t.Publicevent.ID)
}

func IsApproved() pg.BoolExpression {
	return t.Publicevent.Approved
}

func IsPublicApproved() pg.BoolExpression {
	return IsPublic().AND(IsApproved())
}

func IsPrivate() pg.BoolExpression {
	return t.Baseevent.ID.EQ(t.Privateevent.ID)
}

func IsPrivateSameUniversity(universityId string) pg.BoolExpression {
	return IsPrivate().AND(t.Baseevent.UniversityID.EQ(pg.String(universityId)))
}

func IsRso() pg.BoolExpression {
	return t.Baseevent.ID.EQ(t.Rsoevent.ID)
}

func IsRsoSameSource(rsoId string) pg.BoolExpression {
	return IsRso().AND(t.Rsoevent.RsoID.EQ(pg.String(rsoId)))
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

func StudentUniversity() *t.UniversityTable {
	return t.University.AS("StudentUniversity")
}

func CreateBaseUser() pg.InsertStatement { return t.Baseuser.INSERT(t.Baseuser.MutableColumns) }

func FullUserColumns() pg.ProjectionList {
	return []pg.Projection{
		t.Baseuser.AllColumns,
		t.Student.AllColumns,
		t.Superadmin.AllColumns,
		t.Rsomember.AllColumns,
		StudentUniversity().AllColumns,
	}
}

func FullUserTable() pg.ReadableTable {
	return t.Baseuser.LEFT_JOIN(
		t.Student, t.Baseuser.ID.EQ(t.Student.ID),
	).LEFT_JOIN(
		t.Superadmin, t.Baseuser.ID.EQ(t.Superadmin.ID),
	).LEFT_JOIN(
		t.Rsomember, t.Baseuser.ID.EQ(t.Rsomember.ID),
	).LEFT_JOIN(
		StudentUniversity(), StudentUniversity().ID.EQ(t.Student.UniversityID),
	)
}

func ReadUsers() pg.SelectStatement {
	return pg.SELECT(FullUserColumns()).FROM(FullUserTable())
}

type StudentFull struct {
	m.Student
	m.University `alias:"StudentUniversity.*"`
}

type User struct {
	m.Baseuser
	*StudentFull
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

func CreateStudent() pg.InsertStatement {
	return t.Student.INSERT(t.Student.AllColumns)
}

func RsoUniversity() t.UniversityTable {
	return *t.University.AS("RsoUniversity")
}

type Rso struct {
	m.Rso
	m.University `alias:"RsoUniversity.*"`
	Tags         []m.Tag
	Members      []User
}

func FullRsoColumns() pg.ProjectionList {
	return []pg.Projection{
		t.Rso.AllColumns,
		RsoUniversity().AllColumns,
		t.Tag.AllColumns,
		FullUserColumns(),
	}
}

func FullRsoTable() pg.ReadableTable {
	// Match university with RSO
	rsoUniversityBool := t.Rso.UniversityID.EQ(RsoUniversity().ID)

	// Match Rso member to Rso
	rsoMemberBool := t.Rso.ID.EQ(t.Rsomember.RsoID).AND(t.Rsomember.ID.EQ(t.Baseuser.ID))

	// Table that MUST include Rso,
	// MUST include its associated university,
	// MAY include rso members,
	// MAY include tags
	table := t.Rso.INNER_JOIN(
		RsoUniversity(), rsoUniversityBool,
	).LEFT_JOIN(
		FullUserTable(), rsoMemberBool,
	).LEFT_JOIN(
		t.Taggedrso, t.Taggedrso.RsoID.EQ(t.Rso.ID),
	).LEFT_JOIN(
		t.Tag, t.Tag.ID.EQ(t.Taggedrso.TagID).AND(t.Taggedrso.RsoID.EQ(t.Rso.ID)),
	)

	return table
}

// Query that selects all RSOs, their associated university data, and members
func ReadRsos() pg.SelectStatement {
	return t.Rso.SELECT(FullRsoColumns()).FROM(FullRsoTable())
}

// Query that selects all RSOs that have the minimum allowed amount of members to create events
func ReadRsosMinMembers() pg.SelectStatement {

	// Column that stores count of members in a given RSO
	rsoMemberCount := pg.COUNT(t.Rso.ID.EQ(t.Rsomember.RsoID))

	// Condition that accepts only RSOs that have a mebmer count greater than 4
	rsoMemberCountHasMin := rsoMemberCount.GT(pg.Int(4))

	return pg.SELECT(t.Rso.ID, t.Rso.Title, rsoMemberCount).FROM(FullRsoTable()).GROUP_BY(t.Rso.ID).HAVING(rsoMemberCountHasMin)
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

// Bool expression for checking if any users in a given rso do NOT have given email domain
func RsoMembersDifferentEmail(rso Rso, emailDomain string) pg.BoolExpression {
	return t.Rso.ID.EQ(pg.String(rso.Rso.ID)).AND(
		t.Rsomember.RsoID.EQ(t.Rso.ID),
	).AND(
		t.Rsomember.ID.EQ(t.Baseuser.ID),
	).AND(
		pg.COUNT(t.Baseuser.Email.NOT_EQ(pg.String(emailDomain))).EQ(pg.Int(0)),
	)
}
