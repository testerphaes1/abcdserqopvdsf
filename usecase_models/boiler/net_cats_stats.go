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

// NetCatsStat is an object representing the database table.
type NetCatsStat struct {
	Time         time.Time   `boil:"time" json:"time" toml:"time" yaml:"time"`
	SessionID    string      `boil:"session_id" json:"session_id" toml:"session_id" yaml:"session_id"`
	ProjectID    int         `boil:"project_id" json:"project_id" toml:"project_id" yaml:"project_id"`
	NetcatID     int         `boil:"netcat_id" json:"netcat_id" toml:"netcat_id" yaml:"netcat_id"`
	URL          null.String `boil:"url" json:"url,omitempty" toml:"url" yaml:"url,omitempty"`
	DatacenterID int         `boil:"datacenter_id" json:"datacenter_id" toml:"datacenter_id" yaml:"datacenter_id"`
	IsHeartBeat  bool        `boil:"is_heart_beat" json:"is_heart_beat" toml:"is_heart_beat" yaml:"is_heart_beat"`
	Success      int         `boil:"success" json:"success" toml:"success" yaml:"success"`

	R *netCatsStatR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L netCatsStatL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var NetCatsStatColumns = struct {
	Time         string
	SessionID    string
	ProjectID    string
	NetcatID     string
	URL          string
	DatacenterID string
	IsHeartBeat  string
	Success      string
}{
	Time:         "time",
	SessionID:    "session_id",
	ProjectID:    "project_id",
	NetcatID:     "netcat_id",
	URL:          "url",
	DatacenterID: "datacenter_id",
	IsHeartBeat:  "is_heart_beat",
	Success:      "success",
}

var NetCatsStatTableColumns = struct {
	Time         string
	SessionID    string
	ProjectID    string
	NetcatID     string
	URL          string
	DatacenterID string
	IsHeartBeat  string
	Success      string
}{
	Time:         "net_cats_stats.time",
	SessionID:    "net_cats_stats.session_id",
	ProjectID:    "net_cats_stats.project_id",
	NetcatID:     "net_cats_stats.netcat_id",
	URL:          "net_cats_stats.url",
	DatacenterID: "net_cats_stats.datacenter_id",
	IsHeartBeat:  "net_cats_stats.is_heart_beat",
	Success:      "net_cats_stats.success",
}

// Generated where

var NetCatsStatWhere = struct {
	Time         whereHelpertime_Time
	SessionID    whereHelperstring
	ProjectID    whereHelperint
	NetcatID     whereHelperint
	URL          whereHelpernull_String
	DatacenterID whereHelperint
	IsHeartBeat  whereHelperbool
	Success      whereHelperint
}{
	Time:         whereHelpertime_Time{field: "\"net_cats_stats\".\"time\""},
	SessionID:    whereHelperstring{field: "\"net_cats_stats\".\"session_id\""},
	ProjectID:    whereHelperint{field: "\"net_cats_stats\".\"project_id\""},
	NetcatID:     whereHelperint{field: "\"net_cats_stats\".\"netcat_id\""},
	URL:          whereHelpernull_String{field: "\"net_cats_stats\".\"url\""},
	DatacenterID: whereHelperint{field: "\"net_cats_stats\".\"datacenter_id\""},
	IsHeartBeat:  whereHelperbool{field: "\"net_cats_stats\".\"is_heart_beat\""},
	Success:      whereHelperint{field: "\"net_cats_stats\".\"success\""},
}

// NetCatsStatRels is where relationship names are stored.
var NetCatsStatRels = struct {
	Datacenter string
	Netcat     string
}{
	Datacenter: "Datacenter",
	Netcat:     "Netcat",
}

// netCatsStatR is where relationships are stored.
type netCatsStatR struct {
	Datacenter *Datacenter `boil:"Datacenter" json:"Datacenter" toml:"Datacenter" yaml:"Datacenter"`
	Netcat     *NetCat     `boil:"Netcat" json:"Netcat" toml:"Netcat" yaml:"Netcat"`
}

// NewStruct creates a new relationship struct
func (*netCatsStatR) NewStruct() *netCatsStatR {
	return &netCatsStatR{}
}

func (r *netCatsStatR) GetDatacenter() *Datacenter {
	if r == nil {
		return nil
	}
	return r.Datacenter
}

func (r *netCatsStatR) GetNetcat() *NetCat {
	if r == nil {
		return nil
	}
	return r.Netcat
}

// netCatsStatL is where Load methods for each relationship are stored.
type netCatsStatL struct{}

var (
	netCatsStatAllColumns            = []string{"time", "session_id", "project_id", "netcat_id", "url", "datacenter_id", "is_heart_beat", "success"}
	netCatsStatColumnsWithoutDefault = []string{"time", "session_id", "project_id", "netcat_id", "datacenter_id", "is_heart_beat", "success"}
	netCatsStatColumnsWithDefault    = []string{"url"}
	netCatsStatPrimaryKeyColumns     = []string{"time", "netcat_id"}
	netCatsStatGeneratedColumns      = []string{}
)

type (
	// NetCatsStatSlice is an alias for a slice of pointers to NetCatsStat.
	// This should almost always be used instead of []NetCatsStat.
	NetCatsStatSlice []*NetCatsStat
	// NetCatsStatHook is the signature for custom NetCatsStat hook methods
	NetCatsStatHook func(context.Context, boil.ContextExecutor, *NetCatsStat) error

	netCatsStatQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	netCatsStatType                 = reflect.TypeOf(&NetCatsStat{})
	netCatsStatMapping              = queries.MakeStructMapping(netCatsStatType)
	netCatsStatPrimaryKeyMapping, _ = queries.BindMapping(netCatsStatType, netCatsStatMapping, netCatsStatPrimaryKeyColumns)
	netCatsStatInsertCacheMut       sync.RWMutex
	netCatsStatInsertCache          = make(map[string]insertCache)
	netCatsStatUpdateCacheMut       sync.RWMutex
	netCatsStatUpdateCache          = make(map[string]updateCache)
	netCatsStatUpsertCacheMut       sync.RWMutex
	netCatsStatUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var netCatsStatAfterSelectHooks []NetCatsStatHook

var netCatsStatBeforeInsertHooks []NetCatsStatHook
var netCatsStatAfterInsertHooks []NetCatsStatHook

var netCatsStatBeforeUpdateHooks []NetCatsStatHook
var netCatsStatAfterUpdateHooks []NetCatsStatHook

var netCatsStatBeforeDeleteHooks []NetCatsStatHook
var netCatsStatAfterDeleteHooks []NetCatsStatHook

var netCatsStatBeforeUpsertHooks []NetCatsStatHook
var netCatsStatAfterUpsertHooks []NetCatsStatHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *NetCatsStat) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *NetCatsStat) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *NetCatsStat) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *NetCatsStat) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *NetCatsStat) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *NetCatsStat) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *NetCatsStat) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *NetCatsStat) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *NetCatsStat) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netCatsStatAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddNetCatsStatHook registers your hook function for all future operations.
func AddNetCatsStatHook(hookPoint boil.HookPoint, netCatsStatHook NetCatsStatHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		netCatsStatAfterSelectHooks = append(netCatsStatAfterSelectHooks, netCatsStatHook)
	case boil.BeforeInsertHook:
		netCatsStatBeforeInsertHooks = append(netCatsStatBeforeInsertHooks, netCatsStatHook)
	case boil.AfterInsertHook:
		netCatsStatAfterInsertHooks = append(netCatsStatAfterInsertHooks, netCatsStatHook)
	case boil.BeforeUpdateHook:
		netCatsStatBeforeUpdateHooks = append(netCatsStatBeforeUpdateHooks, netCatsStatHook)
	case boil.AfterUpdateHook:
		netCatsStatAfterUpdateHooks = append(netCatsStatAfterUpdateHooks, netCatsStatHook)
	case boil.BeforeDeleteHook:
		netCatsStatBeforeDeleteHooks = append(netCatsStatBeforeDeleteHooks, netCatsStatHook)
	case boil.AfterDeleteHook:
		netCatsStatAfterDeleteHooks = append(netCatsStatAfterDeleteHooks, netCatsStatHook)
	case boil.BeforeUpsertHook:
		netCatsStatBeforeUpsertHooks = append(netCatsStatBeforeUpsertHooks, netCatsStatHook)
	case boil.AfterUpsertHook:
		netCatsStatAfterUpsertHooks = append(netCatsStatAfterUpsertHooks, netCatsStatHook)
	}
}

