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

// Gateway is an object representing the database table.
type Gateway struct {
	ID             int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Baseurl        string    `boil:"baseurl" json:"baseurl" toml:"baseurl" yaml:"baseurl"`
	Title          string    `boil:"title" json:"title" toml:"title" yaml:"title"`
	ConnectionRate null.Int  `boil:"connection_rate" json:"connection_rate,omitempty" toml:"connection_rate" yaml:"connection_rate,omitempty"`
	IsActive       null.Bool `boil:"is_active" json:"is_active,omitempty" toml:"is_active" yaml:"is_active,omitempty"`
	IsDefault      null.Bool `boil:"is_default" json:"is_default,omitempty" toml:"is_default" yaml:"is_default,omitempty"`
	Data           null.JSON `boil:"data" json:"data,omitempty" toml:"data" yaml:"data,omitempty"`
	UpdatedAt      time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	CreatedAt      time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt      null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *gatewayR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L gatewayL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GatewayColumns = struct {
	ID             string
	Baseurl        string
	Title          string
	ConnectionRate string
	IsActive       string
	IsDefault      string
	Data           string
	UpdatedAt      string
	CreatedAt      string
	DeletedAt      string
}{
	ID:             "id",
	Baseurl:        "baseurl",
	Title:          "title",
	ConnectionRate: "connection_rate",
	IsActive:       "is_active",
	IsDefault:      "is_default",
	Data:           "data",
	UpdatedAt:      "updated_at",
	CreatedAt:      "created_at",
	DeletedAt:      "deleted_at",
}

var GatewayTableColumns = struct {
	ID             string
	Baseurl        string
	Title          string
	ConnectionRate string
	IsActive       string
	IsDefault      string
	Data           string
	UpdatedAt      string
	CreatedAt      string
	DeletedAt      string
}{
	ID:             "gateways.id",
	Baseurl:        "gateways.baseurl",
	Title:          "gateways.title",
	ConnectionRate: "gateways.connection_rate",
	IsActive:       "gateways.is_active",
	IsDefault:      "gateways.is_default",
	Data:           "gateways.data",
	UpdatedAt:      "gateways.updated_at",
	CreatedAt:      "gateways.created_at",
	DeletedAt:      "gateways.deleted_at",
}

// Generated where

type whereHelpernull_Bool struct{ field string }

func (w whereHelpernull_Bool) EQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Bool) NEQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Bool) LT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Bool) LTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Bool) GT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Bool) GTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Bool) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Bool) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var GatewayWhere = struct {
	ID             whereHelperint
	Baseurl        whereHelperstring
	Title          whereHelperstring
	ConnectionRate whereHelpernull_Int
	IsActive       whereHelpernull_Bool
	IsDefault      whereHelpernull_Bool
	Data           whereHelpernull_JSON
	UpdatedAt      whereHelpertime_Time
	CreatedAt      whereHelpertime_Time
	DeletedAt      whereHelpernull_Time
}{
	ID:             whereHelperint{field: "\"gateways\".\"id\""},
	Baseurl:        whereHelperstring{field: "\"gateways\".\"baseurl\""},
	Title:          whereHelperstring{field: "\"gateways\".\"title\""},
	ConnectionRate: whereHelpernull_Int{field: "\"gateways\".\"connection_rate\""},
	IsActive:       whereHelpernull_Bool{field: "\"gateways\".\"is_active\""},
	IsDefault:      whereHelpernull_Bool{field: "\"gateways\".\"is_default\""},
	Data:           whereHelpernull_JSON{field: "\"gateways\".\"data\""},
	UpdatedAt:      whereHelpertime_Time{field: "\"gateways\".\"updated_at\""},
	CreatedAt:      whereHelpertime_Time{field: "\"gateways\".\"created_at\""},
	DeletedAt:      whereHelpernull_Time{field: "\"gateways\".\"deleted_at\""},
}

// GatewayRels is where relationship names are stored.
var GatewayRels = struct {
	Orders string
}{
	Orders: "Orders",
}

// gatewayR is where relationships are stored.
type gatewayR struct {
	Orders OrderSlice `boil:"Orders" json:"Orders" toml:"Orders" yaml:"Orders"`
}

// NewStruct creates a new relationship struct
func (*gatewayR) NewStruct() *gatewayR {
	return &gatewayR{}
}

func (r *gatewayR) GetOrders() OrderSlice {
	if r == nil {
		return nil
	}
	return r.Orders
}

// gatewayL is where Load methods for each relationship are stored.
type gatewayL struct{}

var (
	gatewayAllColumns            = []string{"id", "baseurl", "title", "connection_rate", "is_active", "is_default", "data", "updated_at", "created_at", "deleted_at"}
	gatewayColumnsWithoutDefault = []string{"baseurl", "title", "created_at"}
	gatewayColumnsWithDefault    = []string{"id", "connection_rate", "is_active", "is_default", "data", "updated_at", "deleted_at"}
	gatewayPrimaryKeyColumns     = []string{"id"}
	gatewayGeneratedColumns      = []string{}
)

