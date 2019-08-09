// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wfmt_test

import (
	"math"
	"reflect"
	"strings"
	"testing"

	. "github.com/lostsnow/wfmt"
)

var (
	NaN    = math.NaN()
	posInf = math.Inf(1)
	negInf = math.Inf(-1)

	intVar = 0
)

var fmtTests = []struct {
	fmt string
	val interface{}
	out string
}{
	{"%10q", "⌘", `       "⌘"`},
	{"%+10q", "⌘", `  "\u2318"`},
	{"%-10q", "⌘", `"⌘"       `},
	{"%+-10q", "⌘", `"\u2318"  `},
	{"%010q", "⌘", `0000000"⌘"`},
	{"%+010q", "⌘", `00"\u2318"`},
	{"%-010q", "⌘", `"⌘"       `}, // 0 has no effect when - is present.
	{"%+-010q", "⌘", `"\u2318"  `},
	{"%#8q", "\n", `    "\n"`},
	{"%#+8q", "\r", `    "\r"`},
	{"%#-8q", "\t", "`	`     "},
	{"%#+-8q", "\b", `"\b"    `},

	{"%.0c", '⌘', "⌘"}, // Specifying precision should have no effect.
	{"%3c", '⌘', "  ⌘"},
	{"%-3c", '⌘', "⌘  "},

	{"%.0q", '☺', `'☺'`}, // Specifying precision should have no effect.
	{"%10q", '⌘', `       '⌘'`},
	{"%+10q", '⌘', `  '\u2318'`},
	{"%-10q", '⌘', `'⌘'       `},
	{"%+-10q", '⌘', `'\u2318'  `},
	{"%010q", '⌘', `0000000'⌘'`},
	{"%+010q", '⌘', `00'\u2318'`},
	{"%-010q", '⌘', `'⌘'       `}, // 0 has no effect when - is present.
	{"%+-010q", '⌘', `'\u2318'  `},

	// width
	{"%5s", "啊啊", " 啊啊"},
	{"%5s", "abc", "  abc"},
	{"%2s", "\u263a", "☺"},
	{"%-5s", "abc", "abc  "},
	{"%-8q", "abc", `"abc"   `},
	{"%05s", "abc", "00abc"},
	{"%08q", "abc", `000"abc"`},
	{"%5s", "abcdefghijklmnopqrstuvwxyz", "abcdefghijklmnopqrstuvwxyz"},
	{"%.5s", "abcdefghijklmnopqrstuvwxyz", "abcde"},
	{"%.0s", "日本語日本語", ""},
	{"%.4s", "日本語日本語", "日本語日"},
	{"%.10s", "日本語日本語", "日本語日本語"},
	{"%.5s", []byte("日本語日本語"), "日本語日本"},
	{"%.5q", "abcdefghijklmnopqrstuvwxyz", `"abcde"`},
	{"%.5x", "abcdefghijklmnopqrstuvwxyz", "6162636465"},
	{"%.5q", []byte("abcdefghijklmnopqrstuvwxyz"), `"abcde"`},
	{"%.5x", []byte("abcdefghijklmnopqrstuvwxyz"), "6162636465"},
	{"%.3q", "日本語日本語", `"日本語"`},
	{"%.3q", []byte("日本語日本語"), `"日本語"`},
	{"%.1q", "日本語", `"日"`},
	{"%.1q", []byte("日本語"), `"日"`},
	{"%.1x", "日本語", "e6"},
	{"%.1X", []byte("日本語"), "E6"},
	{"%10.1q", "日本語日本語", `      "日"`},
	{"%10v", nil, "     <nil>"},
	{"%-10v", nil, "<nil>     "},

	// unicode format
	{"%#.2U", 'x', `U+0078 'x'`}, // Precisions below 4 should print 4 digits.
	{"%#14.6U", '⌘', "  U+002318 '⌘'"},
	{"%#-14.6U", '⌘', "U+002318 '⌘'  "},
	{"%#014.6U", '⌘', "  U+002318 '⌘'"},
	{"%#-014.6U", '⌘', "U+002318 '⌘'  "},
	{"%.68U", uint(42), zeroFill("U+", 68, "2A")},
	{"%#.68U", '日', zeroFill("U+", 68, "65E5") + " '日'"},

	// floats
	{"%+.3e", 0.0, "+0.000e+00"},
	{"%+.3e", 1.0, "+1.000e+00"},
	{"%+.3f", -1.0, "-1.000"},
	{"%+.3F", -1.0, "-1.000"},
	{"%+.3F", float32(-1.0), "-1.000"},
	{"%+07.2f", 1.0, "+001.00"},
	{"%+07.2f", -1.0, "-001.00"},
	{"%-07.2f", 1.0, "1.00   "},
	{"%-07.2f", -1.0, "-1.00  "},
	{"%+-07.2f", 1.0, "+1.00  "},
	{"%+-07.2f", -1.0, "-1.00  "},
	{"%-+07.2f", 1.0, "+1.00  "},
	{"%-+07.2f", -1.0, "-1.00  "},
	{"%+10.2f", +1.0, "     +1.00"},
	{"%+10.2f", -1.0, "     -1.00"},
	{"% .3E", -1.0, "-1.000E+00"},
	{"% .3e", 1.0, " 1.000e+00"},
	{"%+.3g", 0.0, "+0"},
	{"%+.3g", 1.0, "+1"},
	{"%+.3g", -1.0, "-1"},
	{"% .3g", -1.0, "-1"},
	{"% .3g", 1.0, " 1"},
	// Test sharp flag used with floats.
	{"%#.0f", 1.0, "1."},
	{"%#.0e", 1.0, "1.e+00"},
	{"%#.0g", 1.0, "1."},
	{"%#.0g", 1100000.0, "1.e+06"},
	{"%#.4f", 1.0, "1.0000"},
	{"%#.4e", 1.0, "1.0000e+00"},
	{"%#.4g", 1.0, "1.000"},
	{"%#.4g", 100000.0, "1.000e+05"},
	{"%#.0f", 123.0, "123."},
	{"%#.0e", 123.0, "1.e+02"},
	{"%#.0g", 123.0, "1.e+02"},
	{"%#.4f", 123.0, "123.0000"},
	{"%#.4e", 123.0, "1.2300e+02"},
	{"%#.4g", 123.0, "123.0"},
	{"%#.4g", 123000.0, "1.230e+05"},
	{"%#9.4g", 1.0, "    1.000"},
	// Precision has no effect for binary float format.
	{"%.4b", float32(1.0), "8388608p-23"},
	{"%.4b", -1.0, "-4503599627370496p-52"},
	// Test correct f.intbuf boundary checks.
	{"%.68f", 1.0, zeroFill("1.", 68, "")},
	{"%.68f", -1.0, zeroFill("-1.", 68, "")},
	// float infinites and NaNs
	{"%.1f", negInf, "-Inf"},
	{"%20f", posInf, "                +Inf"},
	{"% 20F", posInf, "                 Inf"},
	{"% 20e", negInf, "                -Inf"},
	{"%+20E", negInf, "                -Inf"},
	{"% +20g", negInf, "                -Inf"},
	{"%+-20G", posInf, "+Inf                "},
	{"%20e", NaN, "                 NaN"},
	{"% +20E", NaN, "                +NaN"},
	{"% -20g", NaN, " NaN                "},
	{"%+-20G", NaN, "+NaN                "},
	// Zero padding does not apply to infinities and NaN.
	{"%+020e", posInf, "                +Inf"},
	{"%-020f", negInf, "-Inf                "},
	{"%-020E", NaN, "NaN                 "},

	// complex values
	{"%+.3e", 0i, "(+0.000e+00+0.000e+00i)"},
	{"%+.3f", 0i, "(+0.000+0.000i)"},
	{"%+.3g", 0i, "(+0+0i)"},
	{"%+.3e", 1 + 2i, "(+1.000e+00+2.000e+00i)"},
	{"%+.3f", 1 + 2i, "(+1.000+2.000i)"},
	{"%+.3g", 1 + 2i, "(+1+2i)"},
	{"%.3e", 0i, "(0.000e+00+0.000e+00i)"},
	{"%.3f", 0i, "(0.000+0.000i)"},
	{"%.3F", 0i, "(0.000+0.000i)"},
	{"%.3F", complex64(0i), "(0.000+0.000i)"},
	{"%.3g", 0i, "(0+0i)"},
	{"%.3e", 1 + 2i, "(1.000e+00+2.000e+00i)"},
	{"%.3f", 1 + 2i, "(1.000+2.000i)"},
	{"%.3g", 1 + 2i, "(1+2i)"},
	{"%.3e", -1 - 2i, "(-1.000e+00-2.000e+00i)"},
	{"%.3f", -1 - 2i, "(-1.000-2.000i)"},
	{"%.3g", -1 - 2i, "(-1-2i)"},
	{"% .3E", -1 - 2i, "(-1.000E+00-2.000E+00i)"},
	{"%+.3g", 1 + 2i, "(+1+2i)"},
	{"%+.3g", complex64(1 + 2i), "(+1+2i)"},
	{"%#.0f", 1.23 + 1.0i, "(1.+1.i)"},
	{"%#.0e", 1.23 + 1.0i, "(1.e+00+1.e+00i)"},
	{"%#.0g", 1.23 + 1.0i, "(1.+1.i)"},
	{"%#.0g", 0 + 100000i, "(0.+1.e+05i)"},
	{"%#.0g", 1230000 + 0i, "(1.e+06+0.i)"},
	{"%#.4f", 1 + 1.23i, "(1.0000+1.2300i)"},
	{"%#.4e", 123 + 1i, "(1.2300e+02+1.0000e+00i)"},
	{"%#.4g", 123 + 1.23i, "(123.0+1.230i)"},
	{"%#12.5g", 0 + 100000i, "(      0.0000 +1.0000e+05i)"},
	{"%#12.5g", 1230000 - 0i, "(  1.2300e+06     +0.0000i)"},
	// Precision has no effect for binary complex format.
	{"%.4b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"},
	{"%.4b", complex64(1 + 2i), "(8388608p-23+8388608p-22i)"},
	// complex infinites and NaNs
	{"%.1f", complex(posInf, posInf), "(+Inf+Infi)"},
	{"% f", complex(posInf, posInf), "( Inf+Infi)"},
	{"% f", complex(negInf, negInf), "(-Inf-Infi)"},
	{"% f", complex(NaN, NaN), "( NaN+NaNi)"},
	{"%8e", complex(posInf, posInf), "(    +Inf    +Infi)"},
	{"% 8E", complex(posInf, posInf), "(     Inf    +Infi)"},
	{"%+8f", complex(negInf, negInf), "(    -Inf    -Infi)"},
	{"% +8g", complex(negInf, negInf), "(    -Inf    -Infi)"},
	{"% -8G", complex(NaN, NaN), "( NaN    +NaN    i)"},
	{"%+-8b", complex(NaN, NaN), "(+NaN    +NaN    i)"},
	// Zero padding does not apply to infinities and NaN.
	{"%08f", complex(posInf, posInf), "(    +Inf    +Infi)"},
	{"%-08g", complex(negInf, negInf), "(-Inf    -Inf    i)"},
	{"%-08G", complex(NaN, NaN), "(NaN     +NaN    i)"},

	// old test/fmt_test.go
	{"%20.5s", "qwertyuiop", "               qwert"},
	{"%.5s", "qwertyuiop", "qwert"},
	{"%-20.5s", "qwertyuiop", "qwert               "},
	{"%20c", 'x', "                   x"},
	{"%-20c", 'x', "x                   "},
	{"%20.6e", 1.2345e3, "        1.234500e+03"},
	{"%20.6e", 1.2345e-3, "        1.234500e-03"},
	{"%20e", 1.2345e3, "        1.234500e+03"},
	{"%20e", 1.2345e-3, "        1.234500e-03"},
	{"%20.8e", 1.2345e3, "      1.23450000e+03"},
	{"%20f", 1.23456789e3, "         1234.567890"},
	{"%20f", 1.23456789e-3, "            0.001235"},
	{"%20f", 12345678901.23456789, "  12345678901.234568"},
	{"%-20f", 1.23456789e3, "1234.567890         "},
	{"%20.8f", 1.23456789e3, "       1234.56789000"},
	{"%20.8f", 1.23456789e-3, "          0.00123457"},

	// byte arrays and slices with %b,%c,%d,%o,%U and %v
	{"%012v", []byte{}, "[]"},
	{"%#012v", []byte{}, "[]byte{}"},
	{"%6v", []byte{1, 11, 111}, "[     1     11    111]"},
	{"%06v", []byte{1, 11, 111}, "[000001 000011 000111]"},
	{"%-6v", []byte{1, 11, 111}, "[1      11     111   ]"},
	{"%-06v", []byte{1, 11, 111}, "[1      11     111   ]"},
	{"%#6v", []byte{1, 11, 111}, "[]byte{   0x1,    0xb,   0x6f}"},
	{"%#06v", []byte{1, 11, 111}, "[]byte{0x000001, 0x00000b, 0x00006f}"},
	{"%#-6v", []byte{1, 11, 111}, "[]byte{0x1   , 0xb   , 0x6f  }"},
	{"%#-06v", []byte{1, 11, 111}, "[]byte{0x1   , 0xb   , 0x6f  }"},
	// f.space should and f.plus should not have an effect with %v.
	{"%# -6v", []byte{1, 11, 111}, "[]byte{ 0x1  ,  0xb  ,  0x6f }"},
	{"%#+-6v", [3]byte{1, 11, 111}, "[3]uint8{0x1   , 0xb   , 0x6f  }"},
	// f.space and f.plus should have an effect with %d.
	{"%# -6d", []byte{1, 11, 111}, "[ 1      11     111  ]"},
	{"%#+-6d", [3]byte{1, 11, 111}, "[+1     +11    +111  ]"},

	{"%-8s", []byte("\t"), "	       "},

	// slices with other formats
	{"% 02x", []byte{1}, "01"},
	{"% 02x", []byte{1, 2, 3}, "01 02 03"},

	// Padding with byte slices.
	{"%2x", []byte{}, "  "},
	{"%#2x", []byte{}, "  "},
	{"% 02x", []byte{}, "00"},
	{"%# 02x", []byte{}, "00"},
	{"%-2x", []byte{}, "  "},
	{"%-02x", []byte{}, "  "},
	{"%8x", []byte{0xab}, "      ab"},
	{"% 8x", []byte{0xab}, "      ab"},
	{"%#8x", []byte{0xab}, "    0xab"},
	{"%# 8x", []byte{0xab}, "    0xab"},
	{"%08x", []byte{0xab}, "000000ab"},
	{"% 08x", []byte{0xab}, "000000ab"},
	{"%#08x", []byte{0xab}, "00000xab"},
	{"%# 08x", []byte{0xab}, "00000xab"},
	{"%10x", []byte{0xab, 0xcd}, "      abcd"},
	{"% 10x", []byte{0xab, 0xcd}, "     ab cd"},
	{"%#10x", []byte{0xab, 0xcd}, "    0xabcd"},
	{"%# 10x", []byte{0xab, 0xcd}, " 0xab 0xcd"},
	{"%010x", []byte{0xab, 0xcd}, "000000abcd"},
	{"% 010x", []byte{0xab, 0xcd}, "00000ab cd"},
	{"%#010x", []byte{0xab, 0xcd}, "00000xabcd"},
	{"%# 010x", []byte{0xab, 0xcd}, "00xab 0xcd"},
	{"%-10X", []byte{0xab}, "AB        "},
	{"% -010X", []byte{0xab}, "AB        "},
	{"%#-10X", []byte{0xab, 0xcd}, "0XABCD    "},
	{"%# -010X", []byte{0xab, 0xcd}, "0XAB 0XCD "},
	// Same for strings
	{"%2x", "", "  "},
	{"%#2x", "", "  "},
	{"% 02x", "", "00"},
	{"%# 02x", "", "00"},
	{"%-2x", "", "  "},
	{"%-02x", "", "  "},
	{"%8x", "\xab", "      ab"},
	{"% 8x", "\xab", "      ab"},
	{"%#8x", "\xab", "    0xab"},
	{"%# 8x", "\xab", "    0xab"},
	{"%08x", "\xab", "000000ab"},
	{"% 08x", "\xab", "000000ab"},
	{"%#08x", "\xab", "00000xab"},
	{"%# 08x", "\xab", "00000xab"},
	{"%10x", "\xab\xcd", "      abcd"},
	{"% 10x", "\xab\xcd", "     ab cd"},
	{"%#10x", "\xab\xcd", "    0xabcd"},
	{"%# 10x", "\xab\xcd", " 0xab 0xcd"},
	{"%010x", "\xab\xcd", "000000abcd"},
	{"% 010x", "\xab\xcd", "00000ab cd"},
	{"%#010x", "\xab\xcd", "00000xabcd"},
	{"%# 010x", "\xab\xcd", "00xab 0xcd"},
	{"%-10X", "\xab", "AB        "},
	{"% -010X", "\xab", "AB        "},
	{"%#-10X", "\xab\xcd", "0XABCD    "},
	{"%# -010X", "\xab\xcd", "0XAB 0XCD "},

	// %T
	{"%6T", &intVar, "  *int"},
	{"%10T", nil, "     <nil>"},
	{"%-10T", nil, "<nil>     "},

	// %p with pointers
	{"%8.2p", (*int)(nil), "    0x00"},
	{"%-20.16p", &intVar, "0xPTR  "},
	// %v on pointers
	{"%8.2v", (*int)(nil), "   <nil>"},
	{"%-20.16v", &intVar, "0xPTR  "},

	// erroneous things
	{"%s %.2", "hello", "hello %!(NOVERB)"},
	// Extra argument errors should format without flags set.
	{"%010.2", "12345", "%!(NOVERB)%!(EXTRA string=12345)"},

	{"%.2f", 1.0, "1.00"},
	{"%.2f", -1.0, "-1.00"},
	{"% .2f", 1.0, " 1.00"},
	{"% .2f", -1.0, "-1.00"},
	{"%+.2f", 1.0, "+1.00"},
	{"%+.2f", -1.0, "-1.00"},
	{"%7.2f", 1.0, "   1.00"},
	{"%7.2f", -1.0, "  -1.00"},
	{"% 7.2f", 1.0, "   1.00"},
	{"% 7.2f", -1.0, "  -1.00"},
	{"%+7.2f", 1.0, "  +1.00"},
	{"%+7.2f", -1.0, "  -1.00"},
	{"% +7.2f", 1.0, "  +1.00"},
	{"% +7.2f", -1.0, "  -1.00"},
	{"%07.2f", 1.0, "0001.00"},
	{"%07.2f", -1.0, "-001.00"},
	{"% 07.2f", 1.0, " 001.00"},
	{"% 07.2f", -1.0, "-001.00"},
	{"%+07.2f", 1.0, "+001.00"},
	{"%+07.2f", -1.0, "-001.00"},
	{"% +07.2f", 1.0, "+001.00"},
	{"% +07.2f", -1.0, "-001.00"},

	// Complex numbers: exhaustively tested in TestComplexFormatting.
	{"%7.2f", 1 + 2i, "(   1.00  +2.00i)"},
	{"%+07.2f", -1 - 2i, "(-001.00-002.00i)"},

	// Use spaces instead of zero if padding to the right.
	{"%0-5s", "abc", "abc  "},
	{"%-05.1f", 1.0, "1.0  "},

	// float and complex formatting should not change the padding width
	// for other elements. See issue 14642.
	{"%06v", []interface{}{+10.0, 10}, "[000010 000010]"},
	{"%06v", []interface{}{-10.0, 10}, "[-00010 000010]"},
	{"%06v", []interface{}{+10.0 + 10i, 10}, "[(000010+00010i) 000010]"},
	{"%06v", []interface{}{-10.0 + 10i, 10}, "[(-00010+00010i) 000010]"},

	// integer formatting should not alter padding for other elements.
	{"%03.6v", []interface{}{1, 2.0, "x"}, "[000001 002 00x]"},
	{"%03.0v", []interface{}{0, 2.0, "x"}, "[    002 000]"},

	// Padding for complex numbers. Has been bad, then fixed, then bad again.
	{"%+10.2f", +104.66 + 440.51i, "(   +104.66   +440.51i)"},
	{"%+10.2f", -104.66 + 440.51i, "(   -104.66   +440.51i)"},
	{"%+10.2f", +104.66 - 440.51i, "(   +104.66   -440.51i)"},
	{"%+10.2f", -104.66 - 440.51i, "(   -104.66   -440.51i)"},
	{"%+010.2f", +104.66 + 440.51i, "(+000104.66+000440.51i)"},
	{"%+010.2f", -104.66 + 440.51i, "(-000104.66+000440.51i)"},
	{"%+010.2f", +104.66 - 440.51i, "(+000104.66-000440.51i)"},
	{"%+010.2f", -104.66 - 440.51i, "(-000104.66-000440.51i)"},

	// verbs apply to the extracted value too.
	{"%#04x", reflect.ValueOf(256), "0x0100"},
}