// One returns a single netCatsStat record from the query.
func (q netCatsStatQuery) One(ctx context.Context, exec boil.ContextExecutor) (*NetCatsStat, error) {
	o := &NetCatsStat{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for net_cats_stats")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all NetCatsStat records from the query.
func (q netCatsStatQuery) All(ctx context.Context, exec boil.ContextExecutor) (NetCatsStatSlice, error) {
	var o []*NetCatsStat

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to NetCatsStat slice")
	}

	if len(netCatsStatAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all NetCatsStat records in the query.
func (q netCatsStatQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count net_cats_stats rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q netCatsStatQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if net_cats_stats exists")
	}

	return count > 0, nil
}

// Datacenter pointed to by the foreign key.
func (o *NetCatsStat) Datacenter(mods ...qm.QueryMod) datacenterQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.DatacenterID),
	}

	queryMods = append(queryMods, mods...)

	return Datacenters(queryMods...)
}

// Netcat pointed to by the foreign key.
func (o *NetCatsStat) Netcat(mods ...qm.QueryMod) netCatQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.NetcatID),
	}

	queryMods = append(queryMods, mods...)

	return NetCats(queryMods...)
}

// LoadDatacenter allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (netCatsStatL) LoadDatacenter(ctx context.Context, e boil.ContextExecutor, singular bool, maybeNetCatsStat interface{}, mods queries.Applicator) error {
	var slice []*NetCatsStat
	var object *NetCatsStat

	if singular {
		var ok bool
		object, ok = maybeNetCatsStat.(*NetCatsStat)
		if !ok {
			object = new(NetCatsStat)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeNetCatsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeNetCatsStat))
			}
		}
	} else {
		s, ok := maybeNetCatsStat.(*[]*NetCatsStat)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeNetCatsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeNetCatsStat))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &netCatsStatR{}
		}
		args = append(args, object.DatacenterID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &netCatsStatR{}
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
		foreign.R.NetCatsStats = append(foreign.R.NetCatsStats, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.DatacenterID == foreign.ID {
				local.R.Datacenter = foreign
				if foreign.R == nil {
					foreign.R = &datacenterR{}
				}
				foreign.R.NetCatsStats = append(foreign.R.NetCatsStats, local)
				break
			}
		}
	}

	return nil
}

