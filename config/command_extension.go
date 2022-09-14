package config

import "github.com/rs/zerolog"

type CommandExtension struct {
	Command string `yaml:"command"`

	Metrics     []MetricExtension  `yaml:"metrics"`
	Labels      []LabelExtension   `yaml:"labels"`
	Variables   []LabelExtension   `yaml:"variables"`
	SubCommands []CommandExtension `yaml:"sub_commands"`
}

func (x *CommandExtension) ExtendCommand(log zerolog.Logger, command Command) {
	commandExtLog := log.With().Str("command", command.CommandBase.Command).Logger()
	x.extendLabels(commandExtLog, x.Labels, command.HasLabels.Labels)
	x.extendLabels(commandExtLog, x.Variables, command.Variables)
	x.extendMetrics(commandExtLog, x.Metrics, command.Metrics)
}

func (x *CommandExtension) extendLabels(log zerolog.Logger, extensions []LabelExtension, originals Labels) {
	for _, extension := range extensions {
		extendLog := log.With().Str("label_name", extension.Label.GetName()).Logger()

		if extension.Extension.ExtensionAction == EXTENSION_ACTION_ADD {
			extendLog.Trace().Msg("add label/variable")
			originals.Add(extension.Label)
		} else {
			label, i := originals.GetByName(extension.Label.GetName())
			if label == nil {
				extendLog.Warn().Msg("label/variable not found")
				continue
			}

			switch extension.Extension.ExtensionAction {
			case EXTENSION_ACTION_OVERWRITE:
				extendLog.Trace().Msg("overwrite label/variable")
				(*originals)[i] = extension.Label
			case EXTENSION_ACTION_REMOVE:
				extendLog.Trace().Msg("remove label/variable")
				originals.RemoveByIndex(i)
			default:
				extendLog.Panic().Str("extension_action", extension.Extension.ExtensionAction).Msg("invalid extension_action")
			}
		}
	}
}

func (x *CommandExtension) extendMetrics(log zerolog.Logger, extensions []MetricExtension, originals Metrics) {
	for _, extension := range extensions {
		extendLog := log.With().Str("metric_name", extension.GetName()).Logger()

		if extension.Extension.ExtensionAction == EXTENSION_ACTION_ADD {
			extendLog.Trace().Msg("add metric")
			originals.Add(extension.Metric)
		} else {
			metric, i := originals.GetByName(extension.Metric.GetName())
			if metric == nil {
				extendLog.Warn().Msg("metric not found")
				continue
			}

			switch extension.Extension.ExtensionAction {
			case EXTENSION_ACTION_OVERWRITE:
				extendLog.Trace().Msg("overwrite metric")
				(*originals)[i] = extension.Metric
			case EXTENSION_ACTION_REMOVE:
				extendLog.Trace().Msg("remove metric")
				originals.RemoveByIndex(i)
			default:
				extendLog.Panic().Str("extension_action", extension.Extension.ExtensionAction).Msg("invalid extension_action")
			}
		}
	}
}
