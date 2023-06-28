package aud

import (
	"github.com/invopop/validation"

	"github.com/infragmo/auditum/internal/util/validate"
)

type Settings struct {
	Records RecordsSettings `yaml:"records" json:"records"`
}

func (s Settings) Validate() error {
	return validate.Each(
		s.Records,
	)
}

type RecordsSettings struct {
	UpdateEnabled bool                `yaml:"updateEnabled" json:"updateEnabled"`
	DeleteEnabled bool                `yaml:"deleteEnabled" json:"deleteEnabled"`
	Restrictions  RecordsRestrictions `yaml:"restrictions" json:"restrictions"`
}

func (r RecordsSettings) Validate() error {
	return validate.Each(
		r.Restrictions,
	)
}

type RecordsRestrictions struct {
	Labels    RestrictionsKeyValue         `yaml:"labels" json:"labels"`
	Resource  RecordsRestrictionsResource  `yaml:"resource" json:"resource"`
	Operation RecordsRestrictionsOperation `yaml:"operation" json:"operation"`
	Actor     RecordsRestrictionsActor     `yaml:"actor" json:"actor"`
}

func (r RecordsRestrictions) Validate() error {
	return validate.Each(
		r.Labels,
		r.Resource,
		r.Operation,
		r.Actor,
	)
}

type RecordsRestrictionsResource struct {
	Type     RestrictionsString                 `yaml:"type" json:"type"`
	ID       RestrictionsString                 `yaml:"id" json:"id"`
	Metadata RestrictionsKeyValue               `yaml:"metadata" json:"metadata"`
	Changes  RecordsRestrictionsResourceChanges `yaml:"changes" json:"changes"`
}

func (r RecordsRestrictionsResource) Validate() error {
	return validate.Each(
		r.Type,
		r.ID,
		r.Metadata,
		r.Changes,
	)
}

type RecordsRestrictionsResourceChanges struct {
	TotalMaxCount int                `yaml:"totalMaxCount" json:"totalMaxCount"`
	Name          RestrictionsString `yaml:"name" json:"name"`
	Description   RestrictionsString `yaml:"description" json:"description"`
	OldValue      RestrictionsBytes  `yaml:"oldValue" json:"oldValue"`
	NewValue      RestrictionsBytes  `yaml:"newValue" json:"newValue"`
}

func (r RecordsRestrictionsResourceChanges) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(
			&r.TotalMaxCount,
			validation.Min(0),
		),
	)

	return validate.Each(
		validate.ErrorAsValidatable(err),
		r.Name,
		r.Description,
		r.OldValue,
		r.NewValue,
	)
}

type RecordsRestrictionsOperation struct {
	Type     RestrictionsString   `yaml:"type" json:"type"`
	ID       RestrictionsString   `yaml:"id" json:"id"`
	Metadata RestrictionsKeyValue `yaml:"metadata" json:"metadata"`
}

func (r RecordsRestrictionsOperation) Validate() error {
	return validate.Each(
		r.Type,
		r.ID,
		r.Metadata,
	)
}

type RecordsRestrictionsActor struct {
	Type     RestrictionsString   `yaml:"type" json:"type"`
	ID       RestrictionsString   `yaml:"id" json:"id"`
	Metadata RestrictionsKeyValue `yaml:"metadata" json:"metadata"`
}

func (r RecordsRestrictionsActor) Validate() error {
	return validate.Each(
		r.Type,
		r.ID,
		r.Metadata,
	)
}

type RestrictionsKeyValue struct {
	KeyMaxSizeBytes   int `yaml:"keyMaxSizeBytes" json:"keyMaxSizeBytes"`
	ValueMaxSizeBytes int `yaml:"valueMaxSizeBytes" json:"valueMaxSizeBytes"`
	TotalMaxSizeBytes int `yaml:"totalMaxSizeBytes" json:"totalMaxSizeBytes"`
}

func (r RestrictionsKeyValue) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.KeyMaxSizeBytes,
			validation.Min(0),
		),
		validation.Field(
			&r.ValueMaxSizeBytes,
			validation.Min(0),
		),
		validation.Field(
			&r.TotalMaxSizeBytes,
			validation.Min(0),
		),
	)
}

type RestrictionsString struct {
	MaxSizeBytes int `yaml:"maxSizeBytes" json:"maxSizeBytes"`
}

func (r RestrictionsString) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.MaxSizeBytes,
			validation.Min(0),
		),
	)
}

type RestrictionsBytes struct {
	MaxSizeBytes int `yaml:"maxSizeBytes" json:"maxSizeBytes"`
}

func (r RestrictionsBytes) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.MaxSizeBytes,
			validation.Min(0),
		),
	)
}

var DefaultSettings = Settings{
	Records: RecordsSettings{
		UpdateEnabled: false,
		DeleteEnabled: false,
		Restrictions: RecordsRestrictions{
			Labels: RestrictionsKeyValue{
				KeyMaxSizeBytes:   64,
				ValueMaxSizeBytes: 256,
				TotalMaxSizeBytes: 2048,
			},
			Resource: RecordsRestrictionsResource{
				Type: RestrictionsString{
					MaxSizeBytes: 256,
				},
				ID: RestrictionsString{
					MaxSizeBytes: 256,
				},
				Metadata: RestrictionsKeyValue{
					KeyMaxSizeBytes:   64,
					ValueMaxSizeBytes: 256,
					TotalMaxSizeBytes: 2048,
				},
				Changes: RecordsRestrictionsResourceChanges{
					TotalMaxCount: 20,
					Name: RestrictionsString{
						MaxSizeBytes: 256,
					},
					Description: RestrictionsString{
						MaxSizeBytes: 1024,
					},
					OldValue: RestrictionsBytes{
						MaxSizeBytes: 4096,
					},
					NewValue: RestrictionsBytes{
						MaxSizeBytes: 4096,
					},
				},
			},
			Operation: RecordsRestrictionsOperation{
				Type: RestrictionsString{
					MaxSizeBytes: 256,
				},
				ID: RestrictionsString{
					MaxSizeBytes: 512,
				},
				Metadata: RestrictionsKeyValue{
					KeyMaxSizeBytes:   64,
					ValueMaxSizeBytes: 256,
					TotalMaxSizeBytes: 2048,
				},
			},
			Actor: RecordsRestrictionsActor{
				Type: RestrictionsString{
					MaxSizeBytes: 256,
				},
				ID: RestrictionsString{
					MaxSizeBytes: 256,
				},
				Metadata: RestrictionsKeyValue{
					KeyMaxSizeBytes:   64,
					ValueMaxSizeBytes: 256,
					TotalMaxSizeBytes: 2048,
				},
			},
		},
	},
}
