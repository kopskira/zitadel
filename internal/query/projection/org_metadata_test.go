package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestOrgMetadataProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduceMetadataSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MetadataSetType),
					org.AggregateType,
					[]byte(`{
						"key": "key",
						"value": "dmFsdWU="
					}`),
				), org.MetadataSetEventMapper),
			},
			reduce: (&orgMetadataProjection{}).reduceMetadataSet,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMetadataProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.org_metadata (instance_id, org_id, key, resource_owner, creation_date, change_date, sequence, value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (instance_id, org_id, key) DO UPDATE SET (resource_owner, creation_date, change_date, sequence, value) = (EXCLUDED.resource_owner, EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.value)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
								"key",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								[]byte("value"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMetadataRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MetadataRemovedType),
					org.AggregateType,
					[]byte(`{
						"key": "key"
					}`),
				), org.MetadataRemovedEventMapper),
			},
			reduce: (&orgMetadataProjection{}).reduceMetadataRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMetadataProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.org_metadata WHERE (org_id = $1) AND (key = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"key",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMetadataRemovedAll",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MetadataRemovedAllType),
					org.AggregateType,
					nil,
				), org.MetadataRemovedAllEventMapper),
			},
			reduce: (&orgMetadataProjection{}).reduceMetadataRemovedAll,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMetadataProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.org_metadata WHERE (org_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMetadataRemovedAll (org removed)",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					nil,
				), org.OrgRemovedEventMapper),
			},
			reduce: (&orgMetadataProjection{}).reduceMetadataRemovedAll,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMetadataProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.org_metadata WHERE (org_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
