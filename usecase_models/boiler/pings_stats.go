// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// PingsStat is an object representing the database table.
type PingsStat struct {
	Time         time.Time   `boil:"time" json:"time" toml:"time" yaml:"time"`
	SessionID    string      `boil:"session_id" json:"session_id" toml:"session_id" yaml:"session_id"`
	ProjectID    int         `boil:"project_id" json:"project_id" toml:"project_id" yaml:"project_id"`
	PingID       int         `boil:"ping_id" json:"ping_id" toml:"ping_id" yaml:"ping_id"`
	URL          null.String `boil:"url" json:"url,omitempty" toml:"url" yaml:"url,omitempty"`
	DatacenterID int         `boil:"datacenter_id" json:"datacenter_id" toml:"datacenter_id" yaml:"datacenter_id"`
	IsHeartBeat  bool        `boil:"is_heart_beat" json:"is_heart_beat" toml:"is_heart_beat" yaml:"is_heart_beat"`
	Success      int         `boil:"success" json:"success" toml:"success" yaml:"success"`

	R *pingsStatR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L pingsStatL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PingsStatColumns = struct {
	Time         string
	SessionID    string
	ProjectID    string
	PingID       string
	URL          string
	DatacenterID string
	IsHeartBeat  string
	Success      string
}{
	Time:         "time",
	SessionID:    "session_id",
	ProjectID:    "project_id",
	PingID:       "ping_id",
	URL:          "url",
	DatacenterID: "datacenter_id",
	IsHeartBeat:  "is_heart_beat",
	Success:      "success",
}

var PingsStatTableColumns = struct {
	Time         string
	SessionID    string
	ProjectID    string
	PingID       string
	URL          string
	DatacenterID string
	IsHeartBeat  string
	Success      string
}{
	Time:         "pings_stats.time",
	SessionID:    "pings_stats.session_id",
	ProjectID:    "pings_stats.project_id",
	PingID:       "pings_stats.ping_id",
	URL:          "pings_stats.url",
	DatacenterID: "pings_stats.datacenter_id",
	IsHeartBeat:  "pings_stats.is_heart_beat",
	Success:      "pings_stats.success",
}

// Generated where

var PingsStatWhere = struct {
	Time         whereHelpertime_Time
	SessionID    whereHelperstring
	ProjectID    whereHelperint
	PingID       whereHelperint
	URL          whereHelpernull_String
	DatacenterID whereHelperint
	IsHeartBeat  whereHelperbool
	Success      whereHelperint
}{
	Time:         whereHelpertime_Time{field: "\"pings_stats\".\"time\""},
	SessionID:    whereHelperstring{field: "\"pings_stats\".\"session_id\""},
	ProjectID:    whereHelperint{field: "\"pings_stats\".\"project_id\""},
	PingID:       whereHelperint{field: "\"pings_stats\".\"ping_id\""},
	URL:          whereHelpernull_String{field: "\"pings_stats\".\"url\""},
	DatacenterID: whereHelperint{field: "\"pings_stats\".\"datacenter_id\""},
	IsHeartBeat:  whereHelperbool{field: "\"pings_stats\".\"is_heart_beat\""},
	Success:      whereHelperint{field: "\"pings_stats\".\"success\""},
}

// PingsStatRels is where relationship names are stored.
var PingsStatRels = struct {
	Datacenter string
	Ping       string
}{
	Datacenter: "Datacenter",
	Ping:       "Ping",
}

// pingsStatR is where relationships are stored.
type pingsStatR struct {
	Datacenter *Datacenter `boil:"Datacenter" json:"Datacenter" toml:"Datacenter" yaml:"Datacenter"`
	Ping       *Ping       `boil:"Ping" json:"Ping" toml:"Ping" yaml:"Ping"`
}

// NewStruct creates a new relationship struct
func (*pingsStatR) NewStruct() *pingsStatR {
	return &pingsStatR{}
}

func (r *pingsStatR) GetDatacenter() *Datacenter {
	if r == nil {
		return nil
	}
	return r.Datacenter
}

func (r *pingsStatR) GetPing() *Ping {
	if r == nil {
		return nil
	}
	return r.Ping
}

// pingsStatL is where Load methods for each relationship are stored.
type pingsStatL struct{}

var (
	pingsStatAllColumns            = []string{"time", "session_id", "project_id", "ping_id", "url", "datacenter_id", "is_heart_beat", "success"}
	pingsStatColumnsWithoutDefault = []string{"time", "session_id", "project_id", "ping_id", "datacenter_id", "is_heart_beat", "success"}
	pingsStatColumnsWithDefault    = []string{"url"}
	pingsStatPrimaryKeyColumns     = []string{"time", "ping_id"}
	pingsStatGeneratedColumns      = []string{}
)

type (
	// PingsStatSlice is an alias for a slice of pointers to PingsStat.
	// This should almost always be used instead of []PingsStat.
	PingsStatSlice []*PingsStat
	// PingsStatHook is the signature for custom PingsStat hook methods
	PingsStatHook func(context.Context, boil.ContextExecutor, *PingsStat) error

	pingsStatQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	pingsStatType                 = reflect.TypeOf(&PingsStat{})
	pingsStatMapping              = queries.MakeStructMapping(pingsStatType)
	pingsStatPrimaryKeyMapping, _ = queries.BindMapping(pingsStatType, pingsStatMapping, pingsStatPrimaryKeyColumns)
	pingsStatInsertCacheMut       sync.RWMutex
	pingsStatInsertCache          = make(map[string]insertCache)
	pingsStatUpdateCacheMut       sync.RWMutex
	pingsStatUpdateCache          = make(map[string]updateCache)
	pingsStatUpsertCacheMut       sync.RWMutex
	pingsStatUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var pingsStatAfterSelectHooks []PingsStatHook

var pingsStatBeforeInsertHooks []PingsStatHook
var pingsStatAfterInsertHooks []PingsStatHook

var pingsStatBeforeUpdateHooks []PingsStatHook
var pingsStatAfterUpdateHooks []PingsStatHook

var pingsStatBeforeDeleteHooks []PingsStatHook
var pingsStatAfterDeleteHooks []PingsStatHook

var pingsStatBeforeUpsertHooks []PingsStatHook
var pingsStatAfterUpsertHooks []PingsStatHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *PingsStat) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *PingsStat) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *PingsStat) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *PingsStat) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *PingsStat) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *PingsStat) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *PingsStat) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *PingsStat) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *PingsStat) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range pingsStatAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPingsStatHook registers your hook function for all future operations.
func AddPingsStatHook(hookPoint boil.HookPoint, pingsStatHook PingsStatHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		pingsStatAfterSelectHooks = append(pingsStatAfterSelectHooks, pingsStatHook)
	case boil.BeforeInsertHook:
		pingsStatBeforeInsertHooks = append(pingsStatBeforeInsertHooks, pingsStatHook)
	case boil.AfterInsertHook:
		pingsStatAfterInsertHooks = append(pingsStatAfterInsertHooks, pingsStatHook)
	case boil.BeforeUpdateHook:
		pingsStatBeforeUpdateHooks = append(pingsStatBeforeUpdateHooks, pingsStatHook)
	case boil.AfterUpdateHook:
		pingsStatAfterUpdateHooks = append(pingsStatAfterUpdateHooks, pingsStatHook)
	case boil.BeforeDeleteHook:
		pingsStatBeforeDeleteHooks = append(pingsStatBeforeDeleteHooks, pingsStatHook)
	case boil.AfterDeleteHook:
		pingsStatAfterDeleteHooks = append(pingsStatAfterDeleteHooks, pingsStatHook)
	case boil.BeforeUpsertHook:
		pingsStatBeforeUpsertHooks = append(pingsStatBeforeUpsertHooks, pingsStatHook)
	case boil.AfterUpsertHook:
		pingsStatAfterUpsertHooks = append(pingsStatAfterUpsertHooks, pingsStatHook)
	}
}