// zeroFill generates zero-filled strings of the specified width. The length
// of the suffix (but not the prefix) is compensated for in the width calculation.
func zeroFill(prefix string, width int, suffix string) string {
	return prefix + strings.Repeat("0", width-len(suffix)) + suffix
}

func TestSprintf(t *testing.T) {
	for _, tt := range fmtTests {
		s := Sprintf(tt.fmt, tt.val)
		i := strings.Index(tt.out, "PTR")
		if i >= 0 && i < len(s) {
			var pattern, chars string
			switch {
			case strings.HasPrefix(tt.out[i:], "PTR_b"):
				pattern = "PTR_b"
				chars = "01"
			case strings.HasPrefix(tt.out[i:], "PTR_o"):
				pattern = "PTR_o"
				chars = "01234567"
			case strings.HasPrefix(tt.out[i:], "PTR_d"):
				pattern = "PTR_d"
				chars = "0123456789"
			case strings.HasPrefix(tt.out[i:], "PTR_x"):
				pattern = "PTR_x"
				chars = "0123456789abcdef"
			case strings.HasPrefix(tt.out[i:], "PTR_X"):
				pattern = "PTR_X"
				chars = "0123456789ABCDEF"
			default:
				pattern = "PTR"
				chars = "0123456789abcdefABCDEF"
			}
			p := s[:i] + pattern
			for j := i; j < len(s); j++ {
				if !strings.ContainsRune(chars, rune(s[j])) {
					p += s[j:]
					break
				}
			}
			s = p
		}
		if s != tt.out {
			if _, ok := tt.val.(string); ok {
				// Don't requote the already-quoted strings.
				// It's too confusing to read the errors.
				t.Errorf("Sprintf(%q, %q) = <%s> want <%s>", tt.fmt, tt.val, s, tt.out)
			} else {
				t.Errorf("Sprintf(%q, %v) = %q want %q", tt.fmt, tt.val, s, tt.out)
			}
		}
	}
}

