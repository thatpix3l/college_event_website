package gen_sql

import (
	s "github.com/go-jet/jet/v2/postgres"
	m "github.com/thatpix3l/collge_event_website/src/gen_sql/college_event_website/cew/model"
	t "github.com/thatpix3l/collge_event_website/src/gen_sql/college_event_website/cew/table"
)

func queryToFunc[T any](query T) func() T {
	return func() T {
		return query
	}
}

var ReadEvents = queryToFunc(s.SELECT(
	t.Baseevent.AllColumns,
	t.Rsoevent.AllColumns,
	t.Publicevent.AllColumns,
	t.Privateevent.AllColumns,
).FROM(
	t.Baseevent.LEFT_JOIN(
		t.Rsoevent, t.Baseevent.ID.EQ(t.Rsoevent.ID),
	).LEFT_JOIN(
		t.Publicevent, t.Baseevent.ID.EQ(t.Publicevent.ID),
	).LEFT_JOIN(
		t.Privateevent, t.Baseevent.ID.EQ(t.Privateevent.ID),
	),
))

type Event struct {
	m.Baseevent
	*m.Rsoevent
	*m.Publicevent
	*m.Privateevent
}

var CreateBaseUser = queryToFunc(t.Baseuser.INSERT(t.Baseuser.MutableColumns))

var ReadUsers = queryToFunc(s.SELECT(
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
))

type User struct {
	m.Baseuser
	*m.Student
	*m.Superadmin
	*m.Rsomember
}

var CreateUniversity = queryToFunc(t.University.INSERT(t.University.MutableColumns))

var ReadUniversities = queryToFunc(s.SELECT(
	t.University.AllColumns,
).FROM(
	t.University,
))

var CreateStudent = queryToFunc(t.Student.INSERT(t.Student.ID))