// One returns a single pingsStat record from the query.
func (q pingsStatQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PingsStat, error) {
	o := &PingsStat{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for pings_stats")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all PingsStat records from the query.
func (q pingsStatQuery) All(ctx context.Context, exec boil.ContextExecutor) (PingsStatSlice, error) {
	var o []*PingsStat

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PingsStat slice")
	}

	if len(pingsStatAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all PingsStat records in the query.
func (q pingsStatQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count pings_stats rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q pingsStatQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if pings_stats exists")
	}

	return count > 0, nil
}

// Datacenter pointed to by the foreign key.
func (o *PingsStat) Datacenter(mods ...qm.QueryMod) datacenterQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.DatacenterID),
	}

	queryMods = append(queryMods, mods...)

	return Datacenters(queryMods...)
}

// Ping pointed to by the foreign key.
func (o *PingsStat) Ping(mods ...qm.QueryMod) pingQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.PingID),
	}

	queryMods = append(queryMods, mods...)

	return Pings(queryMods...)
}

// LoadDatacenter allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (pingsStatL) LoadDatacenter(ctx context.Context, e boil.ContextExecutor, singular bool, maybePingsStat interface{}, mods queries.Applicator) error {
	var slice []*PingsStat
	var object *PingsStat

	if singular {
		var ok bool
		object, ok = maybePingsStat.(*PingsStat)
		if !ok {
			object = new(PingsStat)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePingsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePingsStat))
			}
		}
	} else {
		s, ok := maybePingsStat.(*[]*PingsStat)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePingsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePingsStat))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &pingsStatR{}
		}
		args = append(args, object.DatacenterID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &pingsStatR{}
			}

			for _, a := range args {
				if a == obj.DatacenterID {
					continue Outer
				}
			}

			args = append(args, obj.DatacenterID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`datacenters`),
		qm.WhereIn(`datacenters.id in ?`, args...),
		qmhelper.WhereIsNull(`datacenters.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Datacenter")
	}

	var resultSlice []*Datacenter
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Datacenter")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for datacenters")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for datacenters")
	}

	if len(datacenterAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Datacenter = foreign
		if foreign.R == nil {
			foreign.R = &datacenterR{}
		}
		foreign.R.PingsStats = append(foreign.R.PingsStats, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.DatacenterID == foreign.ID {
				local.R.Datacenter = foreign
				if foreign.R == nil {
					foreign.R = &datacenterR{}
				}
				foreign.R.PingsStats = append(foreign.R.PingsStats, local)
				break
			}
		}
	}

	return nil
}

// LoadPing allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (pingsStatL) LoadPing(ctx context.Context, e boil.ContextExecutor, singular bool, maybePingsStat interface{}, mods queries.Applicator) error {
	var slice []*PingsStat
	var object *PingsStat

	if singular {
		var ok bool
		object, ok = maybePingsStat.(*PingsStat)
		if !ok {
			object = new(PingsStat)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePingsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePingsStat))
			}
		}
	} else {
		s, ok := maybePingsStat.(*[]*PingsStat)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePingsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePingsStat))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &pingsStatR{}
		}
		args = append(args, object.PingID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &pingsStatR{}
			}

			for _, a := range args {
				if a == obj.PingID {
					continue Outer
				}
			}

			args = append(args, obj.PingID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`pings`),
		qm.WhereIn(`pings.id in ?`, args...),
		qmhelper.WhereIsNull(`pings.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Ping")
	}

	var resultSlice []*Ping
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Ping")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for pings")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for pings")
	}

	if len(pingAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Ping = foreign
		if foreign.R == nil {
			foreign.R = &pingR{}
		}
		foreign.R.PingsStats = append(foreign.R.PingsStats, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PingID == foreign.ID {
				local.R.Ping = foreign
				if foreign.R == nil {
					foreign.R = &pingR{}
				}
				foreign.R.PingsStats = append(foreign.R.PingsStats, local)
				break
			}
		}
	}

	return nil
}

