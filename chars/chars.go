package chars

// https://github.com/angular/angular/blob/master/packages/compiler/src/chars.ts

var VEOF = 0
var VBSPACE = 8
var VTAB = 9
var VLF = 10
var VVTAB = 11
var VFF = 12
var VCR = 13
var VSPACE = 32
var VBANG = 33
var VDQ = 34
var VHASH = 35
var VDOLLAR = 36
var VPERCENT = 37
var VAMPERSAND = 38
var VSQ = 39
var VLPAREN = 40
var VRPAREN = 41
var VSTAR = 42
var VPLUS = 43
var VCOMMA = 44
var VMINUS = 45
var VPERIOD = 46
var VSLASH = 47
var VCOLON = 58
var VSEMICOLON = 59
var VLT = 60
var VEQ = 61
var VGT = 62
var VQUESTION = 63

var V0 = 48
var V7 = 55
var V9 = 57

var VA = 65
var VE = 69
var VF = 70
var VX = 88
var VZ = 90

var VLBRACKET = 91
var VBACKSLASH = 92
var VRBRACKET = 93
var VCARET = 94
var V_ = 95

var Va = 97
var Vb = 98
var Ve = 101
var Vf = 102
var Vn = 110
var Vr = 114
var Vt = 116
var Vu = 117
var Vv = 118
var Vx = 120
var Vz = 122

var VLBRACE = 123
var VBAR = 124
var VRBRACE = 125
var VNBSP = 160

var VPIPE = 124
var VTILDA = 126
var VAT = 64

var VBT = 96

func IsWhitespace(code int) bool {
	return (code >= VTAB && code <= VSPACE) || (code == VNBSP)
}

func IsDigit(code int) bool {
	return V0 <= code && code <= V9
}

func IsAsciiLetter(code int) bool {
	return code >= Va && code <= Vz || code >= VA && code <= VZ
}

func IsAsciiHexDigit(code int) bool {
	return code >= Va && code <= Vf || code >= VA && code <= VF || IsDigit(code)
}

func IsNewLine(code int) bool {
	return code == VLF || code == VCR
}

func IsOctalDigit(code int) bool {
	return V0 <= code && code <= V7
}

func IsQuote(code int) bool {
	return code == VSQ || code == VDQ || code == VBT
}
