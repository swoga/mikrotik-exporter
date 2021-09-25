package config

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/swoga/mikrotik-exporter/utils"
)

func (label *Label) AsString(log zerolog.Logger, response map[string]string, variables map[string]string) string {
	labelLog := log.With().Str("label_name", label.GetName()).Logger()
	labelLog.Trace().Msg("get label value")

	if label.Param.ParamType == PARAM_TYPE_STRING {
		value, isStaticOrInResponse := label.Param.PreprocessValue(labelLog, response, variables)
		if isStaticOrInResponse {
			return value
		}

		labelLog.Trace().Str("value", label.Param.Default).Msg("got label value from default")
		return utils.Substitute(labelLog, label.Param.Default, variables)
	} else {
		valueFloat, ok := label.Param.tryGetValue(labelLog, response, variables)
		if !ok {
			return ""
		}
		value := fmt.Sprintf("%g", valueFloat)
		labelLog.Trace().Str("value", value).Msg("got label value")
		return value
	}
}
