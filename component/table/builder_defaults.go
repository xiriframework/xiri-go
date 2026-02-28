package table

// alignPtr returns a pointer to a FieldAlign value.
func alignPtr(a FieldAlign) *FieldAlign { return &a }

// fieldDefault holds the non-formatter defaults for a fieldTypeHint.
type fieldDefault struct {
	fieldType fieldType
	align     *FieldAlign
	decimals  int
	search    bool
	sort      bool
	csv       bool
}

// fieldDefaults maps each fieldTypeHint to its default configuration.
// Formatters are assigned separately because they are closures that may reference field values.
var fieldDefaults = map[fieldTypeHint]fieldDefault{
	idHint:          {fieldTypeID, alignPtr(FieldAlignRight), 0, true, true, false},
	integer:         {fieldTypeNumber, alignPtr(FieldAlignRight), 0, true, true, true},
	float:           {fieldTypeNumber, alignPtr(FieldAlignRight), 2, true, true, true},
	text:            {fieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	boolean:         {fieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	dateTime:        {fieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	date:            {fieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true},
	distanceHint:    {fieldTypeNumber, alignPtr(FieldAlignRight), 2, true, true, true},
	pressureHint:    {fieldTypeNumber, alignPtr(FieldAlignRight), 2, true, true, true},
	speedHint:       {fieldTypeNumber, alignPtr(FieldAlignRight), 1, true, true, true},
	buttons:         {fieldTypeButtons, alignPtr(FieldAlignCenter), 0, false, false, false},
	icon:            {fieldTypeIcon, alignPtr(FieldAlignCenter), 0, false, true, true},
	link:            {fieldTypeLink, alignPtr(FieldAlignLeft), 0, true, true, true},
	html:            {fieldTypeHtml, alignPtr(FieldAlignLeft), 0, true, true, false},
	inputHint:       {fieldTypeInput, alignPtr(FieldAlignLeft), 0, false, true, false},
	text2:           {fieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	text2Int:        {fieldTypeText2, alignPtr(FieldAlignRight), 0, true, true, true},
	text2Float:      {fieldTypeText2, alignPtr(FieldAlignRight), 2, true, true, true},
	text2DateTime:   {fieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	text2Date:       {fieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	text2Distance:   {fieldTypeText2, alignPtr(FieldAlignRight), 2, true, true, true},
	text2Speed:      {fieldTypeText2, alignPtr(FieldAlignRight), 1, true, true, true},
	text2Bool:       {fieldTypeText2, alignPtr(FieldAlignLeft), 0, true, true, true},
	timeLength:      {fieldTypeText, alignPtr(FieldAlignRight), 0, true, true, true},
	text2TimeLength: {fieldTypeText2, alignPtr(FieldAlignRight), 0, true, true, true},
	textN:           {fieldTypeTextN, alignPtr(FieldAlignLeft), 0, true, true, true},
	integerN:        {fieldTypeTextN, alignPtr(FieldAlignRight), 0, true, true, true},
	floatN:          {fieldTypeTextN, alignPtr(FieldAlignRight), 2, true, true, true},
	dateTimeN:       {fieldTypeTextN, alignPtr(FieldAlignLeft), 0, true, true, true},
	dateN:           {fieldTypeTextN, alignPtr(FieldAlignLeft), 0, true, true, true},
	distanceN:       {fieldTypeTextN, alignPtr(FieldAlignRight), 2, true, true, true},
	speedN:          {fieldTypeTextN, alignPtr(FieldAlignRight), 1, true, true, true},
	boolN:           {fieldTypeTextN, alignPtr(FieldAlignLeft), 0, true, true, true},
	timeLengthN:     {fieldTypeTextN, alignPtr(FieldAlignRight), 0, true, true, true},
	header:          {fieldTypeHeader, alignPtr(FieldAlignLeft), 0, false, false, false},
}

// applyFieldTypeDefaults configures a field builder with appropriate defaults
// based on the specified field type hint.
func applyFieldTypeDefaults[T any](builder *FieldBuilder[T], ft fieldTypeHint) *FieldBuilder[T] {
	builder.field.fieldTypeHint = ft

	// Apply struct defaults from the map
	def, ok := fieldDefaults[ft]
	if !ok {
		def = fieldDefault{fieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true}
	}

	builder.field.fieldType = def.fieldType
	builder.field.align = def.align
	builder.field.decimals = def.decimals
	builder.field.search = def.search
	builder.field.sort = def.sort
	builder.field.csv = def.csv

	// Bool-specific defaults
	if ft == boolean {
		builder.field.boolTrueText = "true"
		builder.field.boolFalseText = "false"
	}

	// Assign formatter (closures that may reference decimals)
	switch ft {
	case idHint:
		builder.field.defaultFormatter = createIdFormatter()
	case integer:
		builder.field.defaultFormatter = createIntegerFormatter()
	case float:
		builder.field.defaultFormatter = createFloatFormatter(def.decimals)
	case text, html, inputHint, header:
		builder.field.defaultFormatter = createTextFormatter()
	case boolean:
		builder.field.defaultFormatter = createBoolFormatter("true", "false")
	case dateTime:
		builder.field.defaultFormatter = createDateTimeFormatter()
	case date:
		builder.field.defaultFormatter = createDateFormatter()
	case distanceHint:
		builder.field.defaultFormatter = createDistanceFormatter(def.decimals)
	case pressureHint:
		builder.field.defaultFormatter = createPressureFormatter(def.decimals)
	case speedHint:
		builder.field.defaultFormatter = createSpeedFormatter(def.decimals)
	case buttons:
		builder.field.defaultFormatter = createPassthroughFormatter()
	case icon:
		builder.field.defaultFormatter = createTextFormatter()
	case link:
		builder.field.defaultFormatter = createLinkFormatter()
	case text2:
		builder.field.defaultFormatter = createText2Formatter()
	case text2Int:
		builder.field.defaultFormatter = createText2IntFormatter()
	case text2Float:
		builder.field.defaultFormatter = createText2FloatFormatter(def.decimals)
	case text2DateTime:
		builder.field.defaultFormatter = createText2DateTimeFormatter()
	case text2Date:
		builder.field.defaultFormatter = createText2DateFormatter()
	case text2Distance:
		builder.field.defaultFormatter = createText2DistanceFormatter(def.decimals)
	case text2Speed:
		builder.field.defaultFormatter = createText2SpeedFormatter(def.decimals)
	case text2Bool:
		builder.field.defaultFormatter = createText2BoolFormatter()
	case timeLength:
		builder.field.defaultFormatter = createTimeLengthFormatter()
	case text2TimeLength:
		builder.field.defaultFormatter = createText2TimeLengthFormatter()
	case textN:
		builder.field.defaultFormatter = createTextNFormatter()
	case integerN:
		builder.field.defaultFormatter = createIntegerNFormatter()
	case floatN:
		builder.field.defaultFormatter = createFloatNFormatter(def.decimals)
	case dateTimeN:
		builder.field.defaultFormatter = createDateTimeNFormatter()
	case dateN:
		builder.field.defaultFormatter = createDateNFormatter()
	case distanceN:
		builder.field.defaultFormatter = createDistanceNFormatter(def.decimals)
	case speedN:
		builder.field.defaultFormatter = createSpeedNFormatter(def.decimals)
	case boolN:
		builder.field.defaultFormatter = createBoolNFormatter()
	case timeLengthN:
		builder.field.defaultFormatter = createTimeLengthNFormatter()
	default:
		builder.field.defaultFormatter = createTextFormatter()
	}

	return builder
}
