package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type dateLexer struct {
	lexer.Lexer
}

func ParseDate(str string) model.Date {
	l := dateLexer{Lexer: lexer.NewLexer(str)}
	return l.lexDate()
}

var strMonths = map[string]time.Month{
	"jan": time.January, "feb": time.February, "mar": time.March,
	"apr": time.April, "may": time.May, "jun": time.June,
	"jul": time.July, "aug": time.August, "sep": time.September,
	"oct": time.October, "nov": time.November, "dec": time.December,
	"january": time.January, "february": time.February, "march": time.March,
	"april": time.April, "june": time.June,
	"july": time.July, "august": time.August, "september": time.September,
	"october": time.October, "november": time.November, "december": time.December,
}

var separators = map[string]struct{}{
	".": {}, "-": {}, "/": {}, " ": {},
}

const minParts = 5

func (l *dateLexer) lexDate() model.Date {
	parts := l.lexDateParts()
	isStartOrWordBreak := true

	for i := 0; i < len(parts)-minParts+1; i++ {
		part1 := parts[i]
		if !isStartOrWordBreak {
			if part1.format == datePartNonWordChars {
				isStartOrWordBreak = true
			}

			continue
		}

		if !part1.IsNil() {
			i++

			sep := parts[i]
			if sep.format == datePartNonWordChars {
				if _, ok := separators[sep.literal]; ok {
					i++

					part2 := parts[i]
					if !part2.IsNil() {
						i++

						sep2 := parts[i]
						if sep2.literal != sep.literal {
							isStartOrWordBreak = sep2.format == datePartNonWordChars
							continue
						}

						i++

						part3 := parts[i]
						if !part3.IsNil() &&
							(i == len(parts)-1 || parts[i+1].format == datePartNonWordChars) {
							if date := findFirstValidDate(part1.Date, part2.Date, part3.Date); !date.IsNil() {
								return date
							}

							isStartOrWordBreak = false

							continue
						}

						isStartOrWordBreak = part3.format == datePartNonWordChars

						continue
					}

					isStartOrWordBreak = part2.format == datePartNonWordChars

					continue
				}

				isStartOrWordBreak = true

				continue
			}

			isStartOrWordBreak = false

			continue
		}

		isStartOrWordBreak = part1.format == datePartNonWordChars
	}

	return model.Date{}
}

func findFirstValidDate(part1, part2, part3 model.Date) model.Date {
	// Y-M-D
	if part1.Year != 0 && part2.Month != 0 && part3.Day != 0 {
		d := model.Date{Year: part1.Year, Month: part2.Month, Day: part3.Day}
		if d.IsValid() {
			return d
		}
	}
	// D-M-Y
	if part1.Day != 0 && part2.Month != 0 && part3.Year != 0 {
		d := model.Date{Year: part3.Year, Month: part2.Month, Day: part1.Day}
		if d.IsValid() {
			return d
		}
	}
	// M-D-Y
	if part1.Month != 0 && part2.Day != 0 && part3.Year != 0 {
		d := model.Date{Year: part3.Year, Month: part1.Month, Day: part2.Day}
		if d.IsValid() {
			return d
		}
	}

	return model.Date{}
}

type datePartFormat int

const (
	datePart1Digit datePartFormat = 1 + iota
	datePart2Digits
	datePart4Digits
	datePartStrMonth
	datePartWordChars
	datePartNonWordChars
)

type datePart struct {
	model.Date
	format  datePartFormat
	literal string
}

func (l *dateLexer) lexDateParts() []datePart {
	var parts []datePart
	for !l.IsEOF() {
		parts = append(parts, l.lexDatePart())
	}

	return parts
}

var (
	regex1Digit  = regexp.MustCompile(`^\d$`)
	regex2Digits = regexp.MustCompile(`^\d{2}$`)
	regex4Digits = regexp.MustCompile(`^\d{4}$`)
)

func (l *dateLexer) lexDatePart() datePart {
	str := l.ReadWhile(lexer.IsWordChar)
	if str == "" {
		str = l.ReadWhile(lexer.IsNonWordChar)

		return datePart{
			format:  datePartNonWordChars,
			literal: str,
		}
	}

	if m, ok := strMonths[strings.ToLower(str)]; ok {
		return datePart{
			Date:    model.Date{Month: m},
			format:  datePartStrMonth,
			literal: str,
		}
	}

	if regex1Digit.MatchString(str) {
		i, _ := strconv.Atoi(str)

		return datePart{
			Date:    model.Date{Day: uint8(i), Month: time.Month(i)},
			format:  datePart1Digit,
			literal: str,
		}
	}

	if regex2Digits.MatchString(str) {
		i, _ := strconv.Atoi(str)
		date := model.Date{Year: model.Year(2000 + i)}

		if i >= 1 && i <= 12 {
			date.Month = time.Month(i)
		}

		if i >= 1 && i <= 31 {
			date.Day = uint8(i)
		}

		return datePart{
			Date:    date,
			format:  datePart2Digits,
			literal: str,
		}
	}

	if regex4Digits.MatchString(str) {
		i, _ := strconv.Atoi(str)

		return datePart{
			Date:    model.Date{Year: model.Year(i)},
			format:  datePart4Digits,
			literal: str,
		}
	}

	return datePart{
		format:  datePartWordChars,
		literal: str,
	}
}