type SE []interface{} // slice of empty; notational compactness.

var reorderTests = []struct {
	fmt string
	val SE
	out string
}{
	{"%[1]d", SE{1}, "1"},
	{"%[2]d", SE{2, 1}, "1"},
	{"%[2]d %[1]d", SE{1, 2}, "2 1"},
	{"%[2]*[1]d", SE{2, 5}, "    2"},
	{"%6.2f", SE{12.0}, " 12.00"}, // Explicit version of next line.
	{"%[3]*.[2]*[1]f", SE{12.0, 2, 6}, " 12.00"},
	{"%[1]*.[2]*[3]f", SE{6, 2, 12.0}, " 12.00"},
	{"%10f", SE{12.0}, " 12.000000"},
	{"%[1]*[3]f", SE{10, 99, 12.0}, " 12.000000"},
	{"%.6f", SE{12.0}, "12.000000"}, // Explicit version of next line.
	{"%.[1]*[3]f", SE{6, 99, 12.0}, "12.000000"},
	{"%6.f", SE{12.0}, "    12"}, //  // Explicit version of next line; empty precision means zero.
	{"%[1]*.[3]f", SE{6, 3, 12.0}, "    12"},
	// An actual use! Print the same arguments twice.
	{"%d %d %d %#[1]o %#o %#o", SE{11, 12, 13}, "11 12 13 013 014 015"},
}

