package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SnapshotResponse struct {
	ID           string `json:"id"`
	CreationDate string `json:"creation_date"`
	DeletionDate string `json:"deletion_date"`
}

func NewSnapshotObject(snapshot SnapshotResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"id":            types.StringValue(snapshot.ID),
		"creation_date": types.StringValue(snapshot.CreationDate),
		"deletion_date": types.StringValue(snapshot.DeletionDate),
	}
	return types.ObjectValue(SnapshotObjectType().AttrTypes, attrs)
}

func SnapshotObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":            types.StringType,
			"creation_date": types.StringType,
			"deletion_date": types.StringType,
		},
	}
}

func SnapshotDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Snapshot identifier",
		},
		"creation_date": schema.StringAttribute{
			Computed:    true,
			Description: "Snapshot creation date",
		},
		"deletion_date": schema.StringAttribute{
			Computed:    true,
			Description: "Snapshot deletion date",
		},
	}
}

func SnapshotResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed:    true,
			Description: "Snapshot identifier",
		},
		"creation_date": rschema.StringAttribute{
			Computed:    true,
			Description: "Snapshot creation date",
		},
		"deletion_date": rschema.StringAttribute{
			Computed:    true,
			Description: "Snapshot deletion date",
		},
	}
}
