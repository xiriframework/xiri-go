package table

// alignPtr returns a pointer to a FieldAlign value.
func alignPtr(a FieldAlign) *FieldAlign { return &a }

// fieldDefault holds the non-formatter defaults for a FieldTypeHint.
type fieldDefault struct {
	fieldType FieldType
	align     *FieldAlign
	decimals  int
	search    bool
	sort      bool
	csv       bool
}

// fieldDefaults maps each FieldTypeHint to its default configuration.
// Formatters are assigned separately because they are closures that may reference field values.
var fieldDefaults = map[FieldTypeHint]fieldDefault{
	Id:              {FieldTypeID, alignPtr(FieldAlignRight), 0, true, true, false},
	Integer:         {FieldTypeNumber, alignPtr(FieldAlignRight), 0, true, true, true},
	Float:           {FieldTypeNumber, alignPtr(FieldAlignRight), 2, true, true, true},
	Text:            {FieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	Bool:            {FieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	DateTime:        {FieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	Date:            {FieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	Distance:        {FieldTypeNumber, alignPtr(FieldAlignRight), 2, true, true, true},
	Pressure:        {FieldTypeNumber, alignPtr(FieldAlignRight), 2, true, true, true},
	Speed:           {FieldTypeNumber, alignPtr(FieldAlignRight), 1, true, true, true},
	Buttons:         {FieldTypeButtons, alignPtr(FieldAlignCenter), 0, false, false, false},
	Icon:            {FieldTypeIcon, alignPtr(FieldAlignCenter), 0, false, true, true},
	Link:            {FieldTypeLink, alignPtr(FieldAlignLeft), 0, true, true, true},
	Html:            {FieldTypeHtml, alignPtr(FieldAlignLeft), 0, true, true, false},
	Input:           {FieldTypeInput, alignPtr(FieldAlignLeft), 0, false, true, false},
	Text2:           {FieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	Text2Int:        {FieldTypeText2, alignPtr(FieldAlignRight), 0, true, true, true},
	Text2Float:      {FieldTypeText2, alignPtr(FieldAlignRight), 2, true, true, true},
	Text2DateTime:   {FieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	Text2Date:       {FieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	Text2Distance:   {FieldTypeText2, alignPtr(FieldAlignRight), 2, true, true, true},
	Text2Speed:      {FieldTypeText2, alignPtr(FieldAlignRight), 1, true, true, true},
	Text2Bool:       {FieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	TimeLength:      {FieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	Text2TimeLength: {FieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	Header:          {FieldTypeHeader, alignPtr(FieldAlignLeft), 0, false, false, false},
}

// applyFieldTypeDefaults configures a field builder with appropriate defaults
// based on the specified field type hint.
func applyFieldTypeDefaults[T any](builder *FieldBuilder[T], fieldType FieldTypeHint) *FieldBuilder[T] {
	builder.field.fieldTypeHint = fieldType

	// Apply struct defaults from the map
	def, ok := fieldDefaults[fieldType]
	if !ok {
		def = fieldDefault{FieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true}
	}

	builder.field.fieldType = def.fieldType
	builder.field.align = def.align
	builder.field.decimals = def.decimals
	builder.field.search = def.search
	builder.field.sort = def.sort
	builder.field.csv = def.csv

	// Bool-specific defaults
	if fieldType == Bool {
		builder.field.boolTrueText = "true"
		builder.field.boolFalseText = "false"
	}

	// Assign formatter (closures that may reference decimals)
	switch fieldType {
	case Id:
		builder.field.defaultFormatter = createIdFormatter()
	case Integer:
		builder.field.defaultFormatter = createIntegerFormatter()
	case Float:
		builder.field.defaultFormatter = createFloatFormatter(def.decimals)
	case Text, Html, Input, Header:
		builder.field.defaultFormatter = createTextFormatter()
	case Bool:
		builder.field.defaultFormatter = createBoolFormatter("true", "false")
	case DateTime:
		builder.field.defaultFormatter = createDateTimeFormatter()
	case Date:
		builder.field.defaultFormatter = createDateFormatter()
	case Distance:
		builder.field.defaultFormatter = createDistanceFormatter(def.decimals)
	case Pressure:
		builder.field.defaultFormatter = createPressureFormatter(def.decimals)
	case Speed:
		builder.field.defaultFormatter = createSpeedFormatter(def.decimals)
	case Buttons:
		builder.field.defaultFormatter = createPassthroughFormatter()
	case Icon:
		builder.field.defaultFormatter = createTextFormatter()
	case Link:
		builder.field.defaultFormatter = createLinkFormatter()
	case Text2:
		builder.field.defaultFormatter = createText2Formatter()
	case Text2Int:
		builder.field.defaultFormatter = createText2IntFormatter()
	case Text2Float:
		builder.field.defaultFormatter = createText2FloatFormatter(def.decimals)
	case Text2DateTime:
		builder.field.defaultFormatter = createText2DateTimeFormatter()
	case Text2Date:
		builder.field.defaultFormatter = createText2DateFormatter()
	case Text2Distance:
		builder.field.defaultFormatter = createText2DistanceFormatter(def.decimals)
	case Text2Speed:
		builder.field.defaultFormatter = createText2SpeedFormatter(def.decimals)
	case Text2Bool:
		builder.field.defaultFormatter = createText2BoolFormatter()
	case TimeLength:
		builder.field.defaultFormatter = createTimeLengthFormatter()
	case Text2TimeLength:
		builder.field.defaultFormatter = createText2TimeLengthFormatter()
	default:
		builder.field.defaultFormatter = createTextFormatter()
	}

	return builder
}
