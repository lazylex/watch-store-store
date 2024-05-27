package validators

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"testing"
	"time"
)

func TestOrderNumber(t *testing.T) {
	testCases := []struct {
		testName    string
		order       reservation.OrderNumber
		expectedErr error
	}{
		{
			testName:    "correct order",
			order:       465,
			expectedErr: nil,
		},
		{
			testName:    "zero order",
			order:       0,
			expectedErr: ErrIncorrectOrder,
		},
		{
			testName:    "negative order",
			order:       0,
			expectedErr: ErrIncorrectOrder,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(OrderNumber(tc.order), tc.expectedErr) {
				t.Fail()
			}
		})
	}
}

func TestArticleDTO_Validate(t *testing.T) {
	type testCase struct {
		testName    string
		article     article.Article
		expectedErr error
	}
	testCases := []testCase{
		{
			testName:    "empty article",
			article:     "",
			expectedErr: ErrIncorrectArticle,
		},
		{
			testName:    "too long article",
			article:     "Lorem ipsum dolor sit amet, consectetuer adipi.0000",
			expectedErr: ErrIncorrectArticle,
		},
		{
			testName:    "article length 50 runes",
			article:     "ᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧᾧ.1111",
			expectedErr: nil,
		},
		{
			testName:    "maximum length article",
			article:     "Lorem ipsum dolor sit amet, consectetuer adip.0000",
			expectedErr: nil,
		},
		{
			testName:    "incorrect prefix in case",
			article:     "casio..000",
			expectedErr: ErrIncorrectArticle,
		},
		{
			testName:    "incorrect prefix in display",
			article:     "casio.0a00",
			expectedErr: ErrIncorrectArticle,
		},
		{
			testName:    "incorrect prefix in open state",
			article:     "casio.00!0",
			expectedErr: ErrIncorrectArticle,
		},
		{
			testName:    "incorrect prefix in packaging",
			article:     "casio.000-",
			expectedErr: ErrIncorrectArticle,
		},
	}

	//for c := 0; c <= article.CaseWithHeavyScratches; c++ {
	//	for d := 0; d <= article.DisplayWithHeavyScratches; d++ {
	//		for p := 0; p <= article.PackageOpened; p++ {
	//			for pn := 0; pn <= article.PackagingWithDamage; pn++ {
	//				correctArt := fmt.Sprintf("correct article .%d%d%d%d", c, d, p, pn)
	//				testCases = append(testCases, testCase{
	//					testName:    correctArt,
	//					article:     article.Article(correctArt),
	//					expectedErr: nil,
	//				})
	//			}
	//		}
	//	}
	//}
	//
	//for c := article.CaseWithHeavyScratches + 1; c < 10; c++ {
	//	for d := article.DisplayWithHeavyScratches + 1; d < 10; d++ {
	//		for p := article.PackageOpened + 1; p < 10; p++ {
	//			for pn := article.PackagingWithDamage; pn < 10; pn++ {
	//				incorrectArt := fmt.Sprintf("incorrect article .%d%d%d%d", c, d, p, pn)
	//				testCases = append(testCases, testCase{
	//					testName:    incorrectArt,
	//					article:     article.Article(incorrectArt),
	//					expectedErr: ErrIncorrectArticle,
	//				})
	//			}
	//		}
	//	}
	//}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(Article(tc.article), tc.expectedErr) {
				t.Fail()
			}
		})
	}
}

func TestDateOrder(t *testing.T) {
	testCases := []struct {
		testName    string
		from        time.Time
		to          time.Time
		expectedErr error
	}{
		{
			testName:    "valid order",
			from:        time.Date(1987, 1, 19, 14, 0, 0, 0, time.UTC),
			to:          time.Date(2003, 1, 20, 10, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		{
			testName:    "invalid order",
			from:        time.Date(2003, 1, 20, 10, 0, 0, 0, time.UTC),
			to:          time.Date(1987, 1, 19, 14, 0, 0, 0, time.UTC),
			expectedErr: ErrIncorrectDatesOrder,
		},
		{
			testName:    "equal date",
			from:        time.Date(2003, 1, 20, 10, 0, 0, 0, time.UTC),
			to:          time.Date(2003, 1, 20, 10, 0, 0, 0, time.UTC),
			expectedErr: ErrDatesIsEqual,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(DateOrder(tc.from, tc.to), tc.expectedErr) {
				t.Fail()
			}
		})
	}
}

func TestPrice(t *testing.T) {
	testCases := []struct {
		testName    string
		price       float64
		expectedErr error
	}{
		{
			testName:    "zero price",
			price:       0,
			expectedErr: ErrZeroPrice,
		},
		{
			testName:    "negative price",
			price:       -100,
			expectedErr: ErrNegativePrice,
		},
		{
			testName:    "correct price",
			price:       999.99,
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(Price(tc.price), tc.expectedErr) {
				t.Fail()
			}
		})
	}
}

func TestName(t *testing.T) {
	testCases := []struct {
		testName    string
		name        string
		expectedErr error
	}{
		{
			testName:    "empty testName",
			name:        "",
			expectedErr: ErrEmptyName,
		},
		{
			testName:    "correct testName",
			name:        "Casio F-91W",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if !errors.Is(Name(tc.name), tc.expectedErr) {
				t.Fail()
			}
		})
	}
}