// LoadNetcat allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (netCatsStatL) LoadNetcat(ctx context.Context, e boil.ContextExecutor, singular bool, maybeNetCatsStat interface{}, mods queries.Applicator) error {
	var slice []*NetCatsStat
	var object *NetCatsStat

	if singular {
		var ok bool
		object, ok = maybeNetCatsStat.(*NetCatsStat)
		if !ok {
			object = new(NetCatsStat)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeNetCatsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeNetCatsStat))
			}
		}
	} else {
		s, ok := maybeNetCatsStat.(*[]*NetCatsStat)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeNetCatsStat)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeNetCatsStat))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &netCatsStatR{}
		}
		args = append(args, object.NetcatID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &netCatsStatR{}
			}

			for _, a := range args {
				if a == obj.NetcatID {
					continue Outer
				}
			}

			args = append(args, obj.NetcatID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`net_cats`),
		qm.WhereIn(`net_cats.id in ?`, args...),
		qmhelper.WhereIsNull(`net_cats.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load NetCat")
	}

	var resultSlice []*NetCat
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice NetCat")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for net_cats")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for net_cats")
	}

	if len(netCatAfterSelectHooks) != 0 {
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
		object.R.Netcat = foreign
		if foreign.R == nil {
			foreign.R = &netCatR{}
		}
		foreign.R.NetcatNetCatsStats = append(foreign.R.NetcatNetCatsStats, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.NetcatID == foreign.ID {
				local.R.Netcat = foreign
				if foreign.R == nil {
					foreign.R = &netCatR{}
				}
				foreign.R.NetcatNetCatsStats = append(foreign.R.NetcatNetCatsStats, local)
				break
			}
		}
	}

	return nil
}

// SetDatacenter of the netCatsStat to the related item.
// Sets o.R.Datacenter to related.
// Adds o to related.R.NetCatsStats.
func (o *NetCatsStat) SetDatacenter(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Datacenter) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"net_cats_stats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"datacenter_id"}),
		strmangle.WhereClause("\"", "\"", 2, netCatsStatPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.Time, o.NetcatID}

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
		o.R = &netCatsStatR{
			Datacenter: related,
		}
	} else {
		o.R.Datacenter = related
	}

	if related.R == nil {
		related.R = &datacenterR{
			NetCatsStats: NetCatsStatSlice{o},
		}
	} else {
		related.R.NetCatsStats = append(related.R.NetCatsStats, o)
	}

	return nil
}

// SetNetcat of the netCatsStat to the related item.
// Sets o.R.Netcat to related.
// Adds o to related.R.NetcatNetCatsStats.
func (o *NetCatsStat) SetNetcat(ctx context.Context, exec boil.ContextExecutor, insert bool, related *NetCat) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"net_cats_stats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"netcat_id"}),
		strmangle.WhereClause("\"", "\"", 2, netCatsStatPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.Time, o.NetcatID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.NetcatID = related.ID
	if o.R == nil {
		o.R = &netCatsStatR{
			Netcat: related,
		}
	} else {
		o.R.Netcat = related
	}

	if related.R == nil {
		related.R = &netCatR{
			NetcatNetCatsStats: NetCatsStatSlice{o},
		}
	} else {
		related.R.NetcatNetCatsStats = append(related.R.NetcatNetCatsStats, o)
	}

	return nil
}

// NetCatsStats retrieves all the records using an executor.
func NetCatsStats(mods ...qm.QueryMod) netCatsStatQuery {
	mods = append(mods, qm.From("\"net_cats_stats\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"net_cats_stats\".*"})
	}

	return netCatsStatQuery{q}
}

// FindNetCatsStat retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindNetCatsStat(ctx context.Context, exec boil.ContextExecutor, time time.Time, netcatID int, selectCols ...string) (*NetCatsStat, error) {
	netCatsStatObj := &NetCatsStat{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"net_cats_stats\" where \"time\"=$1 AND \"netcat_id\"=$2", sel,
	)

	q := queries.Raw(query, time, netcatID)

	err := q.Bind(ctx, exec, netCatsStatObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from net_cats_stats")
	}

	if err = netCatsStatObj.doAfterSelectHooks(ctx, exec); err != nil {
		return netCatsStatObj, err
	}

	return netCatsStatObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *NetCatsStat) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no net_cats_stats provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(netCatsStatColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	netCatsStatInsertCacheMut.RLock()
	cache, cached := netCatsStatInsertCache[key]
	netCatsStatInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			netCatsStatAllColumns,
			netCatsStatColumnsWithDefault,
			netCatsStatColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(netCatsStatType, netCatsStatMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(netCatsStatType, netCatsStatMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"net_cats_stats\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"net_cats_stats\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into net_cats_stats")
	}

	if !cached {
		netCatsStatInsertCacheMut.Lock()
		netCatsStatInsertCache[key] = cache
		netCatsStatInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the NetCatsStat.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *NetCatsStat) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	netCatsStatUpdateCacheMut.RLock()
	cache, cached := netCatsStatUpdateCache[key]
	netCatsStatUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			netCatsStatAllColumns,
			netCatsStatPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update net_cats_stats, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"net_cats_stats\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, netCatsStatPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(netCatsStatType, netCatsStatMapping, append(wl, netCatsStatPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update net_cats_stats row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for net_cats_stats")
	}

	if !cached {
		netCatsStatUpdateCacheMut.Lock()
		netCatsStatUpdateCache[key] = cache
		netCatsStatUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q netCatsStatQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for net_cats_stats")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for net_cats_stats")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o NetCatsStatSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), netCatsStatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"net_cats_stats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, netCatsStatPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in netCatsStat slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all netCatsStat")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *NetCatsStat) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no net_cats_stats provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(netCatsStatColumnsWithDefault, o)

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

	netCatsStatUpsertCacheMut.RLock()
	cache, cached := netCatsStatUpsertCache[key]
	netCatsStatUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			netCatsStatAllColumns,
			netCatsStatColumnsWithDefault,
			netCatsStatColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			netCatsStatAllColumns,
			netCatsStatPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert net_cats_stats, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(netCatsStatPrimaryKeyColumns))
			copy(conflict, netCatsStatPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"net_cats_stats\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(netCatsStatType, netCatsStatMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(netCatsStatType, netCatsStatMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert net_cats_stats")
	}

	if !cached {
		netCatsStatUpsertCacheMut.Lock()
		netCatsStatUpsertCache[key] = cache
		netCatsStatUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single NetCatsStat record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *NetCatsStat) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no NetCatsStat provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), netCatsStatPrimaryKeyMapping)
	sql := "DELETE FROM \"net_cats_stats\" WHERE \"time\"=$1 AND \"netcat_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from net_cats_stats")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for net_cats_stats")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q netCatsStatQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no netCatsStatQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from net_cats_stats")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for net_cats_stats")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o NetCatsStatSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(netCatsStatBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), netCatsStatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"net_cats_stats\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, netCatsStatPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from netCatsStat slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for net_cats_stats")
	}

	if len(netCatsStatAfterDeleteHooks) != 0 {
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
func (o *NetCatsStat) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindNetCatsStat(ctx, exec, o.Time, o.NetcatID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *NetCatsStatSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := NetCatsStatSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), netCatsStatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"net_cats_stats\".* FROM \"net_cats_stats\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, netCatsStatPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in NetCatsStatSlice")
	}

	*o = slice

	return nil
}

// NetCatsStatExists checks if the NetCatsStat row exists.
func NetCatsStatExists(ctx context.Context, exec boil.ContextExecutor, time time.Time, netcatID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"net_cats_stats\" where \"time\"=$1 AND \"netcat_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, time, netcatID)
	}
	row := exec.QueryRowContext(ctx, sql, time, netcatID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if net_cats_stats exists")
	}

	return exists, nil
}

// Exists checks if the NetCatsStat row exists.
func (o *NetCatsStat) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return NetCatsStatExists(ctx, exec, o.Time, o.NetcatID)
}