type (
	// GatewaySlice is an alias for a slice of pointers to Gateway.
	// This should almost always be used instead of []Gateway.
	GatewaySlice []*Gateway
	// GatewayHook is the signature for custom Gateway hook methods
	GatewayHook func(context.Context, boil.ContextExecutor, *Gateway) error

	gatewayQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	gatewayType                 = reflect.TypeOf(&Gateway{})
	gatewayMapping              = queries.MakeStructMapping(gatewayType)
	gatewayPrimaryKeyMapping, _ = queries.BindMapping(gatewayType, gatewayMapping, gatewayPrimaryKeyColumns)
	gatewayInsertCacheMut       sync.RWMutex
	gatewayInsertCache          = make(map[string]insertCache)
	gatewayUpdateCacheMut       sync.RWMutex
	gatewayUpdateCache          = make(map[string]updateCache)
	gatewayUpsertCacheMut       sync.RWMutex
	gatewayUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var gatewayAfterSelectHooks []GatewayHook

var gatewayBeforeInsertHooks []GatewayHook
var gatewayAfterInsertHooks []GatewayHook

var gatewayBeforeUpdateHooks []GatewayHook
var gatewayAfterUpdateHooks []GatewayHook

var gatewayBeforeDeleteHooks []GatewayHook
var gatewayAfterDeleteHooks []GatewayHook

var gatewayBeforeUpsertHooks []GatewayHook
var gatewayAfterUpsertHooks []GatewayHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Gateway) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Gateway) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Gateway) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Gateway) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Gateway) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Gateway) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Gateway) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Gateway) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Gateway) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gatewayAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddGatewayHook registers your hook function for all future operations.
func AddGatewayHook(hookPoint boil.HookPoint, gatewayHook GatewayHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		gatewayAfterSelectHooks = append(gatewayAfterSelectHooks, gatewayHook)
	case boil.BeforeInsertHook:
		gatewayBeforeInsertHooks = append(gatewayBeforeInsertHooks, gatewayHook)
	case boil.AfterInsertHook:
		gatewayAfterInsertHooks = append(gatewayAfterInsertHooks, gatewayHook)
	case boil.BeforeUpdateHook:
		gatewayBeforeUpdateHooks = append(gatewayBeforeUpdateHooks, gatewayHook)
	case boil.AfterUpdateHook:
		gatewayAfterUpdateHooks = append(gatewayAfterUpdateHooks, gatewayHook)
	case boil.BeforeDeleteHook:
		gatewayBeforeDeleteHooks = append(gatewayBeforeDeleteHooks, gatewayHook)
	case boil.AfterDeleteHook:
		gatewayAfterDeleteHooks = append(gatewayAfterDeleteHooks, gatewayHook)
	case boil.BeforeUpsertHook:
		gatewayBeforeUpsertHooks = append(gatewayBeforeUpsertHooks, gatewayHook)
	case boil.AfterUpsertHook:
		gatewayAfterUpsertHooks = append(gatewayAfterUpsertHooks, gatewayHook)
	}
}

// One returns a single gateway record from the query.
func (q gatewayQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Gateway, error) {
	o := &Gateway{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for gateways")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Gateway records from the query.
func (q gatewayQuery) All(ctx context.Context, exec boil.ContextExecutor) (GatewaySlice, error) {
	var o []*Gateway

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Gateway slice")
	}

	if len(gatewayAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Gateway records in the query.
func (q gatewayQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count gateways rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q gatewayQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if gateways exists")
	}

	return count > 0, nil
}

// Orders retrieves all the order's Orders with an executor.
func (o *Gateway) Orders(mods ...qm.QueryMod) orderQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"orders\".\"gateway_id\"=?", o.ID),
	)

	return Orders(queryMods...)
}

