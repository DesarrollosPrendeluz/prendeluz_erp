package models

import "time"

type ItemType string

const (
	Father ItemType = "father"
	Son    ItemType = "son"
)

type Item struct {
	ID                       uint64   `gorm:"primaryKey;autoIncrement"`
	MainSKU                  string   `gorm:"size:255;not null;uniqueIndex"`
	EAN                      string   `gorm:"size:255;not null;index"`
	Name                     *string  `gorm:"size:255"`
	ItemType                 ItemType `gorm:"type:enum('father','son');not null;default:'father'"`
	AssignmentCost           float64  `gorm:"not null"`
	Description              *string  `gorm:"type:text"`
	Price                    *float64
	UnitaryBFO               *float64
	CCBFO                    *float64 `gorm:"column:cc_bfo"`
	PerPurchaseBFO           *float64
	AmazonBFO                *float64
	New                      bool `gorm:"not null;default:true"`
	CreatedAt                *time.Time
	UpdatedAt                *time.Time
	Discontinued             bool    `gorm:"not null;default:false"`
	EnergyCertification      *string `gorm:"size:255"`
	EnergyCertificationImage *string `gorm:"size:255"`
	BatteryType              *string `gorm:"size:255"`
	BatteryNumber            *string `gorm:"size:255"`
	IsSincroToPymeSQL        bool    `gorm:"not null;default:false"`
	Normative                *string `gorm:"type:text"`
	Category                 *uint64 `gorm:"index"`
	AmazonCategory           *uint64 `gorm:"index"`
	AmazonSubcategory        *uint64 `gorm:"index"`
	Weight                   *string `gorm:"size:255"`
	High                     *string `gorm:"size:255"`
	Wide                     *string `gorm:"size:255"`
	Long                     *string `gorm:"size:255"`
	SafetyData               *string `gorm:"size:255"`
	InstructionManual        *string `gorm:"size:255"`
	ImagesStatus             bool    `gorm:"not null;default:false"`
	CEMarking                bool    `gorm:"not null;default:false"`
	RohsMarking              bool    `gorm:"not null;default:false"`
	DeclarationOfConformity  string  `gorm:"size:255;not null;default:'0'"`
	CountryProhibitions      string  `gorm:"size:255;not null;default:'0'"`

	//Category          *Category       `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT;"`
	//AmazonCategory    *Category       `gorm:"foreignKey:AmazonCategoryID;constraint:OnDelete:RESTRICT;"`
	//AmazonSubcategory *Category       `gorm:"foreignKey:AmazonSubcategoryID;constraint:OnDelete:RESTRICT;"`
	ChildRel      *ItemsParents   `gorm:"foreignKey:ID;references:ParentItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	FatherRel     *ItemsParents   `gorm:"foreignKey:ID;references:ChildItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AsinRel       *Asin           `gorm:"foreignKey:ID;references:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SupplierItems *[]SupplierItem `gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