// SetDatacenter of the pingsStat to the related item.
// Sets o.R.Datacenter to related.
// Adds o to related.R.PingsStats.
func (o *PingsStat) SetDatacenter(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Datacenter) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"pings_stats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"datacenter_id"}),
		strmangle.WhereClause("\"", "\"", 2, pingsStatPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.Time, o.PingID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.DatacenterID = related.ID
	if o.R == nil {
		o.R = &pingsStatR{
			Datacenter: related,
		}
	} else {
		o.R.Datacenter = related
	}

	if related.R == nil {
		related.R = &datacenterR{
			PingsStats: PingsStatSlice{o},
		}
	} else {
		related.R.PingsStats = append(related.R.PingsStats, o)
	}

	return nil
}

// SetPing of the pingsStat to the related item.
// Sets o.R.Ping to related.
// Adds o to related.R.PingsStats.
func (o *PingsStat) SetPing(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Ping) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"pings_stats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"ping_id"}),
		strmangle.WhereClause("\"", "\"", 2, pingsStatPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.Time, o.PingID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PingID = related.ID
	if o.R == nil {
		o.R = &pingsStatR{
			Ping: related,
		}
	} else {
		o.R.Ping = related
	}

	if related.R == nil {
		related.R = &pingR{
			PingsStats: PingsStatSlice{o},
		}
	} else {
		related.R.PingsStats = append(related.R.PingsStats, o)
	}

	return nil
}

