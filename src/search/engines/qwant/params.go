package qwant

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/rs/zerolog/log"
)

const (
	paramKeyPage       = "offset"
	paramKeyLocale     = "locale"     // Same as Locale, only the last two characters are lowered and not everything is supported.
	paramKeySafeSearch = "safesearch" // Can be "0" or "1".

	paramCount = "count=10"
)

var validLocales = [...]string{"bg_bg", "br_fr", "ca_ad", "ca_es", "ca_fr", "co_fr", "cs_cz", "cy_gb", "da_dk", "de_at", "de_ch", "de_de", "ec_ca", "el_gr", "en_au", "en_ca", "en_gb", "en_ie", "en_my", "en_nz", "en_us", "es_ad", "es_ar", "es_cl", "es_co", "es_es", "es_mx", "es_pe", "et_ee", "eu_es", "eu_fr", "fc_ca", "fi_fi", "fr_ad", "fr_be", "fr_ca", "fr_ch", "fr_fr", "gd_gb", "he_il", "hu_hu", "it_ch", "it_it", "ko_kr", "nb_no", "nl_be", "nl_nl", "pl_pl", "pt_ad", "pt_pt", "ro_ro", "sv_se", "th_th", "zh_cn", "zh_hk"}

func localeParamString(locale options.Locale) string {
	l := strings.ToLower(locale.String())
	for _, vl := range validLocales {
		if l == vl {
			return fmt.Sprintf("%v=%v", paramKeyLocale, l)
		}
	}

	log.Warn().
		Caller().
		Str("locale", locale.String()).
		Strs("validLocales", validLocales[:]).
		Msg("Unsupported locale supplied for this engine, falling back to default")
	return fmt.Sprintf("%v=%v", paramKeyLocale, strings.ToLower(options.LocaleDefault.String()))
}

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "1")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "2")
	}
}