// LoadOrders allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (gatewayL) LoadOrders(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGateway interface{}, mods queries.Applicator) error {
	var slice []*Gateway
	var object *Gateway

	if singular {
		var ok bool
		object, ok = maybeGateway.(*Gateway)
		if !ok {
			object = new(Gateway)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeGateway)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeGateway))
			}
		}
	} else {
		s, ok := maybeGateway.(*[]*Gateway)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeGateway)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeGateway))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gatewayR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gatewayR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`orders`),
		qm.WhereIn(`orders.gateway_id in ?`, args...),
		qmhelper.WhereIsNull(`orders.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load orders")
	}

	var resultSlice []*Order
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice orders")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on orders")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for orders")
	}

	if len(orderAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Orders = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &orderR{}
			}
			foreign.R.Gateway = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.GatewayID {
				local.R.Orders = append(local.R.Orders, foreign)
				if foreign.R == nil {
					foreign.R = &orderR{}
				}
				foreign.R.Gateway = local
				break
			}
		}
	}

	return nil
}

// AddOrders adds the given related objects to the existing relationships
// of the gateway, optionally inserting them as new records.
// Appends related to o.R.Orders.
// Sets related.R.Gateway appropriately.
func (o *Gateway) AddOrders(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Order) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.GatewayID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"orders\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"gateway_id"}),
				strmangle.WhereClause("\"", "\"", 2, orderPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.GatewayID = o.ID
		}
	}

	if o.R == nil {
		o.R = &gatewayR{
			Orders: related,
		}
	} else {
		o.R.Orders = append(o.R.Orders, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &orderR{
				Gateway: o,
			}
		} else {
			rel.R.Gateway = o
		}
	}
	return nil
}

// Gateways retrieves all the records using an executor.
func Gateways(mods ...qm.QueryMod) gatewayQuery {
	mods = append(mods, qm.From("\"gateways\""), qmhelper.WhereIsNull("\"gateways\".\"deleted_at\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"gateways\".*"})
	}

	return gatewayQuery{q}
}

// FindGateway retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGateway(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Gateway, error) {
	gatewayObj := &Gateway{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"gateways\" where \"id\"=$1 and \"deleted_at\" is null", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, gatewayObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from gateways")
	}

	if err = gatewayObj.doAfterSelectHooks(ctx, exec); err != nil {
		return gatewayObj, err
	}

	return gatewayObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Gateway) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no gateways provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(gatewayColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	gatewayInsertCacheMut.RLock()
	cache, cached := gatewayInsertCache[key]
	gatewayInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			gatewayAllColumns,
			gatewayColumnsWithDefault,
			gatewayColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(gatewayType, gatewayMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(gatewayType, gatewayMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"gateways\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"gateways\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into gateways")
	}

	if !cached {
		gatewayInsertCacheMut.Lock()
		gatewayInsertCache[key] = cache
		gatewayInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Gateway.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Gateway) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	gatewayUpdateCacheMut.RLock()
	cache, cached := gatewayUpdateCache[key]
	gatewayUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			gatewayAllColumns,
			gatewayPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update gateways, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"gateways\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, gatewayPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(gatewayType, gatewayMapping, append(wl, gatewayPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update gateways row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for gateways")
	}

	if !cached {
		gatewayUpdateCacheMut.Lock()
		gatewayUpdateCache[key] = cache
		gatewayUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q gatewayQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for gateways")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for gateways")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GatewaySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gatewayPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"gateways\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, gatewayPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in gateway slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all gateway")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Gateway) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no gateways provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(gatewayColumnsWithDefault, o)

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

	gatewayUpsertCacheMut.RLock()
	cache, cached := gatewayUpsertCache[key]
	gatewayUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			gatewayAllColumns,
			gatewayColumnsWithDefault,
			gatewayColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			gatewayAllColumns,
			gatewayPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert gateways, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(gatewayPrimaryKeyColumns))
			copy(conflict, gatewayPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"gateways\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(gatewayType, gatewayMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(gatewayType, gatewayMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert gateways")
	}

	if !cached {
		gatewayUpsertCacheMut.Lock()
		gatewayUpsertCache[key] = cache
		gatewayUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Gateway record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Gateway) Delete(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Gateway provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), gatewayPrimaryKeyMapping)
		sql = "DELETE FROM \"gateways\" WHERE \"id\"=$1"
	} else {
		currTime := time.Now().In(boil.GetLocation())
		o.DeletedAt = null.TimeFrom(currTime)
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"gateways\" SET %s WHERE \"id\"=$2",
			strmangle.SetParamNames("\"", "\"", 1, wl),
		)
		valueMapping, err := queries.BindMapping(gatewayType, gatewayMapping, append(wl, gatewayPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), valueMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from gateways")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for gateways")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q gatewayQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no gatewayQuery provided for delete all")
	}

	if hardDelete {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{"deleted_at": currTime})
	}

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from gateways")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for gateways")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GatewaySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(gatewayBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gatewayPrimaryKeyMapping)
			args = append(args, pkeyArgs...)
		}
		sql = "DELETE FROM \"gateways\" WHERE " +
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gatewayPrimaryKeyColumns, len(o))
	} else {
		currTime := time.Now().In(boil.GetLocation())
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gatewayPrimaryKeyMapping)
			args = append(args, pkeyArgs...)
			obj.DeletedAt = null.TimeFrom(currTime)
		}
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"gateways\" SET %s WHERE "+
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 2, gatewayPrimaryKeyColumns, len(o)),
			strmangle.SetParamNames("\"", "\"", 1, wl),
		)
		args = append([]interface{}{currTime}, args...)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from gateway slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for gateways")
	}

	if len(gatewayAfterDeleteHooks) != 0 {
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
func (o *Gateway) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGateway(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GatewaySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GatewaySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gatewayPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"gateways\".* FROM \"gateways\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gatewayPrimaryKeyColumns, len(*o)) +
		"and \"deleted_at\" is null"

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GatewaySlice")
	}

	*o = slice

	return nil
}

// GatewayExists checks if the Gateway row exists.
func GatewayExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"gateways\" where \"id\"=$1 and \"deleted_at\" is null limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if gateways exists")
	}

	return exists, nil
}

// Exists checks if the Gateway row exists.
func (o *Gateway) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return GatewayExists(ctx, exec, o.ID)
}