// PingsStats retrieves all the records using an executor.
func PingsStats(mods ...qm.QueryMod) pingsStatQuery {
	mods = append(mods, qm.From("\"pings_stats\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"pings_stats\".*"})
	}

	return pingsStatQuery{q}
}

// FindPingsStat retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPingsStat(ctx context.Context, exec boil.ContextExecutor, time time.Time, pingID int, selectCols ...string) (*PingsStat, error) {
	pingsStatObj := &PingsStat{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"pings_stats\" where \"time\"=$1 AND \"ping_id\"=$2", sel,
	)

	q := queries.Raw(query, time, pingID)

	err := q.Bind(ctx, exec, pingsStatObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from pings_stats")
	}

	if err = pingsStatObj.doAfterSelectHooks(ctx, exec); err != nil {
		return pingsStatObj, err
	}

	return pingsStatObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PingsStat) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no pings_stats provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(pingsStatColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	pingsStatInsertCacheMut.RLock()
	cache, cached := pingsStatInsertCache[key]
	pingsStatInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			pingsStatAllColumns,
			pingsStatColumnsWithDefault,
			pingsStatColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(pingsStatType, pingsStatMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(pingsStatType, pingsStatMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"pings_stats\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"pings_stats\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into pings_stats")
	}

	if !cached {
		pingsStatInsertCacheMut.Lock()
		pingsStatInsertCache[key] = cache
		pingsStatInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the PingsStat.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PingsStat) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	pingsStatUpdateCacheMut.RLock()
	cache, cached := pingsStatUpdateCache[key]
	pingsStatUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			pingsStatAllColumns,
			pingsStatPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update pings_stats, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"pings_stats\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, pingsStatPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(pingsStatType, pingsStatMapping, append(wl, pingsStatPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update pings_stats row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for pings_stats")
	}

	if !cached {
		pingsStatUpdateCacheMut.Lock()
		pingsStatUpdateCache[key] = cache
		pingsStatUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q pingsStatQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for pings_stats")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for pings_stats")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PingsStatSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pingsStatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"pings_stats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, pingsStatPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in pingsStat slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all pingsStat")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PingsStat) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no pings_stats provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(pingsStatColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	pingsStatUpsertCacheMut.RLock()
	cache, cached := pingsStatUpsertCache[key]
	pingsStatUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			pingsStatAllColumns,
			pingsStatColumnsWithDefault,
			pingsStatColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			pingsStatAllColumns,
			pingsStatPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert pings_stats, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(pingsStatPrimaryKeyColumns))
			copy(conflict, pingsStatPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"pings_stats\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(pingsStatType, pingsStatMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(pingsStatType, pingsStatMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert pings_stats")
	}

	if !cached {
		pingsStatUpsertCacheMut.Lock()
		pingsStatUpsertCache[key] = cache
		pingsStatUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single PingsStat record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PingsStat) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no PingsStat provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), pingsStatPrimaryKeyMapping)
	sql := "DELETE FROM \"pings_stats\" WHERE \"time\"=$1 AND \"ping_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from pings_stats")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for pings_stats")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q pingsStatQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no pingsStatQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from pings_stats")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for pings_stats")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PingsStatSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(pingsStatBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pingsStatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"pings_stats\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pingsStatPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from pingsStat slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for pings_stats")
	}

	if len(pingsStatAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PingsStat) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPingsStat(ctx, exec, o.Time, o.PingID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PingsStatSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PingsStatSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pingsStatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"pings_stats\".* FROM \"pings_stats\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pingsStatPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PingsStatSlice")
	}

	*o = slice

	return nil
}

// PingsStatExists checks if the PingsStat row exists.
func PingsStatExists(ctx context.Context, exec boil.ContextExecutor, time time.Time, pingID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"pings_stats\" where \"time\"=$1 AND \"ping_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, time, pingID)
	}
	row := exec.QueryRowContext(ctx, sql, time, pingID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if pings_stats exists")
	}

	return exists, nil
}

// Exists checks if the PingsStat row exists.
func (o *PingsStat) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return PingsStatExists(ctx, exec, o.Time, o.PingID)
}
