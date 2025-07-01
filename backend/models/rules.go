package models

type HOARules struct {
	ID            string `datastore:"id"` // single row w/ known ID like "default"
	Renovation    string `datastore:"renovation"`
	PaintColors   string `datastore:"paint_colors"`
	ShingleTypes  string `datastore:"shingle_types"`
	GeneralBylaws string `datastore:"general_bylaws"`
	LastUpdatedBy string `datastore:"last_updated_by"`
}