func TestReorder(t *testing.T) {
	for _, tt := range reorderTests {
		s := Sprintf(tt.fmt, tt.val...)
		if s != tt.out {
			t.Errorf("Sprintf(%q, %v) = <%s> want <%s>", tt.fmt, tt.val, s, tt.out)
		} else {
		}
	}
}

func args(a ...interface{}) []interface{} { return a }

var startests = []struct {
	fmt string
	in  []interface{}
	out string
}{
	{"%*d", args(4, 42), "  42"},
	{"%-*d", args(4, 42), "42  "},
	{"%*d", args(-4, 42), "42  "},
	{"%-*d", args(-4, 42), "42  "},
	{"%.*d", args(4, 42), "0042"},
	{"%*.*d", args(8, 4, 42), "    0042"},
	{"%0*d", args(4, 42), "0042"},
	// Some non-int types for width. (Issue 10732).
	{"%0*d", args(uint(4), 42), "0042"},
	{"%0*d", args(uint64(4), 42), "0042"},
	{"%0*d", args('\x04', 42), "0042"},
	{"%0*d", args(uintptr(4), 42), "0042"},
}

func TestWidthAndPrecision(t *testing.T) {
	for i, tt := range startests {
		s := Sprintf(tt.fmt, tt.in...)
		if s != tt.out {
			t.Errorf("#%d: %q: got %q expected %q", i, tt.fmt, s, tt.out)
		}
	}
}
