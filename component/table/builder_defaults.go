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

// applyFieldTypeDefaults configures a fieldBase with appropriate defaults
// based on the specified field type hint.
// This function is non-generic to avoid monomorphization across all Table[T] instantiations.
func applyFieldTypeDefaults(base *fieldBase, ft fieldTypeHint) {
	base.fieldTypeHint = ft

	// Apply struct defaults from the map
	def, ok := fieldDefaults[ft]
	if !ok {
		def = fieldDefault{fieldTypeText, alignPtr(FieldAlignLeft), 0, true, true, true}
	}

	base.fieldType = def.fieldType
	base.align = def.align
	base.decimals = def.decimals
	base.search = def.search
	base.sort = def.sort
	base.csv = def.csv

	// Bool-specific defaults
	if ft == boolean {
		base.boolTrueText = "true"
		base.boolFalseText = "false"
	}

	// Assign formatter (closures that may reference decimals)
	switch ft {
	case idHint:
		base.defaultFormatter = createIdFormatter()
	case integer:
		base.defaultFormatter = createIntegerFormatter()
	case float:
		base.defaultFormatter = createFloatFormatter(def.decimals)
	case text, html, inputHint, header:
		base.defaultFormatter = createTextFormatter()
	case boolean:
		base.defaultFormatter = createBoolFormatter("true", "false")
	case dateTime:
		base.defaultFormatter = createDateTimeFormatter()
	case date:
		base.defaultFormatter = createDateFormatter()
	case distanceHint:
		base.defaultFormatter = createDistanceFormatter(def.decimals)
	case pressureHint:
		base.defaultFormatter = createPressureFormatter(def.decimals)
	case speedHint:
		base.defaultFormatter = createSpeedFormatter(def.decimals)
	case buttons:
		base.defaultFormatter = createPassthroughFormatter()
	case icon:
		base.defaultFormatter = createTextFormatter()
	case link:
		base.defaultFormatter = createLinkFormatter()
	case text2:
		base.defaultFormatter = createText2Formatter()
	case text2Int:
		base.defaultFormatter = createText2IntFormatter()
	case text2Float:
		base.defaultFormatter = createText2FloatFormatter(def.decimals)
	case text2DateTime:
		base.defaultFormatter = createText2DateTimeFormatter()
	case text2Date:
		base.defaultFormatter = createText2DateFormatter()
	case text2Distance:
		base.defaultFormatter = createText2DistanceFormatter(def.decimals)
	case text2Speed:
		base.defaultFormatter = createText2SpeedFormatter(def.decimals)
	case text2Bool:
		base.defaultFormatter = createText2BoolFormatter()
	case timeLength:
		base.defaultFormatter = createTimeLengthFormatter()
	case text2TimeLength:
		base.defaultFormatter = createText2TimeLengthFormatter()
	case textN:
		base.defaultFormatter = createTextNFormatter()
	case integerN:
		base.defaultFormatter = createIntegerNFormatter()
	case floatN:
		base.defaultFormatter = createFloatNFormatter(def.decimals)
	case dateTimeN:
		base.defaultFormatter = createDateTimeNFormatter()
	case dateN:
		base.defaultFormatter = createDateNFormatter()
	case distanceN:
		base.defaultFormatter = createDistanceNFormatter(def.decimals)
	case speedN:
		base.defaultFormatter = createSpeedNFormatter(def.decimals)
	case boolN:
		base.defaultFormatter = createBoolNFormatter()
	case timeLengthN:
		base.defaultFormatter = createTimeLengthNFormatter()
	default:
		base.defaultFormatter = createTextFormatter()
	}
}
