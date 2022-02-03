package eqi

// Ignore is a list of ideographs which do not have 'true' visual equivalents
// in the EquivalentUnifiedIdeograph.txt dataset. By default these characters
// will be remain unprocessed. If you wish to also replace this characters,
// you can set the IgnoreVisuallyUnique flag to false.
var Unique = map[string]bool{
	//"⺈": true, // CJK RADICAL KNIFE ONE
	"⺌": true, // [2] CJK RADICAL SMALL ONE
	"⺍": true, // [2] CJK RADICAL SMALL TWO
	"⺗": true, // CJK RADICAL HEART ONE
	"⺜": true, // CJK RADICAL SUN
	"⺝": true, // CJK RADICAL MOON
	"⺤": true, // [2] CJK RADICAL PAW ONE
	"⺧": true, // CJK RADICAL COW
	"⺩": true, // CJK RADICAL JADE
	"⺫": true, // CJK RADICAL EYE
	"⺻": true, // CJK RADICAL BRUSH TWO
	"⺼": true, // CJK RADICAL MEAT
	"⻏": true, // CJK RADICAL CITY
	"⻗": true, // CJK RADICAL RAIN
}
